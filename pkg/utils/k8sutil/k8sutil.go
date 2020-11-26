package k8sutil

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsControlledBy(reference []metav1.OwnerReference, kind string, name string) bool {
	for _, ref := range reference {
		if ref.Kind == kind && (name == "" || ref.Name == name) {
			return true
		}
	}
	return false
}
