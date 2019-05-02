package resource

import (
	"github.com/tommenx/coordinator/pkg/db"
)

type AddFunc func(k, v string)
type DelFunc func(k, v string)

type ExecutorInterafce interface {
	Register(*Executor) error
	Watch(name string, put AddFunc, del DelFunc) error
	Stop()
}

type executor struct {
	h *db.EtcdHandler
}
