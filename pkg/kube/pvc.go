package kube

import (
	"errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ErrNotBounded = errors.New("volume not bounded")

type PVCHandler struct {
	c *kubeClient
}

var h *PVCHandler

func init() {
	c := newKubeClient()
	h = &PVCHandler{c: c}
}

// use pvc name to find refered pv name
func GetBoundedPVByPVC(ns, name string) (string, error) {
	var pv string
	pvc, err := h.c.Clientset.CoreV1().PersistentVolumeClaims(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return pv, err
	}
	if pvc.Status.Phase == v1.ClaimBound {
		pv = pvc.Spec.VolumeName
		return pv, nil
	}
	return name, ErrNotBounded
}
