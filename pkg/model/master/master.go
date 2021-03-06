package master

import (
	"fmt"
	"gitlab.globoi.com/tks/gks/control-plane-operator/pkg/apis/gks/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Master struct{
	settings v1alpha1.MasterSettings
	namespacedName types.NamespacedName
	apiServer apiServer
	scheduler Scheduler
	controllerManager ControllerManager
}

func NewMaster(namespacedName types.NamespacedName, settings v1alpha1.MasterSettings)Master {

	advertiseAddress := "192.168.39.42"

	return Master{
		settings: settings,
		namespacedName: namespacedName,
		apiServer: newAPIServer(
			advertiseAddress,
			settings.ServiceClusterIPRange,
			settings.AdmissionPlugins,
		),
		scheduler: NewScheduler(),
		controllerManager: NewControllerManager(namespacedName.Name, settings.ServiceClusterIPRange,settings.ClusterCIDR),
	}
}

func (master *Master) Merge(newMaster Master)Merger{
	return Merger{
		oldMaster: *master,
		newMaster: newMaster,
		namespacedName: master.namespacedName,
	}
}

func (master *Master) BuildDeployment()*appsv1.Deployment{

	replicas := int32(1)

	return &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Namespace: master.namespacedName.Namespace,
			Name: fmt.Sprintf("cluster-%s",master.namespacedName.Name),
			Labels: master.buildPodLabels(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: master.buildPodLabels(),
			},
			Template: master.buildPod(),
		},
	}
}

func (master *Master) buildPod()corev1.PodTemplateSpec{
	return corev1.PodTemplateSpec{
		ObjectMeta: v1.ObjectMeta{
			Namespace: master.namespacedName.Namespace,
			Labels: master.buildPodLabels(),
		},
		Spec: corev1.PodSpec{
			Volumes: master.buildVolumes(),
			Containers: []corev1.Container{
				master.apiServer.BuildContainer(),
				master.scheduler.BuilderContainer(),
				master.controllerManager.BuilderContainer(),
			},
		},
	}
}

func (master *Master) buildPodLabels()map[string]string{
	return map[string]string{
		"app":"master",
		"cluster": master.namespacedName.Name,
		"tier": "control-plane",
	}
}

func (master *Master) buildVolumes()[]corev1.Volume{

	return []corev1.Volume{
		master.buildSecretVolume("ca", "ca-certs"),
		master.buildSecretVolume("kubernetes", master.settings.MasterSecretName),
		master.buildSecretVolume("encryption", master.settings.EncryptionSecretName),
	}
}

func (*Master) buildSecretVolume(volumeName, secretName string)corev1.Volume{
	return corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
			},
		},
	}
}
