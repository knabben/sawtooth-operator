package assets

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewService returns Sawtooth TCP ports service
func (s *Sawtooth) NewService() *corev1.Service {
	typeMeta := metav1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	}
	objectMeta := metav1.ObjectMeta{
		Name:      s.GenerateServiceName(),
		Namespace: "default",
	}

	return &corev1.Service{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
		Spec: corev1.ServiceSpec{
			Selector:  s.GenerateSelector(),
			ClusterIP: "",
			Ports:     s.GenerateTCPPorts(),
		},
	}
}

// GenerateTCPPorts returns the enabled ports for Sawtooth
func (s *Sawtooth) GenerateTCPPorts() []corev1.ServicePort {
	return []corev1.ServicePort{
		{
			Name: "peer",
			Port: 8800,
		},
		{
			Name: "component",
			Port: 4004,
		},
		{
			Name: "http",
			Port: 8008,
		},
	}
}

// GenerateSelector generates the correct POD selector
func (s *Sawtooth) GenerateSelector() map[string]string {
	return map[string]string{
		"sawtooth-node": fmt.Sprintf("sawtooth-%d", s.NodeNumber),
	}
}

// GetServiceName returns the generated name for Sawtooth node
func (s *Sawtooth) GenerateServiceName() string {
	return fmt.Sprintf("service-%d", s.NodeNumber)
}
