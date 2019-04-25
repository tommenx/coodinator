package db

import (
	"fmt"
	"log"
	"testing"
	"time"
)

type KV struct {
	key   string
	value string
}

func TestPutAndGet(t *testing.T) {
	kvs := map[string]string{
		"host1": "cpu100",
		"host2": "cpu80",
		"host3": "cpu110",
	}

	endpoints := []string{"http://127.0.0.1:2379"}
	h, err := NewEtcdHandler(endpoints)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	for k, v := range kvs {
		h.Put(FOLDER_NODE_INFO, k, v)
	}
	for k, v := range kvs {
		res, err := h.Get(FOLDER_NODE_INFO, k)
		if err != nil {
			t.Errorf("get key %s error:%v", k, v)
		}
		tmp := fmt.Sprintf("%s/%s", FOLDER_NODE_INFO, k)
		if val, ok := res[tmp]; !ok {
			t.Errorf("don't have key %s ", tmp)
		} else {
			if val != v {
				t.Errorf("value don't match,put %s,get %s", v, val)
			} else {
				fmt.Printf("key = %s,val = %s\n", tmp, val)
			}
		}

	}

}

func TestDelete(t *testing.T) {
	h, _ := NewEtcdHandler([]string{"http://127.0.0.1:2379"})
	h.Put(FOLDER_NODE_INFO, "test", "aaa")
	res, _ := h.Get(FOLDER_NODE_INFO, "test")
	if len(res) != 1 {
		t.Error("get or put error")
	}
	h.Delete(FOLDER_NODE_INFO, "test")
	res, _ = h.Get(FOLDER_NODE_INFO, "test")
	if len(res) != 0 {
		t.Error("delete error")
	}
}

func TestWatch(y *testing.T) {
	h1, _ := NewEtcdHandler([]string{"http://127.0.0.1:2379"})
	h2, _ := NewEtcdHandler([]string{"http://127.0.0.1:2379"})
	go h1.Watch(FOLDER_NODE_INFO, "watch", func(k, v string) {
		fmt.Printf("add key %s, value %s\n", k, v)
	}, func(k, v string) {
		fmt.Printf("delete key %s, value %s\n", k, v)
	})

	time.Sleep(5 * time.Second)
	h2.Put(FOLDER_NODE_INFO, "watch", "aaaa")
	time.Sleep(3 * time.Second)
	h2.Delete(FOLDER_NODE_INFO, "watch")
	time.Sleep(3 * time.Second)
}

func TestRegister(t *testing.T) {
	h1, _ := NewEtcdHandler([]string{"http://127.0.0.1:2379"})
	h2, _ := NewEtcdHandler([]string{"http://127.0.0.1:2379"})
	h3, _ := NewEtcdHandler([]string{"http://127.0.0.1:2379"})
	h4, _ := NewEtcdHandler([]string{"http://127.0.0.1:2379"})
	go h1.Watch(FOLDER_NODE_INFO, "server", func(k, v string) {
		fmt.Printf("add key %s, value %s\n", k, v)
	}, func(k, v string) {
		fmt.Printf("delete key %s, value %s\n", k, v)
	}, "prefix")
	time.Sleep(1 * time.Second)
	log.Println("start to regist")
	go h2.Register(FOLDER_NODE_INFO, "server1", "123456")
	time.Sleep(5 * time.Second)
	go h3.Register(FOLDER_NODE_INFO, "server2", "123456")
	time.Sleep(10 * time.Second)
	log.Println("start to stop")
	h2.Stop()
	log.Println("stop the server")
	kvs, _ := h4.Get(FOLDER_NODE_INFO, "server")
	for k, v := range kvs {
		log.Println(k, "\t", v)
	}
	log.Println("done")
	time.Sleep(30 * time.Second)
}
