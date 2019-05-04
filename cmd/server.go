package main

import (
	"flag"
	"net"

	"github.com/golang/glog"
	"github.com/tommenx/coordinator/pkg/server"
	cdpb "github.com/tommenx/pvproto/pkg/proto/coordinatorpb"
	"google.golang.org/grpc"
)

func init() {
	flag.Set("logtostderr", "true")
}

var (
	port = ":50051"
)

func main() {
	flag.Parse()
	lst, err := net.Listen("tcp", port)
	if err != nil {
		glog.Errorf("listen port error, err=%v", err)
		panic("listen error")
	}
	s := grpc.NewServer()
	endpoints := []string{"127.0.0.1:2379"}
	srv, err := server.New(endpoints)
	// srv := &fake{}
	if err != nil {
		panic(err)
	}
	cdpb.RegisterCoordinatorServer(s, srv)
	s.Serve(lst)
}
