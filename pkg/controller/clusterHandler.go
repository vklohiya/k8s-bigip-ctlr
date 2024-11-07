package controller

import (
	"context"
	"fmt"
	ficV1 "github.com/F5Networks/f5-ipam-controller/pkg/ipamapis/apis/fic/v1"
	"github.com/F5Networks/f5-ipam-controller/pkg/ipammachinery"
	"github.com/F5Networks/k8s-bigip-ctlr/v2/pkg/clustermanager"
	log "github.com/F5Networks/k8s-bigip-ctlr/v2/pkg/vlogger"
	routeclient "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	v1 "k8s.io/api/core/v1"
	extClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"os"
	"reflect"
	"strings"
	"time"
	"unicode"
)

// NewClusterHandler initializes the ClusterHandler with the required structures for each cluster.
func NewClusterHandler(params Params) *ClusterHandler {
	handler := &ClusterHandler{
		ClusterConfigs:      make(map[string]*ClusterConfig),
		UniqueAppIdentifier: make(map[string]struct{}),
		EventChan:           make(chan *rqKey, 1),
		NodeLabel:           params.NodeLabelSelector,
		NamespaceLabel:      params.NamespaceLabel,
		RouteLabel:          params.RouteLabel,
		CustomResouceLabel:  DefaultCustomResourceLabel,
		LocalClusterName:    params.LocalClusterName,
		resourceQueue: workqueue.NewNamedRateLimitingQueue(
			workqueue.DefaultControllerRateLimiter(), "nextgen-resource-controller"),
		processedHostPath: &ProcessedHostPath{
			processedHostPathMap: make(map[string]metav1.Time),
		},
	}
	return handler
}

// fetch cluster name for given secret if it holds kubeconfig of the cluster.
func (ch *ClusterHandler) getClusterForSecret(name, namespace string) ClusterDetails {
	ch.RLock()
	defer ch.RUnlock()
	for _, mcc := range ch.ClusterConfigs {
		// Skip empty/nil configs processing
		if mcc.clusterDetails == (ClusterDetails{}) {
			continue
		}
		// Check if the secret holds the kubeconfig for a cluster by checking if it's referred in the multicluster config
		// if so then return the cluster name associated with the secret
		if mcc.clusterDetails.Secret == (namespace + "/" + name) {
			return mcc.clusterDetails
		}
	}
	return ClusterDetails{}
}

// addClusterConfig adds a new cluster configuration to the ClusterHandler.
func (ch *ClusterHandler) addClusterConfig(clusterName string, config *ClusterConfig) {
	ch.Lock()
	config.clusterName = clusterName
	ch.ClusterConfigs[clusterName] = config
	ch.Unlock()
}

// deleteClusterConfig removes a cluster configuration from the ClusterHandler.
func (ch *ClusterHandler) deleteClusterConfig(clusterName string) {
	ch.Lock()
	delete(ch.ClusterConfigs, clusterName)
	ch.Unlock()
}

// getClusterConfig returns the cluster configuration for the specified cluster.
func (ch *ClusterHandler) getClusterConfig(clusterName string) *ClusterConfig {
	ch.RLock()
	defer ch.RUnlock()
	if _, exists := ch.ClusterConfigs[clusterName]; !exists {
		return nil
	}
	return ch.ClusterConfigs[clusterName]
}

// addInformerStore adds a new InformerStore to the ClusterHandler.
func (ch *ClusterHandler) addInformerStore(clusterName string, store *InformerStore) {
	ch.Lock()
	if clusterConfig, ok := ch.ClusterConfigs[clusterName]; !ok {
		clusterConfig.InformerStore = store
	}
	ch.Unlock()
}

// getInformerStore returns the InformerStore for the specified cluster.
func (ch *ClusterHandler) getInformerStore(clusterName string) *InformerStore {
	ch.RLock()
	defer ch.RUnlock()
	if clusterConfig, exists := ch.ClusterConfigs[clusterName]; exists {
		return clusterConfig.InformerStore
	}
	return nil
}

// processEvent handles individual events, simulating sending to a controller.
func (ch *ClusterHandler) processEvent(obj interface{}) {
	// Here you would handle the business logic for the event.
	fmt.Printf("Processing event: %v\n", obj)
	// Add actual controller handling logic here.
}

// remove any cluster which is not provided in externalClustersConfig or not part of the HA cluster
func (ch *ClusterHandler) cleanClusterCache(primaryClusterName, secondaryClusterName string, activeClusters map[string]struct{}) {
	ch.Lock()
	defer ch.Unlock()
	for clusterName := range ch.ClusterConfigs {
		// Avoid deleting HA cluster related configs
		if clusterName == primaryClusterName || clusterName == secondaryClusterName || clusterName == "" {
			continue
		}
		// Avoid deleting active clusters
		if _, exists := activeClusters[clusterName]; exists {
			continue
		}
		log.Infof("[MultiCluster] Removing the cluster config for cluster %s from CIS Cache", clusterName)
		delete(ch.ClusterConfigs, clusterName)
	}
}

// function to get the count of edns resources from all the active clusters for a namespace
func (ch *ClusterHandler) getEDNSCount(ns string) int {
	rscCount := 0
	for _, infSet := range ch.ClusterConfigs {
		if comInf, ok := infSet.comInformers[ns]; ok {
			edns, err := comInf.ednsInformer.GetIndexer().ByIndex("namespace", ns)
			if err != nil {
				continue
			}
			rscCount += len(edns)
		}
	}
	return rscCount
}

// function to return if cluster informers are ready
func (ch *ClusterHandler) isClusterInformersReady() bool {
	if ch.ClusterConfigs == nil {
		return false
	}
	return true
}

// function to return the list of all nodes from all the active clusters
func (ch *ClusterHandler) getAllNodesUsingInformers() []interface{} {
	ch.RLock()
	defer ch.RUnlock()
	var nodes []interface{}
	for _, infSet := range ch.ClusterConfigs {
		nodes = append(nodes, infSet.nodeInformer.nodeInformer.GetIndexer().List()...)
	}
	return nodes
}

// function to return the list of all nodes using the rest client
func (ch *ClusterHandler) getAllNodesUsingRestClient() []interface{} {
	ch.RLock()
	defer ch.RUnlock()
	var nodes []interface{}
	for clusterName, config := range ch.ClusterConfigs {
		nodesObj, err := config.kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: config.nodeLabelSelector})
		if err != nil {
			log.Debugf("[MultiCluster] Unable to fetch nodes for cluster %v with err %v", clusterName, err)
		} else {
			for _, node := range nodesObj.Items {
				nodes = append(nodes, &node)
			}
		}
	}
	return nodes
}

// function to return the cluster counts
func (ch *ClusterHandler) getClusterCount() int {
	ch.RLock()
	defer ch.RUnlock()
	return len(ch.ClusterConfigs)
}

// function to return the list of monitored namespaces for ingress resources
func (ch *ClusterHandler) getMonitoredNamespaces(clusterName string) map[string]bool {
	ch.RLock()
	defer ch.RUnlock()
	ns := make(map[string]bool)
	if config, ok := ch.ClusterConfigs[clusterName]; ok {
		ns = config.namespaces
	}
	return ns
}

func (ch *ClusterHandler) getClusterNames() map[string]struct{} {
	ch.RLock()
	defer ch.RUnlock()
	clusterNames := make(map[string]struct{})
	for clusterName, config := range ch.ClusterConfigs {
		if config.clusterDetails.ServiceTypeLBDiscovery {
			clusterNames[clusterName] = struct{}{}
		}
	}
	return clusterNames
}

// setupClientsforCluster sets Kubernetes Clients.
func (ch *ClusterHandler) setupClientsforCluster(config *rest.Config, ipamClient bool, clusterName string, clusterConfig *ClusterConfig, mode ControllerMode) error {
	kubeCRClient, err := clustermanager.CreateKubeCRClientFromKubeConfig(config)
	if err != nil {
		return fmt.Errorf("Failed to create Custom Resource kubeClient: %v", err)
	}

	kubeClient, err := clustermanager.CreateKubeClientFromKubeConfig(config)
	if err != nil {
		return fmt.Errorf("Failed to create kubeClient: %v", err)
	}

	var kubeIPAMClient *extClient.Clientset
	if ipamClient {
		kubeIPAMClient, err = clustermanager.CreateKubeIPAMClientFromKubeConfig(config)
		if err != nil {
			log.Errorf("Failed to create ipam client: %v", err)
		}
	}

	var rclient *routeclient.RouteV1Client
	if mode == OpenShiftMode {
		rclient, err = clustermanager.CreateRouteClientFromKubeconfig(config)
		if nil != err {
			return fmt.Errorf("Failed to create Route Client: %v", err)
		}
	}

	log.Debugf("Clients Created for cluster: %s", clusterName)
	//Update the clusterConfig store
	clusterConfig.kubeClient = kubeClient
	clusterConfig.kubeCRClient = kubeCRClient
	clusterConfig.kubeAPIClient = kubeIPAMClient
	clusterConfig.routeClientV1 = rclient
	return nil
}

func (ch *ClusterHandler) createNamespaceLabeledInformerForCluster(label string, clusterConfig *ClusterConfig) error {
	selector, err := createLabelSelector(label)
	if err != nil {
		return fmt.Errorf("unable to setup namespace-label informer for label: %v, Error:%v", label, err)
	}
	namespaceOptions := func(options *metav1.ListOptions) {
		options.LabelSelector = selector.String()
	}

	if clusterConfig == nil {
		return fmt.Errorf("no cluster config found in cache")
	}
	if clusterConfig != nil && clusterConfig.InformerStore != nil && 0 != len(clusterConfig.InformerStore.crInformers) {
		return fmt.Errorf("cannot set a namespace label informer when informers " +
			"have been setup for one or more namespaces")
	}

	resyncPeriod := 0 * time.Second
	restClientv1 := clusterConfig.kubeClient.CoreV1().RESTClient()
	clusterConfig.InformerStore.nsInformers[label] = &NSInformer{
		stopCh: make(chan struct{}),
		nsInformer: cache.NewSharedIndexInformer(
			cache.NewFilteredListWatchFromClient(
				restClientv1,
				"namespaces",
				"",
				namespaceOptions,
			),
			&v1.Namespace{},
			resyncPeriod,
			cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
		),
	}

	clusterConfig.InformerStore.nsInformers[label].nsInformer.AddEventHandlerWithResyncPeriod(
		&cache.ResourceEventHandlerFuncs{
			AddFunc:    func(obj interface{}) { clusterConfig.enqueueNamespace(obj) },
			DeleteFunc: func(obj interface{}) { clusterConfig.enqueueDeletedNamespace(obj) },
		},
		resyncPeriod,
	)

	return nil
}

// function to initialize the cluster config for a cluster
func (ch *ClusterHandler) initClusterConfig(clusterName string, mode ControllerMode, config *rest.Config, ipamClient bool, namespaces []string, ipamNamespace, ipamClusterLabel string) {
	clusterConfig := newClusterConfig()
	clusterConfig.EventChan = ch.EventChan
	clusterConfig.InformerStore = initInformerStore()
	clusterConfig.nodeLabelSelector = ch.NodeLabel
	clusterConfig.nativeResourceSelector, _ = createLabelSelector(DefaultNativeResourceLabel)
	clusterConfig.customResourceSelector, _ = createLabelSelector(ch.CustomResouceLabel)
	clusterConfig.routeLabel = ch.RouteLabel
	if err := ch.setupClientsforCluster(config, ipamClient, "", clusterConfig, mode); err != nil {
		log.Errorf("Failed to Setup Clients: %v", err)
	}
	// add the cluster config for local cluster
	ch.addClusterConfig(clusterName, clusterConfig)
	if ch.NamespaceLabel == "" {
		if len(namespaces) == 0 {
			clusterConfig.namespaces[""] = true
			log.Debug("No namespaces provided. Watching all namespaces")
		} else {
			for _, ns := range namespaces {
				clusterConfig.namespaces[ns] = true
			}
		}
	} else {
		err2 := ch.createNamespaceLabeledInformerForCluster(ch.NamespaceLabel, clusterConfig)
		if err2 != nil {
			log.Errorf("%v", err2)
			ch.RLock()
			for _, nsInf := range clusterConfig.nsInformers {
				for _, v := range nsInf.nsInformer.GetIndexer().List() {
					ns := v.(*v1.Namespace)
					clusterConfig.namespaces[ns.ObjectMeta.Name] = true
				}
			}
			ch.RUnlock()
		}
	}
	if ipamClient {
		if !ch.validateIPAMConfig(ipamNamespace, clusterConfig) {
			log.Warningf("[IPAM] IPAM Namespace %s not found in the list of monitored namespaces", ipamNamespace)
		}
		ipamParams := ipammachinery.Params{
			Config:        config,
			EventHandlers: clusterConfig.getEventHandlerForIPAM(),
			Namespaces:    []string{ipamNamespace},
		}

		clusterConfig.ipamClient = ipammachinery.NewIPAMClient(ipamParams)
		clusterConfig.ipamClusterLabel = ipamClusterLabel
		if ipamClusterLabel != "" {
			clusterConfig.ipamClusterLabel = ipamClusterLabel + "/"
		}
		ch.registerIPAMCRD(clusterConfig)
		time.Sleep(3 * time.Second)
		_ = ch.createIPAMResource(ipamNamespace, clusterConfig)
	}
}

// validate IPAM configuration
func (ch *ClusterHandler) validateIPAMConfig(ipamNamespace string, clusterConfig *ClusterConfig) bool {
	// verify the ipam configuration
	for ns, _ := range clusterConfig.namespaces {
		if ns == "" {
			return true
		} else {
			if ns == ipamNamespace {
				return true
			}
		}
	}
	return false
}

// Register IPAM CRD
func (ctlr *Controller) registerIPAMCRD() {
	clusterConfig := ctlr.multiClusterHandler.getClusterConfig("")
	err := ipammachinery.RegisterCRD(clusterConfig.kubeAPIClient)
	if err != nil {
		log.Errorf("[IPAM] error while registering CRD %v", err)
	}
}

func (ch *ClusterHandler) getIPAMClient(clusterName string) *ipammachinery.IPAMClient {
	ch.RLock()
	defer ch.RUnlock()
	var ipamClient *ipammachinery.IPAMClient
	if clusterConfig := ch.getClusterConfig(clusterName); clusterConfig != nil {
		ipamClient = clusterConfig.ipamClient
	}
	return ipamClient
}

func (ch *ClusterHandler) setIPAMClient(clusterName string, ipamClient *ipammachinery.IPAMClient) {
	ch.RLock()
	defer ch.RUnlock()
	if clusterConfig := ch.getClusterConfig(clusterName); clusterConfig != nil {
		clusterConfig.ipamClient = ipamClient
	}
}

// function to get the IPAM CR
func (ch *ClusterHandler) getIPAMCR(clusterName string) *ficV1.IPAM {
	ch.RLock()
	defer ch.RUnlock()
	if clusterConfig := ch.getClusterConfig(clusterName); clusterConfig != nil && clusterConfig.ipamClient != nil && clusterConfig.ipamCR != "" {
		cr := strings.Split(clusterConfig.ipamCR, "/")
		ipamCR, err := clusterConfig.ipamClient.Get(cr[0], cr[1])
		if err == nil {
			return ipamCR
		}
		log.Errorf("[IPAM] error while retrieving IPAM custom resource: %s.", err.Error())
	}
	log.Debugf("[IPAM] IPAM Client not found for cluster %s", clusterName)
	return nil
}

// Register IPAM CRD
func (ch *ClusterHandler) registerIPAMCRD(clusterConfig *ClusterConfig) {
	err := ipammachinery.RegisterCRD(clusterConfig.kubeAPIClient)
	if err != nil {
		log.Errorf("[IPAM] error while registering CRD %v", err)
	}
}

// Create IPAM CRD
func (ch *ClusterHandler) createIPAMResource(ipamNamespace string, clusterConfig *ClusterConfig) error {

	frameIPAMResourceName := func() string {
		prtn := ""
		for _, ch := range DEFAULT_PARTITION {
			elem := string(ch)
			if unicode.IsUpper(ch) {
				elem = strings.ToLower(elem) + "-"
			}
			prtn += elem
		}
		if string(prtn[len(prtn)-1]) == "-" {
			prtn = prtn + ipamCRName
		} else {
			prtn = prtn + "." + ipamCRName
		}

		prtn = strings.Replace(prtn, "_", "-", -1)
		prtn = strings.Replace(prtn, "--", "-", -1)

		hostsplit := strings.Split(os.Getenv("HOSTNAME"), "-")
		var host string
		if len(hostsplit) > 2 {
			host = strings.Join(hostsplit[0:len(hostsplit)-2], "-")
		} else {
			host = strings.Join(hostsplit, "-")
		}
		return strings.Join([]string{host, prtn}, ".")
	}

	crName := frameIPAMResourceName()
	f5ipam := &ficV1.IPAM{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName,
			Namespace: ipamNamespace,
		},
		Spec: ficV1.IPAMSpec{
			HostSpecs: make([]*ficV1.HostSpec, 0),
		},
		Status: ficV1.IPAMStatus{
			IPStatus: make([]*ficV1.IPSpec, 0),
		},
	}
	clusterConfig.ipamCR = ipamNamespace + "/" + crName

	ipamCR, err := clusterConfig.ipamClient.Create(f5ipam)
	if err == nil {
		log.Debugf("[IPAM] Created IPAM Custom Resource: \n%v\n", ipamCR)
		return nil
	}

	log.Debugf("[IPAM] error while creating IPAM custom resource %v", err.Error())
	return err
}

// ClusterEventHandler function to handle the events
func (ch *ClusterHandler) ClusterEventHandler() {
	for req := range ch.EventChan {
		switch req.kind {
		case IPAM:
			ch.resourceQueue.Add(req)
		case Namespace:
			ch.resourceQueue.Add(req)
		default:
			log.Errorf("unknown event type received %v", req)
		}
	}
}

func (cc *ClusterConfig) getEventHandlerForIPAM() *cache.ResourceEventHandlerFuncs {
	return &cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { cc.enqueueIPAM(obj) },
		UpdateFunc: func(oldObj, newObj interface{}) { cc.enqueueUpdatedIPAM(oldObj, newObj) },
		DeleteFunc: func(obj interface{}) { cc.enqueueDeletedIPAM(obj) },
	}
}

func (cc *ClusterConfig) enqueueIPAM(obj interface{}) {
	ipamObj := obj.(*ficV1.IPAM)

	if ipamObj.Namespace+"/"+ipamObj.Name != cc.ipamCR {
		return
	}

	log.Debugf("Enqueueing IPAM: %v", ipamObj)
	key := &rqKey{
		namespace:   ipamObj.ObjectMeta.Namespace,
		kind:        IPAM,
		rscName:     ipamObj.ObjectMeta.Name,
		rsc:         obj,
		event:       Create,
		clusterName: cc.clusterName,
	}

	cc.EventChan <- key
}

func (cc *ClusterConfig) enqueueUpdatedIPAM(oldObj, newObj interface{}) {
	oldIpam := oldObj.(*ficV1.IPAM)
	curIpam := newObj.(*ficV1.IPAM)

	if curIpam.Namespace+"/"+curIpam.Name != cc.ipamCR {
		return
	}

	if reflect.DeepEqual(oldIpam.Status, curIpam.Status) {
		return
	}

	log.Debugf("Enqueueing Updated IPAM: %v", curIpam)
	key := &rqKey{
		namespace:   curIpam.ObjectMeta.Namespace,
		kind:        IPAM,
		rscName:     curIpam.ObjectMeta.Name,
		rsc:         newObj,
		event:       Update,
		clusterName: cc.clusterName,
	}
	cc.EventChan <- key
}

func (cc *ClusterConfig) enqueueDeletedIPAM(obj interface{}) {
	ipamObj := obj.(*ficV1.IPAM)

	if ipamObj.Namespace+"/"+ipamObj.Name != cc.ipamCR {
		return
	}

	log.Debugf("Enqueueing IPAM: %v", ipamObj)
	key := &rqKey{
		namespace:   ipamObj.ObjectMeta.Namespace,
		kind:        IPAM,
		rscName:     ipamObj.ObjectMeta.Name,
		rsc:         obj,
		event:       Delete,
		clusterName: cc.clusterName,
	}
	cc.EventChan <- key
}

func (cc *ClusterConfig) enqueueNamespace(obj interface{}) {
	ns := obj.(*v1.Namespace)
	log.Debugf("Enqueueing Namespace: %v", ns)
	key := &rqKey{
		namespace:   ns.ObjectMeta.Namespace,
		kind:        Namespace,
		rscName:     ns.ObjectMeta.Name,
		rsc:         obj,
		event:       Create,
		clusterName: cc.clusterName,
	}
	cc.EventChan <- key
}

func (cc *ClusterConfig) enqueueDeletedNamespace(obj interface{}) {
	ns := obj.(*v1.Namespace)
	log.Debugf("Enqueueing Namespace: %v on Delete", ns)
	key := &rqKey{
		namespace:   ns.ObjectMeta.Namespace,
		kind:        Namespace,
		rscName:     ns.ObjectMeta.Name,
		rsc:         obj,
		event:       Delete,
		clusterName: cc.clusterName,
	}
	cc.EventChan <- key
}
