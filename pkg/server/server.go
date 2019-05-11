package server

import (
	"context"
	"fmt"

	"github.com/golang/glog"

	"github.com/tommenx/coordinator/pkg/kube"
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
// don't know how to find unbounded pod
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

	pods, err := s.c.Pod("").GetAll()
	if err != nil {
		glog.Errorf("get all pods error,err=%v", err)
		rsp.Header.Error.Type = cdpb.ErrorType_INTERNAL_ERROR
		rsp.Header.Error.Message = fmt.Sprintf("get all pods error,err=%v", err)
		return rsp, nil
	}
	// _, err = s.c.Pod("default").Update(pod)
	// if err != nil {
	// 	glog.Errorf("update pod pv error, err=%v", err)
	// }
	pod := s.updateAndGetPod(pv.Name, pods)
	if pod == nil {
		rsp.Header.Error.Type = cdpb.ErrorType_INTERNAL_ERROR
		rsp.Header.Error.Message = "do not find pod"
		return rsp, nil
	}

	for _, a := range pod.Allocations {
		if a.PV.Name == pv.Name {
			a.PV = pv
			break
		}
	}
	_, err = s.c.Pod(pod.Namespace).Update(pod)
	if err != nil {
		rsp.Header.Error.Type = cdpb.ErrorType_INTERNAL_ERROR
		rsp.Header.Error.Message = fmt.Sprintf("pod update error,err=%v,name=%s", err, pod.Name)
		return rsp, nil
	}
	// update succsss
	glog.V(4).Infof("update pv success, pod=%s, pv=%s", pod.Name, pv.Name)
	rsp.Header.Error.Type = cdpb.ErrorType_OK
	rsp.Header.Error.Message = "success"
	return rsp, nil
}

func (s *Server) updateAndGetPod(name string, pods []*resource.Pod) *resource.Pod {
	var findPod *resource.Pod
	for _, p := range pods {
		for i, item := range p.Allocations {
			// 如果pv为空，说明没有绑定pv
			if item.PV == nil {
				if item.PVC != nil {
					//  未设置重试
					if pvname, err := kube.GetBoundedPVByPVC(p.Namespace, item.PVC.Name); err != nil {
						p.Allocations[i].PV.Name = pvname
						_, err := s.c.Pod(p.Namespace).Update(p)
						if err != nil {
							glog.Errorf("put pv update pv name error,err=%v", err)
							continue
						}
						if pvname == name {
							findPod = p
							return findPod
						}
					}
				}
			} else {
				// 只需要校验名字
				if item.PV.Name == name {
					findPod = p
					return findPod
				}

			}

		}
	}
	return nil
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
