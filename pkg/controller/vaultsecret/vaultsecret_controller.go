package vaultsecret

import (
	"bytes"
	"context"
	"html/template"

	"github.com/Masterminds/sprig"
	vaultsecretv1alpha1 "github.com/carlosjgp/vault-secret-operator/pkg/apis/vaultsecret/v1alpha1"
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

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

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
	instance := &vaultsecretv1alpha1.VaultSecret{
		Spec: vaultsecretv1alpha1.VaultSecretSpec{
			VaultAgent: vaultsecretv1alpha1.VaultAgentSpec{
				Image: vaultsecretv1alpha1.ContainerImageSpec{
					Repository:      "vault",
					Tag:             "latest",
					ImagePullPolicy: "Always",
				},
			},
			ConsulTemplate: vaultsecretv1alpha1.ConsulTemplateSpec{
				Image: vaultsecretv1alpha1.ContainerImageSpec{
					Repository:      "hashicorp/consul-template",
					Tag:             "latest",
					ImagePullPolicy: "Always",
				},
			},
		},
	}
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
	allResources := []metav1.Object{}
	consulTemplatesCM := newConsulTemplatesConfigMapForCR(instance)
	vaultAgentCM := newVaultAgentConfigMapForCR(instance)
	pod := newPodForCR(instance, consulTemplatesCM, vaultAgentCM)

	//TODO create serviceaccount

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
func newPodForCR(cr *vaultsecretv1alpha1.VaultSecret, consulTemplates *corev1.ConfigMap, vaultAgent *corev1.ConfigMap) *corev1.Pod {
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
					Image:           cr.Spec.VaultAgent.Image.Repository + ":" + cr.Spec.VaultAgent.Image.Tag,
					ImagePullPolicy: cr.Spec.VaultAgent.Image.ImagePullPolicy,
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
					},
				},
				corev1.Container{
					Name:            "consul-template",
					Image:           cr.Spec.ConsulTemplate.Image.Repository + ":" + cr.Spec.ConsulTemplate.Image.Tag,
					ImagePullPolicy: cr.Spec.ConsulTemplate.Image.ImagePullPolicy,
					Command: []string{
						"consul-template",
						"-config",
						"/etc/consul-template/config.hcl",
					},
					VolumeMounts: []corev1.VolumeMount{
						corev1.VolumeMount{
							Name:      "consul-template",
							ReadOnly:  true,
							MountPath: "/etc/consul-template",
						},
					},
				},
			),

			Volumes: []corev1.Volume{
				corev1.Volume{
					Name: "vault-agent",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: vaultAgent.Name,
							},
						},
					},
				},
				corev1.Volume{
					Name: "consul-template",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: consulTemplates.Name,
							},
						},
					},
				},
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
				"/templates/consul-template.conf.hcl",
				"consul-template",
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
				"/templates/vault-agent.conf.hcl",
				"vault-agent",
				map[string]interface{}{
					"AutoAuthMethod": cr.Spec.VaultAgent.AutoAuthMethod,
				}),
		},
	}
}

func templateFile(path string, tempateName string, data interface{}) string {
	t := template.Must(
		template.New(tempateName).
			Funcs(sprig.FuncMap()).
			ParseFiles(path))

	var templateBuffer bytes.Buffer
	if err := t.Execute(&templateBuffer, data); err != nil {
		// TODO
		//return err
	}

	return templateBuffer.String()
}
