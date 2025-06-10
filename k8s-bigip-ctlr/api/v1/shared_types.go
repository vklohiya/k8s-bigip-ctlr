package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type CustomResourceStatus struct {
	VSAddress   string      `json:"vsAddress,omitempty"`
	Status      string      `json:"status,omitempty"`
	LastUpdated metav1.Time `json:"lastUpdated,omitempty"`
	Error       string      `json:"error,omitempty"`
}
type AlternateBackend struct {
	// +kubebuilder:validation:Pattern=`[a-z]([-a-z0-9]*[a-z0-9])?`
	Service string `json:"service"`
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]+([-A-z0-9_.+:])*([A-z0-9])+$`
	ServiceNamespace string `json:"serviceNamespace,omitempty"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=256
	Weight *int32 `json:"weight,omitempty"`
}

type MultiClusterServiceReference struct {
	// +kubebuilder:validation:Required
	ClusterName string `json:"clusterName"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`[a-z]([-a-z0-9]*[a-z0-9])?`
	SvcName string `json:"service"`
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`
	// +kubebuilder:validation:Required
	ServicePort intstr.IntOrString `json:"servicePort"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=256
	Weight *int `json:"weight,omitempty"`
}

type Monitor struct {
	// +kubebuilder:validation:Enum=tcp;udp;http;https
	Type       string `json:"type,omitempty"`
	Send       string `json:"send,omitempty"`
	Recv       string `json:"recv,omitempty"`
	Interval   int    `json:"interval,omitempty"`
	Timeout    int    `json:"timeout,omitempty"`
	TargetPort int32  `json:"targetPort,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Enum=bigip
	Reference string `json:"reference,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/([A-z0-9-_+]+\/)+([A-z0-9]+\/?)*$`
	SSLProfile string `json:"sslProfile,omitempty"`
}

type ProfileVSSpec struct {
	TCP   ProfileTCP   `json:"tcp,omitempty"`
	HTTP2 ProfileHTTP2 `json:"http2,omitempty"`
}

type ProfileTSSpec struct {
	TCP ProfileTCP `json:"tcp,omitempty"`
}

type ProfileTCP struct {
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	Client string `json:"client,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	Server string `json:"server,omitempty"`
}

type ProfileHTTP2 struct {
	// +kubebuilder:validation:Pattern=`^\/([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	Client string `json:"client,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	Server string `json:"server,omitempty"`
}

type ProfileAdapt struct {
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	Request string `json:"request,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/[a-zA-Z]+([A-z0-9-_+]+\/)+([-A-z0-9_.:]+\/?)*$`
	Response string `json:"response,omitempty"`
}

type ClientSSLParams struct {
	RenegotiationEnabled *bool `json:"renegotiationEnabled,omitempty"`
	// +kubebuilder:validation:Enum=bigip;secret
	ProfileReference string `json:"profileReference,omitempty"`
}

type ServerSSLParams struct {
	RenegotiationEnabled *bool `json:"renegotiationEnabled,omitempty"`
	// +kubebuilder:validation:Enum=bigip;secret
	ProfileReference string `json:"profileReference,omitempty"`
}
