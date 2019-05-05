package kube

import (
	"errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PVCHandler struct {
	c *kubeClient
}

func NewPVCHandler() {
	c := newKubeClient()
	return &PVCHandler{c: c}
}

// use pvc name to find refered pv name
func (h *PVCHandler) GetPodPVByPVCName(ns, name string) (string, error) {
	pvc := h.c.Clientset.CoreV1().PersistentVolumeClaims(ns).Get(name, metav1.GetOptions{})
	if pvc.Status.Phase == v1.VolumeBound {
		return pvc.Spec.VolumeName, nil
	} else {
		return "", errors.New("volueme not bound")
	}

}
