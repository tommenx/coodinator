package main

import (
	"context"
	"log"
	"net"

	ecpb "github.com/tommenx/pvproto/pkg/proto/executorpb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{}

func (*server) PutIsolation(ctx context.Context, req *ecpb.PutIsolationRequest) (*ecpb.PutIsolationResponse, error) {
	log.Printf("%v\n", req)
	return &ecpb.PutIsolationResponse{
		Header: &ecpb.ResponseHeader{
			NodeId: "server",
		},
	}, nil
}

func main() {
	lst, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("error listen")
	}
	s := grpc.NewServer()
	ecpb.RegisterExecutorServer(s, &server{})
	s.Serve(lst)
}
