package master

import (
	"gitlab.globoi.com/tks/gks/control-plane-operator/pkg/apis/gks/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Master struct{
	environment v1alpha1.Environment
	clusterName string
	namespace string
	masterSecretName string
	apiServer apiServer
	scheduler Scheduler
	controllerManager ControllerManager
}

func NewMaster(environment v1alpha1.Environment, clusterName, namespace,
	advertiseAddress, serviceClusterIpRange, clusterCIDRS, masterSecretName string,
	admissionPlugins []string )Master {

	return Master{
		environment: environment,
		clusterName: clusterName,
		namespace: namespace,
		masterSecretName: masterSecretName,
		apiServer: newAPIServer(advertiseAddress,serviceClusterIpRange,admissionPlugins,environment.Spec.MasterCount),
		scheduler: NewScheduler(),
		controllerManager: NewControllerManager(clusterName, serviceClusterIpRange,clusterCIDRS),
	}
}

func (master *Master) BuildPod()corev1.Pod{
	return corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Namespace: master.environment.Namespace,
			Name: master.clusterName,
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
		"cluster": master.clusterName,
		"group": "control-plane",
	}
}

func (master *Master) buildVolumes()[]corev1.Volume{

	return []corev1.Volume{
		master.buildSecretVolume("ca", "ca-certs"),
		master.buildSecretVolume("kubernetes", master.masterSecretName),
		master.buildSecretVolume("encryption", master.environment.Spec.EncryptionSecretName),
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
