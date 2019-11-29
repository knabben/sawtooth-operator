package utils

import (
	"fmt"
	sawtoothv1alpha1 "github.com/knabben/sawtooth-operator/pkg/apis/sawtooth/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateLabel(name string) map[string]string {
	return map[string]string{
		"app": name,
	}
}

func newHostPathType(hostType string) *corev1.HostPathType {
	hostPathType := new(corev1.HostPathType)
	*hostPathType = corev1.HostPathType(hostType)
	return hostPathType
}

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
func CreatePod(cr *sawtoothv1alpha1.Sawtooth) *corev1.Pod {
	labels := generateLabel(cr.Name)
	imageName := fmt.Sprintf("hyperledger/sawtooth-validator:%s", cr.Spec.Version)

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			InitContainers: createInitContainer(imageName),
			Containers: []corev1.Container{
				{
					Name:    "sawtooth-1",
					Image:   imageName,
					Command: []string{
						"sawtooth-validator", "-vv",
						"--endpoint", "tcp://$SAWTOOTH_0_SERVICE_HOST:8800",
						"--bind", "component:tcp://eth0:4004",
						"--bind", "consensus:tcp://eth0:5050",
						"--bind", "network:tcp://eth0:8800",
					},
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
