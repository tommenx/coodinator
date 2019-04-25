package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tommenx/coordinator/pkg/db"
	"go.etcd.io/etcd/clientv3"
)

func main() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")
	defer cli.Close()
	//设置1秒超时，访问etcd有超时控制
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// //操作etcd
	// _, err = cli.Put(ctx, "/logagent/conf/", "sample_value")
	// //操作完毕，取消etcd
	// cancel()
	// if err != nil {
	// 	fmt.Println("put failed, err:", err)
	// 	return
	// }

	h, err := db.NewEtcdHandler([]string{"localhost:2379"})
	if err != nil {
		fmt.Println("err", err)
	}
	h.Put("/logagent", "conf", "aaaa")

	//取值，设置超时为1秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "/logagent/conf/")
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
	// h, err := db.NewEtcdHandler([]string{"localhost:2379"})
	// if err != nil {
	// 	fmt.Println("err", err)
	// }
	// res, err := h.Get("/logagent", "conf")
	// if err != nil {
	// 	fmt.Println("err", err)
	// }
	// for k, v := range res {
	// 	fmt.Println(k)
	// 	fmt.Println(string(v))
	// }

}
