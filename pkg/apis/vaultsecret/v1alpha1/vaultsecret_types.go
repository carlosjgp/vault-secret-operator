package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VaultSecretSpec defines the desired state of VaultSecret
// +k8s:openapi-gen=true
type VaultSecretSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ServiceAccount string             `json:"serviceAccount,omitempty"`
	VaultAddress   string             `json:"vaultAddress"`
	VaultCA        string             `json:"vaultCA,omitempty"`
	VaultAgent     VaultAgentSpec     `json:"vaultAgent"`
	ConsulTemplate ConsulTemplateSpec `json:"consulTemplate"`
	Secret         SecretSpec         `json:"secret"`
	// InitContainers to be used on the POD
	// +listType=set
	InitContainers []corev1.Container `json:"initContainers,omitempty"`
	// ExtraContainers to be used on the POD along side the main containers
	// +listType=set
	ExtraContainers []corev1.Container `json:"extraContainers,omitempty"`
}

// VaultSecretStatus defines the observed state of VaultSecret
// +k8s:openapi-gen=true
type VaultSecretStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VaultAgentSpec is the Schema for the vaultsecrets API
type VaultAgentSpec struct {
	Image    ContainerImageSpec `json:"image,omitempty"`
	AutoAuth string             `json:"autoAuth,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretSpec is the Schema for the vaultsecrets API
type SecretSpec struct {
	Name string          `json:"name,omitempty"`
	Keys []SecretKeySpec `json:"keys"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretKeySpec is the Schema for the vaultsecrets API
type SecretKeySpec struct {
	File string `json:"file"`
	Key  string `json:"key"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConsulTemplateSpec is the Schema for the vaultsecrets API
type ConsulTemplateSpec struct {
	Image     ContainerImageSpec `json:"image,omitempty"`
	Templates map[string]string  `json:"templates"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContainerImageSpec is the Schema for the vaultsecrets API
type ContainerImageSpec struct {
	Repository      string `json:"repository,omitempty"`
	Tag             string `json:"tag,omitempty"`
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VaultSecret is the Schema for the vaultsecrets API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=vaultsecrets,scope=Namespaced
type VaultSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VaultSecretSpec   `json:"spec,omitempty"`
	Status VaultSecretStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VaultSecretList contains a list of VaultSecret
type VaultSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VaultSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VaultSecret{}, &VaultSecretList{})
}
