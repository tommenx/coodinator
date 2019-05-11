package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/tommenx/coordinator/pkg/db"

	"github.com/golang/glog"
	"github.com/tommenx/coordinator/pkg/resource"
	pb "github.com/tommenx/pvproto/pkg/proto/executorpb"
)

func init() {
	flag.Set("logtostderr", "true")
}

var conns map[string]pb.ExecutorClient

func main() {
	flag.Parse()
	client, err := resource.New([]string{"127.0.0.1:2379"})
	if err != nil {
		glog.Errorf("new client error")
	}
	conns = make(map[string]pb.ExecutorClient)
	go client.Executor().Watch(db.FOLDER_EXECUTOR_INFO,
		func(k, v string) {
			log.Printf("add new executor,key=%s, val=%s\n", k, v)
			conn, _ := grpc.Dial(v, grpc.WithInsecure())
			c := pb.NewExecutorClient(conn)
			conns[k] = c
		}, func(k, v string) {
			log.Printf("remove executor,key=%s\n", k)
			delete(conns, k)
		},
	)

	time.Sleep(15 * time.Second)
	for _, v := range conns {
		res, _ := v.PutIsolation(context.TODO(), &pb.PutIsolationRequest{
			Header: &pb.RequestHeader{
				NodeId: "client",
			},
			Deivice: &pb.Device{
				Maj: 1,
				Min: 2,
				Id:  "id1",
			},
			Resource: []*pb.Resource{
				&pb.Resource{
					Type: pb.StorageType_STORAGE,
					Kind: "size",
					Size: uint64(1000),
					Unit: pb.Unit_B,
				},
				&pb.Resource{
					Type: pb.StorageType_LIMIT,
					Kind: "throttle.write_bps_device",
					Size: 10,
					Unit: pb.Unit_MB,
				},
				&pb.Resource{
					Type: pb.StorageType_LIMIT,
					Kind: "throttle.read_iops_device",
					Size: 20,
					Unit: pb.Unit_MB,
				},
			},
		})
		log.Printf("got res,%+v\n", res)
	}
	time.Sleep(3 * time.Second)
	log.Println("DONE")

}
