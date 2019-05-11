package kube

import (
	"log"
	"testing"
)

func TestGetBoundedPVByPVC(*testing.T) {
	ns := "default"
	pvc := "lvm-pvc-1"
	pv, err := GetBoundedPVByPVC(ns, pvc)
	log.Println(pv, err)
}
