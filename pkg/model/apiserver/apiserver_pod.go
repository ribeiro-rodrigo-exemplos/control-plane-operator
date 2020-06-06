package apiserver

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const POD_NAME = "kube-apiservers"

func NewAPIServer()*corev1.Pod{
	return &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: POD_NAME,
			Namespace: "",
			Labels: newPodLabels(),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{newAPIServerContainer()},
			Volumes: newPodVolumes(),
		},
	}
}

func newPodLabels()map[string]string{
	return map[string]string{
		"app":POD_NAME,
		"cluster": "",
		"group": "control-plane",
	}
}

func newPodVolumes()[]corev1.Volume{
	return []corev1.Volume{
		newVolume("ca", "gks1-kubernetes-certs"),
		newVolume("kubernetes", "gks1-kubernetes-certs"),
		newVolume("encryptation", "encryption-config"),
	}
}

func newVolume(volumeName, secretName string)corev1.Volume{
	return corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
			},
		},
	}
}