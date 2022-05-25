package api

import "time"

// The data structure returned by the Fetch() method
type Group struct {
	ID            string    `json:"id"`            // The group ID
	LastMessageAt time.Time `json:"lastMessageAt"` // The last time a message was received regarding this group
}

type FullGroup struct {
	Group
	Nodes []FullNode `json:"nodes"` // The state of the nodes in the group
}

// The data structure returned by the Fetch() method
type Node struct {
	ID            string    `json:"id"`            // The node ID
	GroupID       string    `json:"groupId"`       // The group ID
	Online        bool      `json:"online"`        // Whether the node is online
	LastMessageAt time.Time `json:"lastMessageAt"` // The last time a message was received regarding this node
}

type FullNode struct {
	Node
	Devices []FullDevice `json:"devices"` // The state of the devices
	Metrics []Metric     `json:"metrics"` // The metrics of this node
}

// The data structure returned by the Fetch() method
type Device struct {
	ID            string    `json:"id"`            // The device ID
	NodeID        string    `json:"nodeId"`        // The node ID
	GroupID       string    `json:"groupId"`       // The group ID
	Online        bool      `json:"online"`        // Whether the device is online
	LastMessageAt time.Time `json:"lastMessageAt"` // The last time a message was received regarding this device
}

type FullDevice struct {
	Device
	Metrics []Metric `json:"metrics"` // The metrics of this device
}

type Metric struct {
	Name      string    `json:"name"`
	Alias     uint64    `json:"alias"`
	Stale     bool      `json:"stale"`
	DataType  string    `json:"dataType"`
	Timestamp time.Time `json:"timestamp"`
	IsNull    bool      `json:"isNull"`
	Value     any       `json:"value"`
}
