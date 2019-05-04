package server

import (
	"context"
	"fmt"

	"github.com/golang/glog"

	"github.com/tommenx/coordinator/pkg/resource"
	"github.com/tommenx/coordinator/pkg/util"
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

var pvcpv = map[string]string{
	"pvc1": "pv1",
}
var podpvc = map[string]string{
	"nginx": "pvc1",
}

// pod not exist
// scheduler -> coordinator
func (s *Server) PutPodResource(ctx context.Context, req *cdpb.PutPodResourceRequest) (*cdpb.PutPodResourceResponse, error) {
	rsp := &cdpb.PutPodResourceResponse{
		Header: &cdpb.ResponseHeader{
			Error: &cdpb.Error{},
		},
	}
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
	rsp := &cdpb.GetNodeResourceResponse{
		Header: &cdpb.ResponseHeader{
			Error: &cdpb.Error{},
		},
	}
	nodes, err := s.c.Node().GetAll()
	if err != nil {
		glog.Errorf("get all node info error,err=%v", err)
		rsp.Header.Error.Type = cdpb.ErrorType_INTERNAL_ERROR
		rsp.Header.Error.Message = fmt.Sprintf("get all node info error,err=%v", err)
		return rsp, err
	}
	glog.Infof("GetNodeResource have %d nodes", len(nodes))
	for _, node := range nodes {
		rspNode := util.ToRPCNode(node)
		rsp.Node = append(rsp.Node, rspNode)
	}
	rsp.Header.Error.Type = cdpb.ErrorType_OK
	rsp.Header.Error.Message = "success"
	return rsp, nil
}

// executor -> coordinator
func (s *Server) PutNodeResource(ctx context.Context, req *cdpb.PutNodeResourceRequest) (*cdpb.PutNodeResourceResponse, error) {
	rsp := &cdpb.PutNodeResourceResponse{}
	rsp.Header = &cdpb.ResponseHeader{
		Error: &cdpb.Error{},
	}
	node := util.ToResouceNode(req.Node)
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

// executor -> coordinator
// TODO
// use podName to find pvc then get pv
// have not write to ectd
func (s *Server) PutPV(ctx context.Context, req *cdpb.PutPVRequest) (*cdpb.PutPVResponse, error) {
	rsp := &cdpb.PutPVResponse{}
	rsp.Header = &cdpb.ResponseHeader{
		Error: &cdpb.Error{},
	}
	pv := &resource.PV{
		Name: req.Name,
		Device: &resource.Device{
			Id:          req.Device.Id,
			Name:        req.Device.Name,
			Maj:         req.Device.Maj,
			Min:         req.Device.Min,
			VolumeGroup: req.Device.Vg,
			DevicePath:  req.Device.Path,
		},
	}

	podName, pvcName := findPodForPV(pv.Name)

	pod, _ := s.c.Pod("default").Get(podName)
	for _, v := range pod.Allocations {
		if v.PVC.Name == pvcName {
			v.PV = pv
		}
	}
	_, err := s.c.Pod("default").Update(pod)
	if err != nil {
		glog.Errorf("update pod pv error, err=%v", err)
	}
	pod, _ = s.c.Pod("default").Get(podName)
	for _, item := range pod.Allocations {
		if item.PV != nil {
			glog.V(4).Infoln(item.PV.Name)
			glog.V(4).Infoln(*item.PV.Device)
			glog.V(4).Infoln(item.PVC.Name)
		}

	}

	rsp.Header.Error.Type = cdpb.ErrorType_OK
	rsp.Header.Error.Message = "success"
	return rsp, nil
}

// input  PV name
// output pod name
func findPodForPV(name string) (string, string) {
	var pvc string
	var pod string
	for k, v := range pvcpv {
		if v == name {
			pvc = k
			break
		}
	}
	for k, v := range podpvc {
		if v == pvc {
			pod = k
			break
		}
	}
	return pod, pvc
}

func (s *Server) PutPodNodeInfo(ctx context.Context, req *cdpb.PutPodNodeInfoRequest) (*cdpb.PutPodNodeInfoResponse, error) {
	rsp := &cdpb.PutPodNodeInfoResponse{
		Header: &cdpb.ResponseHeader{
			Error: &cdpb.Error{},
		},
	}
	pod, err := s.c.Pod(req.Namespace).Get(req.Pod)
	if err != nil {
		glog.Errorf("get pod node info error, err=%v", err)
		rsp.Header.Error.Type = cdpb.ErrorType_INTERNAL_ERROR
		rsp.Header.Error.Message = fmt.Sprintf("get pod info error, err=%v", err)
		return rsp, err
	}
	pod.Node = req.Node
	_, err = s.c.Pod(req.Namespace).Update(pod)
	if err != nil {
		glog.Errorf("update pod node info error, err=%v", err)
		rsp.Header.Error.Type = cdpb.ErrorType_INTERNAL_ERROR
		rsp.Header.Error.Message = fmt.Sprintf("get pod info error, err=%v", err)
		return rsp, err
	}
	rsp.Header.Error.Type = cdpb.ErrorType_OK
	rsp.Header.Error.Message = "success"
	return rsp, nil
}
