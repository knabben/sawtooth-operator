package assets

import (
	"fmt"
	sawtoothv1alpha1 "github.com/knabben/sawtooth-operator/pkg/apis/sawtooth/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Sawtooth struct {
	Schema     *sawtoothv1alpha1.Sawtooth
	Pod        *corev1.Pod
	Service    *corev1.Service
	NodeNumber int
}

// ValidatorImage generates the image with version tagging
func (s *Sawtooth) ValidatorImage() string {
	return fmt.Sprintf("hyperledger/sawtooth-validator:%s", s.Schema.Spec.Version)
}

// RestAPIImage generates the image with version tagging
func (s *Sawtooth) RestAPIImage() string {
	return fmt.Sprintf("hyperledger/sawtooth-rest-api:%s", s.Schema.Spec.Version)
}

func (s *Sawtooth) Endpoint() string {
	return fmt.Sprintf("tcp://service-%d:8800", s.NodeNumber)
}

// GenerateInitContainers generates system key and genesis block initial container
func (s *Sawtooth) GenerateInitContainers() []corev1.Container {
	mountPath := "/etc/sawtooth/keys"
	command := []string{"sawadm", "keygen", "--force"}

	initContainers := []corev1.Container{
		{
			Name:    "genesis-init",
			Image:   s.ValidatorImage(),
			Command: command,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "validator-priv",
					MountPath: mountPath,
				},
			},
		},
	}
	return initContainers
}

// GeneratePrivateKeyVolume returns the validator shared volume
func (s *Sawtooth) GeneratePrivateKeyVolume() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "validator-priv",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}
}
