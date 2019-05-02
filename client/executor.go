package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"github.com/tommenx/coordinator/pkg/resource"
	pb "github.com/tommenx/pvproto/pkg/proto/executorpb"
)

// func main() {
// 	ServerAddr := "127.0.0.1:50051"
// 	conn, err := grpc.Dial(ServerAddr, grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatal("error")
// 	}
// 	c := ecpb.NewExecutorClient(conn)
// 	res, err := c.PutIsolation(context.Background(), &ecpb.PutIsolationRequest{
// 		Header: &ecpb.RequestHeader{
// 			NodeId: "client",
// 		},
// 		Resource: []*ecpb.Resource{
// 			&ecpb.Resource{
// 				Type: ecpb.StorageType_STORAGE,
// 				Kind: "size",
// 				Size: uint64(1000),
// 				Unit: ecpb.Unit_B,
// 			},
// 		},
// 	})
// 	if err != nil {
// 		log.Printf("error,%v", err)
// 	}
// 	log.Printf("res:%v", res)
// }

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

func (s *server) PutIsolation(ctx context.Context, req *pb.PutIsolationRequest) (*pb.PutIsolationResponse, error) {
	glog.V(4).Infof("get req %+v\n", req)
	return &pb.PutIsolationResponse{
		Header: &pb.ResponseHeader{
			NodeId: nodeId,
			Error: &pb.Error{
				Type:    pb.ErrorType_OK,
				Message: "success",
			},
		},
	}, nil
}

func main() {
	flag.Parse()
	port = ":" + port
	lst, err := net.Listen("tcp", port)
	if err != nil {
		glog.Error("listen error")
	}
	s := grpc.NewServer()
	pb.RegisterExecutorServer(s, &server{})
	// register to etcd
	client, err := resource.New()
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
