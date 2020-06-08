package controlplane

import (
	"context"
	gksv1alpha1 "gitlab.globoi.com/tks/gks/control-plane-operator/pkg/apis/gks/v1alpha1"
	"gitlab.globoi.com/tks/gks/control-plane-operator/pkg/model/master"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

var log = logf.Log.WithName("controller_controlplane")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileControlPlane{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {

	c, err := controller.New("controlplane-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &gksv1alpha1.ControlPlane{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &gksv1alpha1.ControlPlane{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileControlPlane{}

type ReconcileControlPlane struct {
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileControlPlane) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ControlPlane")

	instance := &gksv1alpha1.ControlPlane{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err){
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	environment := &gksv1alpha1.Environment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: instance.Namespace,
		Name: instance.Spec.EnvironmentName,
	}, environment)

	if err != nil {
		if errors.IsNotFound(err){
			reqLogger.Info("Environment not found")
			//create default environment settings
		}
		return reconcile.Result{}, err
	}

	masterModel := buildMaster(instance, environment)

	masterDeployment := masterModel.BuildDeployment()

	if err := controllerutil.SetControllerReference(instance, masterDeployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	err = r.client.Create(context.TODO(), masterDeployment)

	if err != nil {
		reqLogger.Error(err,"Error in create pod master")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func buildMaster(instance *gksv1alpha1.ControlPlane, environment *gksv1alpha1.Environment)master.Master{
	return master.NewMaster(
		*environment,
		instance.Name,
		instance.Namespace,
		"192.168.39.42",
		instance.Spec.ServiceClusterIPRange,
		instance.Spec.ClusterCIDR,
		instance.Spec.MasterSecretName,
		instance.Spec.AdmissionPlugins,
	)
}
