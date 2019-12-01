package assets

import (
	"fmt"
	sawtoothv1alpha1 "github.com/knabben/sawtooth-operator/pkg/apis/sawtooth/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewPod returns a busybox pod with the same name/namespace as the cr
func (s *Sawtooth) NewPod(cr *sawtoothv1alpha1.Sawtooth, podName string, number int, peerArgs []string) *corev1.Pod {
	objectMeta := metav1.ObjectMeta{
		Name:      podName,
		Namespace: cr.Namespace,
		Labels:    s.GenerateSelector(),
	}

	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "validator-priv",
			MountPath: "/etc/sawtooth/keys",
		},
	}

	return &corev1.Pod{
		ObjectMeta: objectMeta,
		Spec: corev1.PodSpec{
			InitContainers: createInitContainer(validatorImage),
			Containers: []corev1.Container{
				ValidatorContainer(volumeMounts, peerArgs),
				RestAPIContainer(volumeMounts),
			},
			Volumes: s.GeneratePrivateKeyVolume(),
		},
	}
}

// ValidatorContainer returns validator container spec
func (s *Sawtooth) ValidatorContainer(volumeMount []corev1.VolumeMount, peerArgs []string) corev1.Container {
	command := append([]string{
		"sawtooth-validator", "-vv",
		"--endpoint", s.Endpoint(),
		"--peering", "dynamic",
		"--bind", "component:tcp://eth0:4004",
		"--bind", "consensus:tcp://eth0:5050",
		"--bind", "network:tcp://eth0:8800",
	}, peerArgs...)

	containerName := fmt.Sprintf("validator-%d", s.NodeNumber)

	return corev1.Container{
		Name:         containerName,
		Image:        s.ValidatorImage(),
		Command:      command,
		VolumeMounts: volumeMount,
	}
}

// RestAPIContainer returns the Rest API container spec
func (s *Sawtooth) RestAPIContainer(volumeMount []corev1.VolumeMount) corev1.Container {
	containerName := fmt.Sprintf("rest-api-%d", s.NodeNumber)
	restApiCommand := append([]string{
		"sawtooth-rest-api",
		"-vv",
		"-C",
		fmt.Sprintf("tcp://service-%d:4004", s.NodeNumber),
	})

	return corev1.Container{
		Name:         containerName,
		Image:        s.RestAPIImage(),
		Command:      restApiCommand,
		VolumeMounts: volumeMount,
	}
}
