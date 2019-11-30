package assets

import (
	"fmt"
	sawtoothv1alpha1 "github.com/knabben/sawtooth-operator/pkg/apis/sawtooth/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createVolumes generates batch.config file and validator.priv key
func createVolumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "validator-priv",
			VolumeSource: corev1.VolumeSource {
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}
}

// createInitContainer generates system key and genesis block
func createInitContainer(imageName string) []corev1.Container {
	initContainers := []corev1.Container{
		{
			Name: "genesis-init",
			Image: imageName,
			Command: []string{"sawadm", "keygen", "--force"},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name: "validator-priv",
					MountPath: "/etc/sawtooth/keys",
				},
			},
		},
	}
	return initContainers
}


// newPodForCR returns a busybox pod with the same name/namespace as the cr
func CreatePodSpec(cr *sawtoothv1alpha1.Sawtooth, podName string, number int, peerArgs []string) *corev1.Pod {
	labels := GetLabel()
	labels["version"] = fmt.Sprintf("sawtooth-%d", number)
	
	command := append([]string{
		"sawtooth-validator", "-vv",
		"--endpoint", fmt.Sprintf("tcp://service-%d:8800", number),
		"--peering", "dynamic",
		"--bind", "component:tcp://eth0:4004",
		"--bind", "consensus:tcp://eth0:5050",
		"--bind", "network:tcp://eth0:8800",
	}, peerArgs...)

	validatorImage := fmt.Sprintf("hyperledger/sawtooth-validator:%s", cr.Spec.Version)

	restApiCommand := append([]string{
		"sawtooth-rest-api", "-C", "tcp://eth0:4004",
	})

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			InitContainers: createInitContainer(validatorImage),
			Containers: []corev1.Container{
				{
					Name:    fmt.Sprintf("validator-%s", podName),
					Image:   validatorImage,
					Command: command,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name: "validator-priv",
							MountPath: "/etc/sawtooth/keys",
						},
					},
				},
				{
					Name:    fmt.Sprintf("rest-api-%s", podName),
					Image:   fmt.Sprintf("hyperledger/sawtooth-rest-api:%s", cr.Spec.Version),
					Command: restApiCommand,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name: "validator-priv",
							MountPath: "/etc/sawtooth/keys",
						},
					},
				},
			},
			Volumes: createVolumes(),
		},
	}
}
