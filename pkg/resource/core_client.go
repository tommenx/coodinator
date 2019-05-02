package resource

import (
	"github.com/tommenx/coordinator/pkg/db"
)

type coreClient struct {
	hanlder *db.EtcdHandler
}

func NewCoreClient() (*coreClient, error) {
	endpoints := []string{"http://127.0.0.1:2379"}
	handler, err := db.NewEtcdHandler(endpoints)
	if err != nil {
		return nil, err
	}
	return &coreClient{hanlder: handler}, nil
}

func (c *coreClient) Node() NodeInfoInterface {
	return newNodeInfo(c.hanlder)
}

func (c *coreClient) Pod(ns string) PodInfoInterface {
	return NewPodInfo(c.hanlder, ns)
}
