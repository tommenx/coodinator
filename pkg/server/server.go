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
		pvc := resource.PVC{}
		pvc.Name = v.Name
		for _, t := range v.Resource {
			r := resource.Resource{
				Type: resource.Type[t.Type],
				Kind: t.Kind,
				Size: t.Size,
				Unit: resource.Unit[t.Unit],
			}
			pvc.Resources = append(pvc.Resources, r)
		}
		allo := resource.Allocation{}
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
	return nil, nil
}

// executor -> coordinator
func (s *Server) PutNodeResource(ctx context.Context, req *cdpb.PutNodeResourceRequest) (*cdpb.PutNodeResourceRequest, error) {
	return nil, nil
}

func (s *Server) PutPV(ctx context.Context, req *cdpb.PutPVRequest) (*cdpb.PutPVResponse, error) {
	return nil, nil
}
