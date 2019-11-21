package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VaultSecretSpec defines the desired state of VaultSecret
// +k8s:openapi-gen=true
type VaultSecretSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ServiceAccount string             `json:"serviceAccount,omitempty" protobuf:"bytes,1,opt,name=serviceAccount"`
	VaultAddress   string             `json:"vaultAddress" protobuf:"bytes,2,name=vaultAddress"`
	VaultCA        string             `json:"vaultCA,omitempty" protobuf:"bytes,3,opt,name=vaultCA"`
	VaultAgent     VaultAgentSpec     `json:"vaultAgent,omitempty" protobuf:"bytes,4,name=vaultAgent"`
	ConsulTemplate ConsulTemplateSpec `json:"consulTemplate,omitempty" protobuf:"bytes,5,name=consulTemplate"`
	KubectlVersion string             `json:"kubectlVersion,omitempty" protobuf:"bytes,5,name=kubectlVersion"`
	Secret         SecretSpec         `json:"secret" protobuf:"bytes,6,opt,name=secret"`
	// InitContainers to be used on the POD
	// +listType=set
	InitContainers []corev1.Container `json:"initContainers,omitempty" protobuf:"bytes,7,rep,name=initContainers"`
	// ExtraContainers to be used on the POD along side the main containers
	// +listType=set
	ExtraContainers []corev1.Container `json:"extraContainers,omitempty" protobuf:"bytes,8,rep,name=extraContainers"`
}

// VaultSecretStatus defines the observed state of VaultSecret
// +k8s:openapi-gen=true
type VaultSecretStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// VaultAgentSpec is the Schema for the vaultsecrets API
// +k8s:openapi-gen=true
type VaultAgentSpec struct {
	Image ContainerImageSpec `json:"image,omitempty" protobuf:"bytes,1,name=image"`
	// Entrypoint array. Not executed within a shell.
	// The docker image's ENTRYPOINT is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax
	// can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	// +listType=set
	Command        []string `json:"command,omitempty" protobuf:"bytes,2,rep,name=command"`
	AutoAuthMethod string   `json:"autoAuthMethod,omitempty" protobuf:"bytes,3,name=autoAuthMethod"`
}

// SecretSpec is the Schema for the vaultsecrets API
// +k8s:openapi-gen=true
type SecretSpec struct {
	Name string `json:"name,omitempty" protobuf:"bytes,1,name=name"`
	Path string `json:"path,omitempty" protobuf:"bytes,2,rep,name=path"`
}

// ConsulTemplateSpec is the Schema for the vaultsecrets API
// +k8s:openapi-gen=true
type ConsulTemplateSpec struct {
	Image ContainerImageSpec `json:"image,omitempty" protobuf:"bytes,1,name=image"`
	// Entrypoint array. Not executed within a shell.
	// The docker image's ENTRYPOINT is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax
	// can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	// +listType=set
	Command   []string `json:"command,omitempty" protobuf:"bytes,2,rep,name=command"`
	Templates string   `json:"templates,omitempty" protobuf:"bytes,3,name=templates"`
}

// ContainerImageSpec is the Schema for the vaultsecrets API
// +k8s:openapi-gen=true
type ContainerImageSpec struct {
	Repository      string            `json:"repository,omitempty" protobuf:"bytes,1,name=repository"`
	Tag             string            `json:"tag,omitempty" protobuf:"bytes,2,name=tag"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty" protobuf:"bytes,3,name=imagePullPolicy"`
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
	SchemeBuilder.Register(
		&VaultSecret{},
		&VaultSecretList{})
}
