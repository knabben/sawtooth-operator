package assets

import (
	//sawtoothv1alpha1 "github.com/knabben/sawtooth-operator/pkg/apis/sawtooth/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetLabel() map[string]string {
	return map[string]string{
		"app": "sawtooth",
	}
}

func newHostPathType(hostType string) *corev1.HostPathType {
	hostPathType := new(corev1.HostPathType)
	*hostPathType = corev1.HostPathType(hostType)
	return hostPathType
}