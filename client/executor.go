package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"github.com/tommenx/coordinator/pkg/isolate"
	"github.com/tommenx/coordinator/pkg/resource"
	"github.com/tommenx/coordinator/pkg/util"
	ecpb "github.com/tommenx/pvproto/pkg/proto/executorpb"
)

type server struct{}

var (
	nodeId string
	ip     string
	port   string
)

func init() {
	flag.Set("logtostderr", "true")
	flag.StringVar(&nodeId, "nodeid", "host1", "use to identify host")
	flag.StringVar(&ip, "ip", "127.0.0.1", "identify ip address")
	flag.StringVar(&port, "port", "50051", "identify port")

}

func (s *server) PutIsolation(ctx context.Context, req *ecpb.PutIsolationRequest) (*ecpb.PutIsolationResponse, error) {
	rsp := &ecpb.PutIsolationResponse{
		Header: &ecpb.ResponseHeader{
			Error: &ecpb.Error{},
		},
	}
	glog.V(4).Infof("get req %+v\n", req)
	device := req.Deivice
	limits := req.Resource
	settings := []*isolate.BlkioSetting{}
	for _, limit := range limits {
		if limit.Type == ecpb.StorageType_STORAGE {
			continue
		}
		setting := &isolate.BlkioSetting{
			Type:  limit.Kind,
			Maj:   device.Maj,
			Min:   device.Min,
			Count: util.ToCount(int64(limit.Size), limit.Unit),
		}
		settings = append(settings, setting)
	}
	tempPath := "/Users/tommenx/Desktop/cgroup"
	err := isolate.NewBlkio(tempPath).Update(isolate.DefaultDir, settings)
	if err != nil {
		rsp.Header.Error.Type = ecpb.ErrorType_INTERNAL_ERROR
		rsp.Header.Error.Message = fmt.Sprintf("update blkio error,err=%v", err)
	}
	rsp.Header.Error.Type = ecpb.ErrorType_OK
	rsp.Header.Error.Message = "success"
	return rsp, nil

}

func main() {
	flag.Parse()
	port = ":" + port
	lst, err := net.Listen("tcp", port)
	if err != nil {
		glog.Error("listen error")
	}
	s := grpc.NewServer()
	ecpb.RegisterExecutorServer(s, &server{})
	// register to etcd
	client, err := resource.New([]string{"127.0.0.1:2379"})
	if err != nil {
		glog.Errorf("create resource client error, err=%v", err)
	}
	address := fmt.Sprintf("%s%s", ip, port)
	executor := &resource.Executor{
		Hostname: nodeId,
		Address:  address,
	}
	go client.Executor().Register(executor)
	glog.V(4).Infof("success register, host=%s, addr=%s", nodeId, address)
	s.Serve(lst)
}
