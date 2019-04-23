package main

import (
	"context"
	"log"

	ecpb "github.com/tommenx/pvproto/pkg/proto/executorpb"
	"google.golang.org/grpc"
)

func main() {
	ServerAddr := "127.0.0.1:50051"
	conn, err := grpc.Dial(ServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("error")
	}
	c := ecpb.NewExecutorClient(conn)
	res, err := c.PutIsolation(context.Background(), &ecpb.PutIsolationRequest{
		Header: &ecpb.RequestHeader{
			NodeId: "client",
		},
		Resource: []*ecpb.Resource{
			&ecpb.Resource{
				Type: ecpb.StorageType_STORAGE,
				Kind: "size",
				Size: uint64(1000),
				Unit: ecpb.Unit_B,
			},
		},
	})
	if err != nil {
		log.Printf("error,%v", err)
	}
	log.Printf("res:%v", res)
}
