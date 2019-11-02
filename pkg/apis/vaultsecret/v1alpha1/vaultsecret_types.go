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

// VaultAgentSpec is the Schema for the vaultsecrets API
type VaultAgentSpec struct {
	Image ContainerImageSpec `json:"image,omitempty"`
	// Entrypoint array. Not executed within a shell.
	// The docker image's ENTRYPOINT is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax
	// can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	Command        []string `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	AutoAuthMethod string   `json:"autoAuthMethod"`
}

// SecretSpec is the Schema for the vaultsecrets API
type SecretSpec struct {
	Name string          `json:"name,omitempty"`
	Keys []SecretKeySpec `json:"keys"`
}

// SecretKeySpec is the Schema for the vaultsecrets API
type SecretKeySpec struct {
	File string `json:"file"`
	Key  string `json:"key"`
}

// ConsulTemplateSpec is the Schema for the vaultsecrets API
type ConsulTemplateSpec struct {
	Image     ContainerImageSpec `json:"image,omitempty"`
	Templates string             `json:"templates"`
}

// ContainerImageSpec is the Schema for the vaultsecrets API
type ContainerImageSpec struct {
	Repository      string            `json:"repository,omitempty"`
	Tag             string            `json:"tag,omitempty"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
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
