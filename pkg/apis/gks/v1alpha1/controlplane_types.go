package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ETCDSettings struct {
	InstancesCount int `json:"instancesCount,omitempty"`
}

type MasterResources struct {
	RequiredMemory string `json:"requiredMemory,omitempty"`
	RequiredCPU string `json:"requiredCPU,omitempty"`
	LimitMemory string `json:"limitMemory,omitempty"`
	LimitCPU string `json:"limitCPU,omitempty"`
}

type MasterScaleSettings struct {
	MaxInstances int `json:"maxInstances,omitempty"`
	MinInstances int `json:"minInstances,omitempty"`
	LimitCPUPercent int `json:"limitCPUPercent,omitempty"`
	LimitMemoryPercent string `json:"limitMemoryPercent,omitempty"`
}

type MasterSettings struct {
	MasterClusterSettings `json:"settings,omitempty"`
	MasterScaleSettings   `json:"scale,omitempty"`
}

type MasterClusterSettings struct {
	MasterSecretName string `json:"certsSecret,omitempty"`
	MasterResources string `json:"resources,omitempty"`
	AdmissionPlugins []string `json:"admissionPlugins,omitempty"`
	ServiceClusterIPRange string `json:"serviceClusterIpRange,omitempty"`
	ClusterCIDR string `json:"clusterCidr,omitempty"`
	EncryptionSecretName string `json:"encryptionSecret,omitempty"`
}

// ControlPlaneSpec defines the desired state of ControlPlane
type ControlPlaneSpec struct {
	EnvironmentName string `json:"environment,omitempty"`
	MasterSettings `json:"master,omitempty"`
}

// ControlPlaneStatus defines the observed state of ControlPlane
type ControlPlaneStatus struct {
	LastMasterSettings *MasterSettings `json:"master,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControlPlane is the Schema for the controlplanes API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=controlplanes,scope=Namespaced
type ControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ControlPlaneSpec   `json:"spec,omitempty"`
	Status ControlPlaneStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControlPlaneList contains a list of ControlPlane
type ControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ControlPlane `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ControlPlane{}, &ControlPlaneList{})
}
