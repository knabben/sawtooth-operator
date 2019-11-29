package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SawtoothSpec defines the desired state of Sawtooth
// +k8s:openapi-gen=true
type SawtoothSpec struct {
	Nodes     int64  `json:"nodes"`
	Version   string `json:"version"`
	Consensus string `json:"consensus"`
}

// SawtoothStatus defines the observed state of Sawtooth
// +k8s:openapi-gen=true
type SawtoothStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Sawtooth is the Schema for the sawtooths API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sawtooths,scope=Namespaced
type Sawtooth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SawtoothSpec   `json:"spec,omitempty"`
	Status SawtoothStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SawtoothList contains a list of Sawtooth
type SawtoothList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Sawtooth `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Sawtooth{}, &SawtoothList{})
}
