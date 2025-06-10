/*
Copyright 2025 F5 Networks.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// VirtualServerStatus defines the observed state of VirtualServer.
type VirtualServerStatus struct {
	CustomResourceStatus
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// VirtualServer is the Schema for the virtualservers API.
type VirtualServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualServerSpec   `json:"spec,omitempty"`
	Status VirtualServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualServerList contains a list of VirtualServer.
type VirtualServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualServer{}, &VirtualServerList{})
}

// +kubebuilder:validation:XValidation:rule="!has(self.partition) || self.partition != 'Common'",message="The partition cannot be 'Common' if specified."
// +kubebuilder:validation:XValidation:rule="has(self.partition) == has(oldSelf.partition) && (!has(self.partition) || self.partition == oldSelf.partition)",message="partition cannot be modified. Delete the resource and recreate with new partition"
// +kubebuilder:validation:XValidation:rule="!(has(self.serviceAddress) && !has(oldSelf.serviceAddress))",message="'serviceAddress' cannot be added when it is not already present."
// +kubebuilder:validation:XValidation:rule="!(has(oldSelf.serviceAddress) && !has(self.serviceAddress))",message="'serviceAddress' cannot be deleted when it is present."
// +kubebuilder:validation:XValidation:rule="has(self.ipamLabel) || has(self.virtualServerAddress)",message="either ipamLabel or virtualServerAddress needs to be specified."
type VirtualServerSpec struct {
	// +kubebuilder:validation:Pattern=`^(([a-zA-Z0-9\*]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`
	Host string `json:"host,omitempty"`
	// +kubebuilder:validation:Pattern=`^(([a-zA-Z0-9\*]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`
	HostAliases []string `json:"hostAliases,omitempty"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]+[-A-z0-9_.:]*[A-z0-9]*$`
	HostGroup string `json:"hostGroup,omitempty"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]+([A-z0-9-._+])*([A-z0-9])$`
	HostGroupVirtualServerName string `json:"hostGroupVirtualServerName,omitempty"`
	// +kubebuilder:validation:Pattern=`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])|(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[ PSYCH!0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	VirtualServerAddress string `json:"virtualServerAddress,omitempty"`
	// +kubebuilder:validation:Pattern=`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])|(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	AdditionalVirtualServerAddresses []string `json:"additionalVirtualServerAddresses,omitempty"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]+[-A-z0-9_.:]+[A-z0-9]+$`
	IPAMLabel string `json:"ipamLabel,omitempty"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65535
	BigIPRouteDomain int32 `json:"bigipRouteDomain,omitempty"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]+([A-z0-9-._+])*([A-z0-9])$`
	VirtualServerName string `json:"virtualServerName,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	VirtualServerHTTPPort int32 `json:"virtualServerHTTPPort,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	VirtualServerHTTPSPort int32       `json:"virtualServerHTTPSPort,omitempty"`
	DefaultPool            DefaultPool `json:"defaultPool,omitempty"`
	Pools                  []VSPool    `json:"pools,omitempty"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]+[-A-z0-9_.:]+[A-z0-9]+$`
	TLSProfileName string `json:"tlsProfileName,omitempty"`
	// +kubebuilder:validation:Enum=allow;none;redirect
	HTTPTraffic string `json:"httpTraffic,omitempty"`
	// +kubebuilder:validation:Pattern=`^$|^\/?[a-zA-Z]+([-A-z0-9_+]+\/)*([-A-z0-9_.:]+\/?)+$`
	SNAT string `json:"snat,omitempty"`
	// +kubebuilder:validation:Enum=none;L4
	ConnectionMirroring string `json:"connectionMirroring,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	WAF string `json:"waf,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)*([-A-z0-9_.:]+\/?)*$`
	RewriteAppRoot string `json:"rewriteAppRoot,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.]+\/?)*$`
	AllowVLANs []string `json:"allowVlans,omitempty"`
	// +kubebuilder:validation:Pattern=`^none$|^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	IRules []string `json:"iRules,omitempty"`
	// +kubebuilder:validation:MaxItems=1
	ServiceIPAddress []ServiceAddress `json:"serviceAddress,omitempty"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]+[-A-z0-9_.:]+[A-z0-9]+$`
	PolicyName string `json:"policyName,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/?[a-zA-Z]+([-A-z0-9_+]+\/)*([-A-z0-9_.:]+\/?)*$`
	PersistenceProfile string `json:"persistenceProfile,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	HTTPCompressionProfile string `json:"httpCompressionProfile,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	ProfileMultiplex string `json:"profileMultiplex,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	DOS string `json:"dos,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	BotDefense            string        `json:"botDefense,omitempty"`
	Profiles              ProfileVSSpec `json:"profiles,omitempty"`
	AllowSourceRange      []string      `json:"allowSourceRange,omitempty"`
	HttpMrfRoutingEnabled *bool         `json:"httpMrfRoutingEnabled,omitempty"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]+[-A-z0-9_.]+$`
	Partition string `json:"partition,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	HTMLProfile     string          `json:"htmlProfile,omitempty"`
	HostPersistence HostPersistence `json:"hostPersistence,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	ProfileAccess string `json:"profileAccess,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	PolicyPerRequestAccess string       `json:"policyPerRequestAccess,omitempty"`
	ProfileAdapt           ProfileAdapt `json:"profileAdapt,omitempty"`
}

type HostPersistence struct {
	// +kubebuilder:validation:Enum=sourceAddress;destinationAddress;cookieInsert;cookieRewrite;cookiePassive;cookieHash;universal;hash;carp;none
	Method          string          `json:"method,omitempty"`
	PersistMetaData PersistMetaData `json:"metaData,omitempty"`
}

type PersistMetaData struct {
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Pattern=`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])|(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	Netmask string `json:"netmask,omitempty"`
	Key     string `json:"key,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Timeout int32 `json:"timeout,omitempty"`
	// +kubebuilder:validation:Pattern=`^((?:(?:[0-9]+d))|(?:(?:[0-9]+d)?((?:[01]?[0-9]|2[0-3]):[0-5][0-9](?::[0-5][0-9])?)))$`
	Expiry string `json:"expiry,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Offset int32 `json:"offset,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Length int32 `json:"length,omitempty"`
}

type ServiceAddress struct {
	ArpEnabled bool `json:"arpEnabled,omitempty"`
	// +kubebuilder:validation:Enum=enable;disable;selective
	ICMPEcho string `json:"icmpEcho,omitempty"`
	// +kubebuilder:validation:Enum=enable;disable;selective;always;any;all
	RouteAdvertisement string `json:"routeAdvertisement,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	TrafficGroup    string `json:"trafficGroup,omitempty"`
	SpanningEnabled bool   `json:"spanningEnabled,omitempty"`
}

type DefaultPool struct {
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-._+]+\/)+([-A-z0-9_.:]+\/?)*$`
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Pattern=`[a-z]([-a-z0-9]*[a-z0-9])?`
	Service     string             `json:"service"`
	ServicePort intstr.IntOrString `json:"servicePort"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z0-9][-A-Za-z0-9_.\/]{0,61}[a-zA-Z0-9]=(\s?|""|[a-zA-Z0-9][-A-Za-z0-9_.]{0,61}[a-zA-Z0-9])$`
	NodeMemberLabel string    `json:"nodeMemberLabel,omitempty"`
	Monitors        []Monitor `json:"monitors"`
	// +kubebuilder:validation:Pattern=`^[a-z]+[a-z_-]+[a-z]+$`
	Balance string `json:"loadBalancingMethod,omitempty"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]+([-A-z0-9_.+:])*([A-z0-9])+$`
	ServiceNamespace string `json:"serviceNamespace,omitempty"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65535
	ReselectTries     int32  `json:"reselectTries,omitempty"`
	ServiceDownAction string `json:"serviceDownAction,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=bigip;service
	Reference string `json:"reference"`
}

type VSPool struct {
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]+([-A-z0-9_.+:])*([A-z0-9])+$`
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/([A-z0-9-_+]+\/)*([-A-z0-9_.:]+\/?)*$`
	Path string `json:"path,omitempty"`
	// +kubebuilder:validation:Pattern=`[a-z]([-a-z0-9]*[a-z0-9])?`
	Service     string             `json:"service"`
	ServicePort intstr.IntOrString `json:"servicePort"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z0-9][-A-Za-z0-9_.\/]{0,61}[a-zA-Z0-9]=(\s?|""|[a-zA-Z0-9][-A-Za-z0-9_.]{0,61}[a-zA-Z0-9])$`
	NodeMemberLabel string             `json:"nodeMemberLabel,omitempty"`
	Monitor         Monitor            `json:"monitor"`
	Monitors        []Monitor          `json:"monitors"`
	MinimumMonitors intstr.IntOrString `json:"minimumMonitors,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/([A-z0-9-_+]+\/)*([-A-z0-9_.:]+\/?)*$`
	Rewrite string `json:"rewrite,omitempty"`
	// +kubebuilder:validation:Pattern=`^[a-z]+[a-z_-]+[a-z]+$`
	Balance string `json:"loadBalancingMethod,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/([A-z0-9-_+]+\/)+([A-z0-9]+\/?)*$`
	WAF string `json:"waf,omitempty"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]+([-A-z0-9_.+:])*([A-z0-9])+$`
	ServiceNamespace string `json:"serviceNamespace,omitempty"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65535
	ReselectTries     int32  `json:"reselectTries,omitempty"`
	ServiceDownAction string `json:"serviceDownAction,omitempty"`
	// +kubebuilder:validation:Pattern=`^(([a-zA-Z0-9\*]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`
	HostRewrite string `json:"hostRewrite,omitempty"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=256
	Weight               *int32                         `json:"weight,omitempty"`
	AlternateBackends    []AlternateBackend             `json:"alternateBackends"`
	MultiClusterServices []MultiClusterServiceReference `json:"multiClusterServices,omitempty"`
}
