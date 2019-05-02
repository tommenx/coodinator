package resource

import "testing"

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

func TestNodeDelete(t *testing.T) {
	client, _ := New()
	node := &Node{
		Name: "zx",
	}
	client.Node().Delete("zx")
	if err := client.Node().Delete("zx"); err != nil {
		if err != ErrKeyNotExist {
			t.Errorf("key %s have not deleted", node.Name)
		}
	} else {
		t.Errorf("key %s have not deleted", node.Name)
	}
}
