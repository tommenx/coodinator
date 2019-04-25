package resource

// scheduler send to coordinator when it need to scheduler

type ResourceType int

type ResourceUnit int

const (
	STORAGE ResourceType = iota
	LIMIT
)

const (
	B ResourceUnit = iota
	KB
	MB
	GB
)

type Resource struct {
	Type ResourceType `json:"type"`
	Kind string       `json:"kind"`
	Size uint64       `json:"size"`
	Unit ResourceUnit `json:"unit"`
}

type Device struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Maj        string `json:"maj"`
	Min        string `json:"min"`
	DevicePath string `json:"device_path"`
}

type PVC struct {
}

type Request struct {
}

type Pod struct {
	Node      string `json:"node"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
