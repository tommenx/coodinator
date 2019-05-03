package server

import (
	"context"
	"fmt"

	"github.com/golang/glog"

	"github.com/tommenx/coordinator/pkg/resource"
	cdpb "github.com/tommenx/pvproto/pkg/proto/coordinatorpb"
)

type Server struct {
	c *resource.CoreClient
}

func New(endpoonts []string) (*Server, error) {
	c, err := resource.New(endpoonts)
	if err != nil {
		return nil, err
	}
	return &Server{c: c}, nil
}

// pod not exist
// scheduler -> coordinator
func (s *Server) PutPodResource(ctx context.Context, req *cdpb.PutPodResourceRequest) (*cdpb.PutPodResourceResponse, error) {
	rsp := &cdpb.PutPodResourceResponse{}
	pod := &resource.Pod{}
	pod.Name = req.Name
	pod.Namespace = req.Namespace
	for _, v := range req.Pvc {
		pvc := &resource.PVC{}
		pvc.Name = v.Name
		for _, t := range v.Resource {
			r := &resource.Resource{
				Type: resource.Type[t.Type],
				Kind: t.Kind,
				Size: t.Size,
				Unit: resource.Unit[t.Unit],
			}
			pvc.Resources = append(pvc.Resources, r)
		}
		allo := &resource.Allocation{}
		allo.PVC = pvc
		pod.Allocations = append(pod.Allocations, allo)
	}
	// add to etcd
	_, err := s.c.Pod(pod.Namespace).Create(pod)
	if err != nil {
		glog.Errorf("create pod error, err=%+v", err)
		rsp.Header.Error.Type = cdpb.ErrorType_INTERNAL_ERROR
		rsp.Header.Error.Message = fmt.Sprintf("create pod error, err=%+v", err)
	}

	rsp.Header.Error.Type = cdpb.ErrorType_OK
	rsp.Header.Error.Message = "success"
	return rsp, nil
}

// scheduler get all node resource
func (s *Server) GetNodeResource(ctx context.Context, req *cdpb.GetNodeResourceRequest) (*cdpb.GetNodeResourceResponse, error) {
	rsp := &cdpb.GetNodeResourceResponse{}
	nodes, err := s.c.Node().GetAll()
	if err != nil {
		glog.Errorf("get all node info error,err=%v", err)
		rsp.Header.Error.Type = cdpb.ErrorType_INTERNAL_ERROR
		rsp.Header.Error.Message = fmt.Sprintf("get all node info error,err=%v", err)
		return rsp, err
	}
	for _, node := range nodes {
		rspNode := &cdpb.Node{}
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
		rsp.Node = append(rsp.Node, rspNode)
	}
	rsp.Header.Error.Type = cdpb.ErrorType_OK
	rsp.Header.Error.Message = "success"
	return rsp, nil
}

// executor -> coordinator
func (s *Server) PutNodeResource(ctx context.Context, req *cdpb.PutNodeResourceRequest) (*cdpb.PutNodeResourceResponse, error) {
	rsp := &cdpb.PutNodeResourceResponse{}
	node := &resource.Node{Name: req.Node.Name}
	for _, v := range req.Node.Storage {
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
		_, err := s.c.Node().Create(node)
		if err != nil {
			glog.Errorf("etcd create node %s error, err=%v", node.Name, err)
			rsp.Header.Error.Type = cdpb.ErrorType_INTERNAL_ERROR
			rsp.Header.Error.Message = fmt.Sprintf("etcd create node %s error, err=%v", node.Name, err)
			return rsp, err
		}
		rsp.Header.Error.Type = cdpb.ErrorType_OK
		rsp.Header.Error.Message = "success"
		return rsp, nil

	}

	return nil, nil
}

func (s *Server) PutPV(ctx context.Context, req *cdpb.PutPVRequest) (*cdpb.PutPVResponse, error) {
	return nil, nil
}
