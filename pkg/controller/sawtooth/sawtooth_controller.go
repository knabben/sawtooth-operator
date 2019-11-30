package sawtooth

import (
	"context"
	"fmt"
	"github.com/knabben/sawtooth-operator/pkg/controller/assets"
	"strconv"

	sawtoothv1alpha1 "github.com/knabben/sawtooth-operator/pkg/apis/sawtooth/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_sawtooth")


// Add creates a new Sawtooth Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSawtooth{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sawtooth-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Sawtooth
	err = c.Watch(&source.Kind{Type: &sawtoothv1alpha1.Sawtooth{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Sawtooth
	src := &source.Kind{Type: &corev1.Pod{}}
	pred := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
	}
	err = c.Watch(src, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sawtoothv1alpha1.Sawtooth{},
	}, pred)
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSawtooth implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSawtooth{}

// ReconcileSawtooth reconciles a Sawtooth object
type ReconcileSawtooth struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile creates a new Sawtooth cluster per CR added to the system
func (r *ReconcileSawtooth) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Sawtooth")

	// Fetch the Sawtooth instance
	instance := &sawtoothv1alpha1.Sawtooth{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("No Sawtooth CR found.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	podList, err := r.fetchPodListItems()
	if err != nil {
		return reconcile.Result{}, err
	}

	numberPods := len(podList)
	reqLogger.Info("Number of pods", "numberPods", numberPods, "spec", instance.Spec.Nodes)

	// Update status.Nodes if needed
	if numberPods != instance.Spec.Nodes {
		podName := fmt.Sprintf("%s-pod-%d", instance.Name, numberPods)


		peerArgs := []string{}
		for _, svc := range instance.Status.Services {
			reqLogger.Info("Service", "svc", svc)
			peerArgs = append(peerArgs, "--seeds", svc)
		}

		// CreatePod starts a new pod
		pod := assets.CreatePodSpec(instance, podName, numberPods, peerArgs)

		// Set Sawtooth instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
			reqLogger.Error(err, "Failed to set reference.")
			return reconcile.Result{}, err
		}

		err := r.GetPod(pod)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Pod not found, creating.", "podName", podName)

			err = r.CreatePod(pod)
			if err != nil {
				reqLogger.Error(err, "Failed to create pod.")
				return reconcile.Result{}, err
			}

			// Create service
			service := assets.CreateService(strconv.Itoa(numberPods))
			err := r.client.Create(context.TODO(), service)
			if err != nil {
				reqLogger.Error(err, "Failed to create service.")
				return reconcile.Result{}, err
			}

			if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
				reqLogger.Error(err, "Failed to set reference.")
				return reconcile.Result{}, err
			}

			instance.Status.Services = append(instance.Status.Services, service.Name)
			err = r.updateStatus(instance, numberPods)
			if err != nil {
				reqLogger.Error(err, "Failed to update Memcached status")
				return reconcile.Result{}, err
			}

			return reconcile.Result{Requeue: true}, nil
		}
	}

	err = r.updateStatus(instance, numberPods)
	if err != nil {
		reqLogger.Error(err, "Failed to update Memcached status")
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue but save state
	reqLogger.Info(strconv.Itoa(int(instance.Spec.Nodes)))

	return reconcile.Result{}, nil
}

func (r *ReconcileSawtooth) updateStatus(instance *sawtoothv1alpha1.Sawtooth, number int) error {
	instance.Status.NodeNumber = number

	err := r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReconcileSawtooth) GetPod(pod *corev1.Pod) error {
	found := &corev1.Pod{}
	return r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
}

func (r *ReconcileSawtooth) CreatePod(pod *corev1.Pod) error {
	err := r.client.Create(context.TODO(), pod)
	if err != nil {
		return err
	}

	return nil
}

// fetchPodListItem return the items with the Sawtooth label
func (r *ReconcileSawtooth) fetchPodListItems() ([]corev1.Pod, error) {
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace("default"),
		client.MatchingLabels(assets.GetLabel()),
	}

	if err := r.client.List(context.TODO(), podList, listOpts...); err != nil {
		return nil, err
	}
	return podList.Items, nil
}

func (r *ReconcileSawtooth) verifyNodes(nodes int, status sawtoothv1alpha1.SawtoothStatus) bool {
	return nodes == status.NodeNumber
}