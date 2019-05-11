package util

import (
	"github.com/tommenx/coordinator/pkg/resource"
	cdpb "github.com/tommenx/pvproto/pkg/proto/coordinatorpb"
	ecpb "github.com/tommenx/pvproto/pkg/proto/executorpb"
)

func ToRPCNode(node *resource.Node) *cdpb.Node {
	rspNode := &cdpb.Node{Name: node.Name}
	for _, s := range node.Storages {
		rspStorage := &cdpb.Storage{
			Name:  s.Name,
			Level: resource.ReLeval[s.Level],
		}
		for _, r := range s.Resources {
			rspResouce := &cdpb.Resource{
				Type: resource.ReType[r.Type],
				Kind: r.Kind,
				Size: r.Size,
				Unit: resource.ReUnit[r.Unit],
			}
			rspStorage.Resource = append(rspStorage.Resource, rspResouce)
		}
		rspNode.Storage = append(rspNode.Storage, rspStorage)
	}
	return rspNode
}

func ToResouceNode(req *cdpb.Node) *resource.Node {
	node := &resource.Node{Name: req.Name}
	for _, v := range req.Storage {
		store := &resource.Storage{Name: v.Name}
		for _, r := range v.Resource {
			resource := &resource.Resource{
				Type: resource.Type[r.Type],
				Kind: r.Kind,
				Size: r.Size,
				Unit: resource.Unit[r.Unit],
			}
			store.Resources = append(store.Resources, resource)
		}
		node.Storages = append(node.Storages, store)
	}
	return node
}

func ToRPCPod(pod *resource.Pod) *cdpb.PutPodResourceRequest {
	req := &cdpb.PutPodResourceRequest{
		Name:      pod.Name,
		Namespace: pod.Namespace,
	}
	pvcs := []*cdpb.PVC{}
	for _, a := range pod.Allocations {
		pvc := &cdpb.PVC{Name: a.PVC.Name}
		resources := []*cdpb.Resource{}
		for _, r := range a.PVC.Resources {
			resource := &cdpb.Resource{
				Type: resource.ReType[r.Type],
				Kind: r.Kind,
				Size: r.Size,
				Unit: resource.ReUnit[r.Unit],
			}
			resources = append(resources, resource)
		}
		pvc.Resource = resources
		pvcs = append(pvcs, pvc)
	}
	req.Pvc = pvcs
	return req
}

func ToCount(num int64, unit ecpb.Unit) int64 {
	var KB, MB, GB int64
	KB = 1 << 10
	MB = KB << 10
	GB = MB << 10
	if unit == ecpb.Unit_B {
		return num
	} else if unit == ecpb.Unit_KB {
		return num * KB
	} else if unit == ecpb.Unit_MB {
		return num * MB
	} else if unit == ecpb.Unit_GB {
		return num * GB
	}
	return 0
}
