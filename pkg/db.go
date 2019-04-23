package db

import (
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

type Handler struct {
	name string
	db   *bolt.DB
}

const (
	DB_NAME = "my.db"
)

var (
	BUCKET_NODE_INFO     = []byte("node")
	BUCKET_POD_INFO      = []byte("pod")
	BUCKET_EXECUTOR_INFO = []byte("executor")
)

func New() *Handler {
	db, err := bolt.Open(DB_NAME, 0600, nil)
	if err != nil {
		glog.Errorf("open db error,%v", err)
	}
	return &Handler{name: "my.db", db: db}
}

// add or update k-v pair to the specify bucket
func (h *Handler) Update(bucket, key, value []byte) error {
	err := h.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			glog.V(4).Infof("create bucket error,%v", err)
			return err
		}
		err = b.Put(key, value)
		if err != nil {
			glog.V(4).Infof("put k-v error,err:%v,key = %s", err, string(key))
		}
		return err
	})
	return err
}

// delete the k-v pair with the specify key
func (h *Handler) Delete(bucket, key []byte) error {
	err := h.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			glog.V(4).Infof("create bucket error,%v", err)
			return err
		}
		err = b.Delete(key)
		if err != nil {
			glog.V(4).Infof("delete k-v error,err:%v,key = %s", err, string(key))
		}
		return err
	})
	return err
}

// if not exist, return nil
func (h *Handler) Get(bucket, key []byte) ([]byte, error) {
	var value []byte
	err := h.db.View(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			glog.V(4).Infof("create bucket error,%v", err)
			return err
		}
		v := b.Get(key)
		copy(v, value)
		return nil
	})
	return value, err
}
func (h *Handler) Close() {
	glog.V(4).Infof("close bolt database")
	h.db.Close()
}
