package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ControlPlaneSettings struct {
	MasterCount int `json:"masterCount,omitempty"`
	ETCDCount int `json:"etcdCount"`
	MasterSettings `json:"master,omitempty"`
}

type MasterSettings struct {
	RequiredMemory string `json:"requiredMemory,omitempty"`
	RequiredCPU string `json:"requiredCPU,omitempty"`
	MaxMemory string `json:"maxMemory,omitempty"`
	MaxCPU string `json:"maxCPU,omitempty"`
}

type SecuritySettings struct{
	EncryptionSecretName string `json:"encryptionSecret,omitempty"`
}

// EnvironmentSpec defines the desired state of Environment
type EnvironmentSpec struct {
	ControlPlaneSettings `json:"controlPlane,omitempty"`
	SecuritySettings `json:"security,omitempty"`
}

// EnvironmentStatus defines the observed state of Environment
type EnvironmentStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Environment is the Schema for the environments API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=environments,scope=Namespaced
type Environment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EnvironmentSpec   `json:"spec,omitempty"`
	Status EnvironmentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EnvironmentList contains a list of Environment
type EnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Environment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Environment{}, &EnvironmentList{})
}
