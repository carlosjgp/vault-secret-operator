package vaultsecret

import (
	"bytes"
	"context"
	"fmt"

	"text/template"

	"github.com/Masterminds/sprig"
	vaultsecretv1alpha1 "github.com/carlosjgp/vault-secret-operator/pkg/apis/vaultsecret/v1alpha1"
	"github.com/carlosjgp/vault-secret-operator/version"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_vaultsecret")

// Add creates a new VaultSecret Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileVaultSecret{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("vaultsecret-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource VaultSecret
	err = c.Watch(&source.Kind{Type: &vaultsecretv1alpha1.VaultSecret{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner VaultSecret
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &vaultsecretv1alpha1.VaultSecret{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileVaultSecret implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileVaultSecret{}

// ReconcileVaultSecret reconciles a VaultSecret object
type ReconcileVaultSecret struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a VaultSecret object and makes changes based on the state read
// and what is in the VaultSecret.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileVaultSecret) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling VaultSecret")

	// Fetch the VaultSecret instance with defaults
	instance := &vaultsecretv1alpha1.VaultSecret{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	consulTemplatesCM := newConsulTemplatesConfigMapForCR(instance)
	vaultAgentCM := newVaultAgentConfigMapForCR(instance)
	pod := newPodForCR(instance, consulTemplatesCM, vaultAgentCM)
	allConfigMaps := []*corev1.ConfigMap{
		consulTemplatesCM,
		vaultAgentCM,
	}

	//TODO create serviceaccount?

	allResources := []metav1.Object{}

	allResources = append(
		allResources,
		pod,
		consulTemplatesCM,
		vaultAgentCM,
	)

	for _, k8sResource := range allResources {
		// Set VaultSecret instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, k8sResource, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
	}

	// Check if ConfigMaps already exists
	for _, cm := range allConfigMaps {
		found := &corev1.ConfigMap{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
			err = r.client.Create(context.TODO(), cm)
			if err != nil {
				return reconcile.Result{}, err
			}

			// ConfigMap created successfully - don't requeue
		} else if err != nil {
			return reconcile.Result{}, err
		}

		// Pod already exists - don't requeue
		reqLogger.Info("Skip reconcile: ConfigMap already exists", "PoConfigMap.Namespace", found.Namespace, "ConfigMap.Name", found.Name)
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

func resourceLabels(cr *vaultsecretv1alpha1.VaultSecret) map[string]string {
	return map[string]string{
		"app": cr.Name,
	}
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(
	cr *vaultsecretv1alpha1.VaultSecret,
	consulTemplates *corev1.ConfigMap,
	vaultAgent *corev1.ConfigMap) *corev1.Pod {

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    resourceLabels(cr),
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: cr.Spec.ServiceAccount,
			InitContainers:     cr.Spec.InitContainers,

			Containers: append(
				cr.Spec.ExtraContainers,
				corev1.Container{
					Name:            "vault-agent",
					Image:           getVaultAgentImage(&cr.Spec),
					ImagePullPolicy: getVaultAgentImagePullPolicy(&cr.Spec),
					Command: []string{
						"vault",
						"agent",
						"-config=/etc/vault/config.hcl",
					},
					VolumeMounts: []corev1.VolumeMount{
						corev1.VolumeMount{
							Name:      "vault-agent",
							ReadOnly:  true,
							MountPath: "/etc/vault",
						},
						corev1.VolumeMount{
							Name:      "vault-token",
							MountPath: "/tmp/vault/agent",
						},
					},
				},
				corev1.Container{
					Name:            "consul-template",
					Image:           getConsulTemplateImage(&cr.Spec),
					ImagePullPolicy: getConsulTemplateImagePullPolicy(&cr.Spec),
					Command: []string{
						"/consul-template",
						"-config",
						"/etc/consul-template/config.hcl",
					},
					VolumeMounts: []corev1.VolumeMount{
						corev1.VolumeMount{
							Name:      "consul-template",
							ReadOnly:  true,
							MountPath: "/etc/consul-template",
						},
						corev1.VolumeMount{
							Name:      "vault-token",
							ReadOnly:  true,
							MountPath: "/tmp/vault/agent",
						},
						corev1.VolumeMount{
							Name:      "templated-secrets",
							ReadOnly:  false,
							MountPath: getTemplatedSecretsMountPath(&cr.Spec),
						},
					},
				},
				corev1.Container{
					Name:  "kubectl",
					Image: fmt.Sprintf("%s:%s", "carlosjgp/vault-secret-operator-kubectl", getKubectlVersion(&cr.Spec)),
					Env: []corev1.EnvVar{
						corev1.EnvVar{
							Name:  "SECRET",
							Value: cr.Spec.Secret.Name,
						},
						corev1.EnvVar{
							Name:  "FOLDER",
							Value: getTemplatedSecretsMountPath(&cr.Spec),
						},
						corev1.EnvVar{
							Name: "NAMESPACE",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									FieldPath: "metadata.namespace",
								},
							},
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						corev1.VolumeMount{
							Name:      "templated-secrets",
							ReadOnly:  true,
							MountPath: getTemplatedSecretsMountPath(&cr.Spec),
						},
					},
				},
			),

			Volumes: []corev1.Volume{
				volumeFromConfigMap("vault-agent", vaultAgent.Name),
				volumeFromConfigMap("consul-template", consulTemplates.Name),
				newEmptyDirInMemory("vault-token"),
				newEmptyDirInMemory("templated-secrets"),
			},
		},
	}
}

func volumeFromConfigMap(name string, configMapName string) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configMapName,
				},
			},
		},
	}
}

func newEmptyDirInMemory(name string) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium: corev1.StorageMediumMemory,
			},
		},
	}
}

func newConsulTemplatesConfigMapForCR(cr *vaultsecretv1alpha1.VaultSecret) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-consul-templates",
			Namespace: cr.Namespace,
			Labels:    resourceLabels(cr),
		},
		Data: map[string]string{
			"config.hcl": templateFile(
				"consul-template.conf.hcl",
				map[string]interface{}{
					"ConsulTemplates": cr.Spec.ConsulTemplate.Templates,
				}),
		},
	}
}

func newVaultAgentConfigMapForCR(cr *vaultsecretv1alpha1.VaultSecret) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-vault-agent",
			Namespace: cr.Namespace,
			Labels:    resourceLabels(cr),
		},
		Data: map[string]string{
			"config.hcl": templateFile(
				"vault-agent.conf.hcl",
				map[string]interface{}{
					"VaultAddress":   cr.Spec.VaultAddress,
					"AutoAuthMethod": cr.Spec.VaultAgent.AutoAuthMethod,
				}),
		},
	}
}

func templateFile(tempate string, data interface{}) string {
	t := template.Must(
		template.New(tempate).
			Funcs(sprig.TxtFuncMap()).
			ParseGlob("/templates/*"))

	var templateBuffer bytes.Buffer
	if err := t.Execute(&templateBuffer, data); err != nil {
		// TODO
		panic(err)
	}

	return templateBuffer.String()
}

func getVaultAgentImage(vs *vaultsecretv1alpha1.VaultSecretSpec) string {
	repo := "vault"
	tag := "latest"

	if vs.VaultAgent.Image.Repository != "" {
		repo = vs.VaultAgent.Image.Repository
	}
	if vs.VaultAgent.Image.Tag != "" {
		tag = vs.VaultAgent.Image.Tag
	}
	return fmt.Sprintf("%s:%s", repo, tag)
}

func getVaultAgentImagePullPolicy(vs *vaultsecretv1alpha1.VaultSecretSpec) corev1.PullPolicy {
	policy := corev1.PullAlways
	policy = vs.VaultAgent.Image.ImagePullPolicy
	return policy
}

func getConsulTemplateImage(vs *vaultsecretv1alpha1.VaultSecretSpec) string {
	repo := "hashicorp/consul-template"
	tag := "latest"

	if vs.ConsulTemplate.Image.Repository != "" {
		repo = vs.ConsulTemplate.Image.Repository
	}
	if vs.ConsulTemplate.Image.Tag != "" {
		tag = vs.ConsulTemplate.Image.Tag
	}
	return fmt.Sprintf("%s:%s", repo, tag)
}

func getKubectlVersion(vs *vaultsecretv1alpha1.VaultSecretSpec) string {
	tag := "v1.16.0"

	if vs.KubectlVersion != "" {
		tag = vs.KubectlVersion
	}
	return fmt.Sprintf("%s-v%s", tag, version.Version)
}

func getConsulTemplateImagePullPolicy(vs *vaultsecretv1alpha1.VaultSecretSpec) corev1.PullPolicy {
	policy := corev1.PullAlways
	policy = vs.ConsulTemplate.Image.ImagePullPolicy
	return policy
}

func getTemplatedSecretsMountPath(vs *vaultsecretv1alpha1.VaultSecretSpec) string {
	mountPath := "/tmp/templated"
	if vs.Secret.Path != "" {
		mountPath = vs.Secret.Path
	}
	return mountPath
}
