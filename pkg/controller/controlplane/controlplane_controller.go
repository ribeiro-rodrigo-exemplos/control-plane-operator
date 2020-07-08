package controlplane

import (
	"context"
	"fmt"
	gksv1alpha1 "gitlab.globoi.com/tks/gks/control-plane-operator/pkg/apis/gks/v1alpha1"
	autoscalingv2 "k8s.io/api/autoscaling/v2beta2"
	"gitlab.globoi.com/tks/gks/control-plane-operator/pkg/model/master"
	appsv1 "k8s.io/api/apps/v1"
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

	clusterName := fmt.Sprintf("cluster-%s",instance.Name)
	clusterNamespacedName := types.NamespacedName{Name: clusterName, Namespace: request.Namespace}

	masterDeployment := &appsv1.Deployment{}

	err = r.client.Get(context.TODO(), clusterNamespacedName ,masterDeployment)

	if err != nil {
		if errors.IsNotFound(err){
			err = r.createMaster(request.NamespacedName, instance)
			if err != nil {
				return reconcile.Result{}, err
			}else{
				return reconcile.Result{Requeue: true}, nil
			}
		}
		return reconcile.Result{}, err
	}

	masterHPA := &autoscalingv2.HorizontalPodAutoscaler{}
	err = r.client.Get(context.TODO(), clusterNamespacedName, masterHPA)

	if err != nil {
		if errors.IsNotFound(err){
			if err = r.createHPA(request.NamespacedName, instance); err != nil {
				return reconcile.Result{}, err
			}else{
				return reconcile.Result{}, nil
			}
		}
		return reconcile.Result{}, err
	}

	if instance.Status.LastMasterSettings != nil {

		masterMerged, mergedSettings, mergedScaleSettings := r.ensureMerge(request.NamespacedName,instance)

		if mergedSettings {
			updateDeploy := masterMerged.BuildDeployment()
			if err = r.client.Update(context.TODO(),updateDeploy); err != nil {
				return reconcile.Result{}, err
			}
		}

		if mergedScaleSettings {
			updateHPA := masterMerged.BuildAutoScaling()
			if err = r.client.Update(context.TODO(), updateHPA); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	 if err := r.updateMasterStatus(instance); err != nil {
	 	return reconcile.Result{}, err
	 }

	return reconcile.Result{}, nil
}

func (r *ReconcileControlPlane) ensureMerge(namespacedName types.NamespacedName,
	instance *gksv1alpha1.ControlPlane) (master.Master, bool, bool){

	currentMaster := master.NewMaster(namespacedName, *instance.Status.LastMasterSettings)

	merger := currentMaster.Merge(instance.Spec.MasterSettings)

	return merger.MergeSettings()
}

func (r *ReconcileControlPlane) createMaster(namspacedName types.NamespacedName, instance *gksv1alpha1.ControlPlane)error{
	masterModel := master.NewMaster(namspacedName, instance.Spec.MasterSettings)

	masterDeployment := masterModel.BuildDeployment()

	if err := controllerutil.SetControllerReference(instance, masterDeployment, r.scheme); err != nil {
		return err
	}

	if err := r.client.Create(context.TODO(), masterDeployment); err != nil {
		return err
	}

	return nil
}

func (r *ReconcileControlPlane) createHPA(namespacedName types.NamespacedName, instance *gksv1alpha1.ControlPlane)error{
	return nil
}

func (r *ReconcileControlPlane) updateMasterStatus(instance *gksv1alpha1.ControlPlane)error {
	instance.Status.LastMasterSettings = &instance.Spec.MasterSettings
	err := r.client.Status().Update(context.TODO(), instance)
	return err
}

