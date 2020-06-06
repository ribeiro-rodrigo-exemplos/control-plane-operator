package apiserver

import corev1 "k8s.io/api/core/v1"

func newAPIServerContainer()corev1.Container{
	return corev1.Container{
		Name: "kube-apiserver",
		Image: "rodrigoribeiro/globo-kube-apiserver",
		VolumeMounts: newVolumeMounts(),
		Command: newCommands(),
	}
}

func newVolumeMounts()[]corev1.VolumeMount{
	return []corev1.VolumeMount{
		{Name: "ca", MountPath: "/var/lib/kubernetes/ca", ReadOnly: true},
		{Name: "kubernetes", MountPath: "/var/lib/kubernetes", ReadOnly: true},
		{Name: "encryption", MountPath: "/var/lib/kubernetes/encryption", ReadOnly: true},
	}
}

func newCommands()[]string{
	return []string{
		"kube-apiserver",
		"--advertise-address=192.168.39.42",
		"--allow-privileged=true",
		"--apiserver-count=1",
		"--audit-log-maxage=30",
		"--audit-log-maxbackup=3",
		"--audit-log-maxsize=100",
		"--audit-log-path=/var/log/audit.log",
		"--authorization-mode=Node,RBAC",
		"--bind-address=0.0.0.0",
		"--client-ca-file=/var/lib/kubernetes/ca/ca.pem",
		"--enable-admission-plugins=NamespaceLifecycle,NodeRestriction,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota",
		"--etcd-cafile=/var/lib/kubernetes/ca/ca.pem",
		"--etcd-certfile=/var/lib/kubernetes/kubernetes.pem",
		"--etcd-keyfile=/var/lib/kubernetes/kubernetes-key.pem",
		"--etcd-servers=https://161.35.116.213:2379",
		"--event-ttl=1h",
		"--encryption-provider-config=/var/lib/kubernetes/encryptation/encryption-config.yaml",
		"--kubelet-certificate-authority=/var/lib/kubernetes/ca/ca.pem",
		"--kubelet-client-certificate=/var/lib/kubernetes/kubernetes.pem",
		"--kubelet-client-key=/var/lib/kubernetes/kubernetes-key.pem",
		"--kubelet-https=true",
		"--runtime-config=api/all",
		"--service-account-key-file=/var/lib/kubernetes/serviceaccount/service-account.pem",
		"--service-cluster-ip-range=10.32.0.0/24",
		"--service-node-port-range=30000-32767",
		"--tls-cert-file=/var/lib/kubernetes/kubernetes.pem",
		"--tls-private-key-file=/var/lib/kubernetes/kubernetes-key.pem",
		"--v=2",
	}
}