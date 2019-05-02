package resource

import "errors"

// scheduler send to coordinator when it need to scheduler

type ResourceType int

type ResourceUnit int

type StorageLevel int

const (
	STORAGE ResourceType = iota
	LIMIT
)

const (
	B ResourceUnit = iota
	KB
	MB
	GB
	C
)

const (
	HDD StorageLevel = iota
	SSD
	NVM
)

type Resource struct {
	Type ResourceType `json:"type"`
	Kind string       `json:"kind"`
	Size uint64       `json:"size"`
	Unit ResourceUnit `json:"unit"`
}

type Storage struct {
	Name      string       `json:"name"` // vg
	Level     StorageLevel `json:"level"`
	Resources []Resource   `json:"resources"`
}

type Device struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Maj         string `json:"maj"`
	Min         string `json:"min"`
	VolumeGroup string `json:"volume_group"` //来源于Storage中的name
	DevicePath  string `json:"device_path"`
}

type PVC struct {
	Name      string     `json:"name"`
	Resources []Resource `json:"resources"`
}

type PV struct {
	Name   string `json:"name"`
	Device Device `json:"device"`
}

type Allocation struct {
	PVC PVC `json:"pvc"`
	PV  PV  `json:"pv"`
}

type Pod struct {
	Node        string       `json:"node"`
	Name        string       `json:"name"`
	Namespace   string       `json:"namespace"`
	Allocations []Allocation `json:"allocation"`
}

type Node struct {
	Name     string    `json:"name"`
	Storages []Storage `json:"storages"`
}

type Executor struct {
	Hostname string `json:"host_name"`
	Address  string `json:"address"`
}

var (
	ErrKeyNotExist     = errors.New("key not exist")
	ErrKeyAlreadyExist = errors.New("key already exist")
)
