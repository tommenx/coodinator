package resource

import (
	"fmt"
	"log"
	"testing"
)

// func TestNodeCreate(t *testing.T) {
// 	client, err := NewCoreClient()
// 	if err != nil {
// 		t.Errorf("create client error,%v", err)
// 	}
// 	node := &Node{
// 		Name: "test",
// 		Storages: []Storage{
// 			{
// 				Name:  "ssd-1",
// 				Level: SSD,
// 				Resources: []Resource{
// 					{
// 						Type: STORAGE,
// 						Kind: "space",
// 						Size: 100,
// 						Unit: GB,
// 					},
// 					{
// 						Type: LIMIT,
// 						Kind: "read_bps",
// 						Size: 1024,
// 						Unit: MB,
// 					},
// 				},
// 			},
// 			{
// 				Name:  "hdd-1",
// 				Level: HDD,
// 				Resources: []Resource{
// 					{
// 						Type: STORAGE,
// 						Kind: "space",
// 						Size: 1000,
// 						Unit: GB,
// 					},
// 					{
// 						Type: LIMIT,
// 						Kind: "read_bps",
// 						Size: 512,
// 						Unit: MB,
// 					},
// 				},
// 			},
// 		},
// 	}
// 	_, err = client.Node("test").Create(node)
// 	if err != nil {
// 		t.Errorf("error,%v", err)
// 	}
// }

// func TestNodeDelete(t *testing.T) {
// 	client, _ := New()
// 	node := &Node{
// 		Name: "zx",
// 	}
// 	client.Node().Delete("zx")
// 	if err := client.Node().Delete("zx"); err != nil {
// 		if err != ErrKeyNotExist {
// 			t.Errorf("key %s have not deleted", node.Name)
// 		}
// 	} else {
// 		t.Errorf("key %s have not deleted", node.Name)
// 	}
// }

// func TestDeleteNode(t *testing.T) {
// 	endpoints := []string{"127.0.0.1:2379"}
// 	cli, _ := New(endpoints)
// 	err1 := cli.Node().Delete("node1")
// 	err2 := cli.Node().Delete("node2")
// 	err3 := cli.Node().Delete("node3")
// 	log.Printf("%v%v%v", err1, err2, err3)
// }

func TestNode(t *testing.T) {
	endpoints := []string{"127.0.0.1:2379"}
	cli, _ := New(endpoints)
	nodes := []Node{
		{
			Name: "node1",
		},
		{
			Name: "node2",
		},
		{
			Name: "node3",
		},
	}
	for _, node := range nodes {
		cli.Node().Create(&node)
	}
	got, _ := cli.Node().GetAll()
	log.Printf("len is %d", len(got))
	for _, v := range got {
		fmt.Printf("%+v\n", v)
	}

	for _, node := range nodes {
		cli.Node().Delete(node.Name)
	}

}
