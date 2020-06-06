package master

import (
	"fmt"
	"gitlab.globoi.com/tks/gks/control-plane-operator/pkg/apis/gks/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Master struct{
	environment v1alpha1.Environment
	clusterName string
	namespace string
	apiServer apiServer
}

func NewMaster(environment v1alpha1.Environment, clusterName, namespace,
	advertiseAddress, serviceClusterIpRange string, admissionPlugins []string )Master {

	return Master{
		environment: environment,
		clusterName: clusterName,
		namespace: namespace,
		apiServer: newAPIServer(advertiseAddress,serviceClusterIpRange,admissionPlugins),
	}
}

func (master *Master) BuildPod()corev1.Pod{
	return corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: master.clusterName,
			Labels: master.buildPodLabels(),
		},
		Spec: corev1.PodSpec{
			Volumes: master.buildVolumes(),
			Containers: []corev1.Container{
				master.apiServer.BuildContainer(),
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

	printCfg := func(clusterName, cfgName string)string{
		return fmt.Sprintf("%s-%s",master.clusterName, cfgName)
	}

	return []corev1.Volume{
		master.buildSecretVolume("ca", printCfg(master.clusterName,"ca-certs")),
		master.buildSecretVolume("kubernetes", printCfg(master.clusterName,"kubernetes-certs")),
		master.buildSecretVolume("encryptation", printCfg(master.environment.Name,"encryption-config")),
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
