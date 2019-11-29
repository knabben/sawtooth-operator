package assets

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateService returns service spec
func CreateService(serviceNumber string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("service-%s", serviceNumber),
			Namespace: "default",

		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"version": fmt.Sprintf("sawtooth-%s", serviceNumber)},
			ClusterIP: "",
			Ports: []corev1.ServicePort{
				{
					Name: "peer",
					Port: 8800,
				},
			},
		},
	}
}