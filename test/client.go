package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/golang/glog"
	"github.com/tommenx/coordinator/pkg/resource"
	"github.com/tommenx/coordinator/pkg/util"

	cdpb "github.com/tommenx/pvproto/pkg/proto/coordinatorpb"
	"google.golang.org/grpc"
)

func init() {
	flag.Set("logtostderr", "true")
}

var (
	address = "127.0.0.1:50051"
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic("create conn error")
	}
	client := cdpb.NewCoordinatorClient(conn)
	req := &cdpb.PutPVRequest{
		Name: "pv1",
		Device: &cdpb.Device{
			Name: "lovl1",
			Maj:  "1",
			Min:  "255",
			Id:   "2131231231231232",
			Path: "/dev/vgdata/lvol1/",
			Vg:   "vgdata",
		},
	}
	rsp, err := client.PutPV(context.TODO(), req)
	fmt.Println(rsp, err)
}

func pod() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic("create conn error")
	}
	client := cdpb.NewCoordinatorClient(conn)
	pod := &resource.Pod{
		Name:      "nginx",
		Namespace: "default",
		Allocations: []*resource.Allocation{
			&resource.Allocation{
				PVC: &resource.PVC{
					Name: "pvc1",
					Resources: []*resource.Resource{
						&resource.Resource{
							Type: resource.STORAGE,
							Kind: "space",
							Size: 10,
							Unit: resource.GB,
						},
						&resource.Resource{
							Type: resource.LIMIT,
							Kind: "read_bps",
							Size: 100,
							Unit: resource.MB,
						},
					},
				},
				PV: &resource.PV{},
			},
			&resource.Allocation{
				PVC: &resource.PVC{
					Name: "pvc2",
					Resources: []*resource.Resource{
						&resource.Resource{
							Type: resource.STORAGE,
							Kind: "space",
							Size: 20,
							Unit: resource.GB,
						},
						&resource.Resource{
							Type: resource.LIMIT,
							Kind: "read_bps",
							Size: 200,
							Unit: resource.MB,
						},
					},
				},
				PV: &resource.PV{},
			},
		},
	}
	rsp1, err := client.PutPodResource(context.TODO(), util.ToRPCPod(pod))
	fmt.Println(rsp1, err)
	rsp2, err := client.PutPodNodeInfo(context.TODO(), &cdpb.PutPodNodeInfoRequest{
		Pod:       pod.Name,
		Namespace: pod.Namespace,
		Node:      "host1",
	})
	fmt.Println(rsp2, err)
}

func node() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic("create conn error")
	}
	client := cdpb.NewCoordinatorClient(conn)
	host1 := &resource.Node{
		Name: "host1",
		Storages: []*resource.Storage{
			&resource.Storage{
				Name:  "ssd1",
				Level: resource.SSD,
				Resources: []*resource.Resource{
					&resource.Resource{
						Type: resource.STORAGE,
						Kind: "space",
						Size: 1000,
						Unit: resource.GB,
					},
					&resource.Resource{
						Type: resource.LIMIT,
						Kind: "read_bps",
						Size: 500,
						Unit: resource.MB,
					},
					&resource.Resource{
						Type: resource.LIMIT,
						Kind: "write_bps",
						Size: 100,
						Unit: resource.MB,
					},
				},
			},
			&resource.Storage{
				Name:  "hdd1",
				Level: resource.HDD,
				Resources: []*resource.Resource{
					&resource.Resource{
						Type: resource.STORAGE,
						Kind: "space",
						Size: 10000,
						Unit: resource.GB,
					},
					&resource.Resource{
						Type: resource.LIMIT,
						Kind: "read_bps",
						Size: 200,
						Unit: resource.MB,
					},
					&resource.Resource{
						Type: resource.LIMIT,
						Kind: "write_bps",
						Size: 50,
						Unit: resource.MB,
					},
				},
			},
		},
	}
	host2 := &resource.Node{
		Name: "host2",
		Storages: []*resource.Storage{
			&resource.Storage{
				Name:  "ssd1",
				Level: resource.SSD,
				Resources: []*resource.Resource{
					&resource.Resource{
						Type: resource.STORAGE,
						Kind: "space",
						Size: 700,
						Unit: resource.GB,
					},
					&resource.Resource{
						Type: resource.LIMIT,
						Kind: "read_bps",
						Size: 300,
						Unit: resource.MB,
					},
					&resource.Resource{
						Type: resource.LIMIT,
						Kind: "write_bps",
						Size: 100,
						Unit: resource.MB,
					},
				},
			},
			&resource.Storage{
				Name:  "hdd1",
				Level: resource.HDD,
				Resources: []*resource.Resource{
					&resource.Resource{
						Type: resource.STORAGE,
						Kind: "space",
						Size: 20000,
						Unit: resource.GB,
					},
					&resource.Resource{
						Type: resource.LIMIT,
						Kind: "read_bps",
						Size: 200,
						Unit: resource.MB,
					},
					&resource.Resource{
						Type: resource.LIMIT,
						Kind: "write_bps",
						Size: 50,
						Unit: resource.MB,
					},
				},
			},
		},
	}
	req1 := &cdpb.PutNodeResourceRequest{
		Node: util.ToRPCNode(host1),
	}
	req2 := &cdpb.PutNodeResourceRequest{
		Node: util.ToRPCNode(host2),
	}
	rsp, err := client.PutNodeResource(context.TODO(), req1)
	log.Println(rsp, err)
	rsp, err = client.PutNodeResource(context.TODO(), req2)
	log.Println(rsp, err)
	req3 := &cdpb.GetNodeResourceRequest{}
	r, _ := client.GetNodeResource(context.TODO(), req3)
	glog.V(4).Infof("%+v\n", r)
	glog.V(4).Infof("len = %d", len(r.Node))

}
