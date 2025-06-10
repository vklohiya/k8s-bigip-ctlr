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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TLSProfileSpec defines the desired state of TLSProfile.
type TLSProfileSpec struct {
	// +kubebuilder:validation:Pattern=`^(([a-zA-Z0-9\*]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`
	Hosts     []string         `json:"hosts"`
	TLS       TLS              `json:"tls"`
	TLSCipher TLSProfileCipher `json:"tlsCipher"`
}

type TLSProfileCipher struct {
	// +kubebuilder:validation:Enum="1.0";"1.1";"1.2";"1.3"
	TLSVersion  string `json:"tlsVersion"`
	Ciphers     string `json:"ciphers"`
	CipherGroup string `json:"cipherGroup"`
	// +kubebuilder:validation:Enum="1.0";"1.1";"1.2";"1.3"
	DisableTLSVersions []string `json:"disableTLSVersions"`
}

type TLS struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=edge;reencrypt;passthrough
	Termination string `json:"termination"`
	// +kubebuilder:validation:Pattern=`^\/?[a-zA-Z]+([-A-z0-9_+]+\/)*([-A-z0-9_.:]+\/?)*$`
	ClientSSL string `json:"clientSSL,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/?[a-zA-Z]+([-A-z0-9_+]+\/)*([-A-z0-9_.:]+\/?)*$`
	ClientSSLs []string `json:"clientSSLs,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/?[a-zA-Z]+([-A-z0-9_+]+\/)*([-A-z0-9_.:]+\/?)*$`
	ServerSSL string `json:"serverSSL,omitempty"`
	// +kubebuilder:validation:Pattern=`^\/?[a-zA-Z]+([-A-z0-9_+]+\/)*([-A-z0-9_.:]+\/?)*$`
	ServerSSLs []string `json:"serverSSLs,omitempty"`
	// +kubebuilder:validation:Enum=bigip;secret;hybrid
	Reference       string          `json:"reference,omitempty"`
	ClientSSLParams ClientSSLParams `json:"clientSSLParams,omitempty"`
	ServerSSLParams ServerSSLParams `json:"serverSSLParams,omitempty"`
}

// TLSProfileStatus defines the observed state of TLSProfile.
type TLSProfileStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TLSProfile is the Schema for the tlsprofiles API.
type TLSProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TLSProfileSpec   `json:"spec,omitempty"`
	Status TLSProfileStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TLSProfileList contains a list of TLSProfile.
type TLSProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TLSProfile `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TLSProfile{}, &TLSProfileList{})
}
