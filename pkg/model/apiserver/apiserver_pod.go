package apiserver

import (
	"fmt"
	"gitlab.globoi.com/tks/gks/control-plane-operator/pkg/apis/gks/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type APIServer struct {
	clusterName string
	namespace string
	applicationName string
	image string
	advertiseAddress string
	serviceClusterIPRange string
	admissionPlugins []string
	loadBalancerAddress string
	environment v1alpha1.Environment
}

func NewAPIServer(environment v1alpha1.Environment, clusterName, namespace,
	advertiseAddress, serviceClusterIpRange, loadBalancerAddress string,
	admissionPlugins []string )APIServer{
	apiServer := APIServer{
		environment: environment,
		clusterName: clusterName,
		namespace: namespace,
		applicationName: "kube-apiserver",
		image: "rodrigoribeiro/globo-kube-apiserver",
		advertiseAddress: advertiseAddress,
		serviceClusterIPRange: serviceClusterIpRange,
		loadBalancerAddress: loadBalancerAddress,
		admissionPlugins: admissionPlugins,
	}

	return apiServer
}

func (apiServer *APIServer) BuildPod()*corev1.Pod{
	return &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: apiServer.applicationName,
			Namespace: apiServer.namespace,
			Labels: apiServer.buildPodLabels(),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{apiServer.buildContainer()},
			Volumes: apiServer.buildPodVolumes(),
		},
	}
}

func (apiServer *APIServer) buildPodLabels()map[string]string{
	return map[string]string{
		"app":apiServer.applicationName,
		"cluster": apiServer.clusterName,
		"group": "control-plane",
	}
}

func (apiServer *APIServer) buildPodVolumes()[]corev1.Volume{

	printCfg := func(clusterName, cfgName string)string{
		return fmt.Sprintf("%s-%s",apiServer.clusterName, cfgName)
	}

	return []corev1.Volume{
		apiServer.buildSecretVolume("ca", printCfg(apiServer.clusterName,"ca-certs")),
		apiServer.buildSecretVolume("kubernetes", printCfg(apiServer.clusterName,"kubernetes-certs")),
		apiServer.buildSecretVolume("encryptation", printCfg(apiServer.environment.Name,"encryption-config")),
	}
}

func (apiServer *APIServer) buildSecretVolume(volumeName, secretName string)corev1.Volume{
	return corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
			},
		},
	}
}

func (apiServer *APIServer) buildContainer()corev1.Container{
	return corev1.Container{
		Name: apiServer.applicationName,
		Image: apiServer.image,
		VolumeMounts: apiServer.buildVolumeMounts(),
		Command: apiServer.buildCommands(),
	}
}

func (apiServer *APIServer) buildVolumeMounts()[]corev1.VolumeMount{
	return []corev1.VolumeMount{
		{Name: "ca", MountPath: "/var/lib/kubernetes/ca", ReadOnly: true},
		{Name: "kubernetes", MountPath: "/var/lib/kubernetes", ReadOnly: true},
		{Name: "encryption", MountPath: "/var/lib/kubernetes/encryption", ReadOnly: true},
	}
}

func (apiServer *APIServer) buildCommands()[]string{

	printFlag := func(flag, value interface{})string{
		return fmt.Sprintf("--%s=%v",flag, value)
	}

	printAdmissionPlugins := func ()string{
		return strings.Join(apiServer.admissionPlugins,",")
	}

	return []string{
		apiServer.applicationName,
		printFlag("advertise-address",apiServer.advertiseAddress),
		printFlag("allow-privileged",true),
		printFlag("apiserver-count", 1),
		printFlag("audit-log-maxage",30),
		printFlag("audit-log-maxbackup",3),
		printFlag("audit-log-maxsize",100),
		printFlag("audit-log-path=","/var/log/audit.log"),
		printFlag("authorization-mode","Node,RBAC"),
		printFlag("bind-address","0.0.0.0"),
		printFlag("client-ca-file","/var/lib/kubernetes/ca/ca.pem"),
		printFlag("enable-admission-plugins", printAdmissionPlugins()),
		printFlag("etcd-cafile","/var/lib/kubernetes/ca/ca.pem"),
		printFlag("etcd-certfile","/var/lib/kubernetes/kubernetes.pem"),
		printFlag("etcd-keyfile","/var/lib/kubernetes/kubernetes-key.pem"),
		printFlag("etcd-servers","https://161.35.116.213:2379"),
		printFlag("event-ttl","1h"),
		printFlag("encryption-provider-config","/var/lib/kubernetes/encryptation/encryption-config.yaml"),
		printFlag("kubelet-certificate-authority","/var/lib/kubernetes/ca/ca.pem"),
		printFlag("kubelet-client-certificate","/var/lib/kubernetes/kubernetes.pem"),
		printFlag("kubelet-client-key","/var/lib/kubernetes/kubernetes-key.pem"),
		printFlag("kubelet-https",true),
		printFlag("runtime-config","api/all"),
		printFlag("service-account-key-file","/var/lib/kubernetes/serviceaccount/service-account.pem"),
		printFlag("service-cluster-ip-range", apiServer.serviceClusterIPRange),
		printFlag("service-node-port-range","30000-32767"),
		printFlag("tls-cert-file","/var/lib/kubernetes/kubernetes.pem"),
		printFlag("tls-private-key-file","/var/lib/kubernetes/kubernetes-key.pem"),
		printFlag("v",2),
	}
}