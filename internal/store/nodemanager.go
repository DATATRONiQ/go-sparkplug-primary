package store

import (
	"sync"
	"time"
)

type NodeManager struct {
	GroupID       string
	NodeID        string
	mu            sync.RWMutex
	LastMessageAt time.Time
	Online        bool
	Devices       map[string]*DeviceManager
}

type FetchedNode struct {
	ID            string          `json:"id"`
	GroupID       string          `json:"groupId"`
	Online        bool            `json:"online"`
	LastMessageAt time.Time       `json:"lastMessageAt"`
	Devices       []FetchedDevice `json:"devices"`
}

func NewNodeManager(groupID, nodeID string) *NodeManager {
	return &NodeManager{
		GroupID:       groupID,
		NodeID:        nodeID,
		LastMessageAt: time.Now(),
		Devices:       make(map[string]*DeviceManager),
	}
}

func (nm *NodeManager) nodeBirth(msg Message) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	// TODO: Check for bdSeq

	if msg.ReceivedAt.After(nm.LastMessageAt) {
		nm.LastMessageAt = msg.ReceivedAt
	}
	nm.Online = true

	// TODO: Add Metrics
}

func (nm *NodeManager) nodeDeath(msg Message) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	// TODO: Check for bdSeq

	if msg.ReceivedAt.After(nm.LastMessageAt) {
		nm.LastMessageAt = msg.ReceivedAt
	}
	nm.Online = false

	for _, device := range nm.Devices {
		device.Offline()
	}

	// TODO: Make metrics stale
}

func (nm *NodeManager) Fetch() *FetchedNode {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	devices := make([]FetchedDevice, 0, len(nm.Devices))
	for _, device := range nm.Devices {
		devices = append(devices, *device.Fetch())
	}

	return &FetchedNode{
		ID:            nm.NodeID,
		GroupID:       nm.GroupID,
		Online:        nm.Online,
		LastMessageAt: nm.LastMessageAt,
		Devices:       devices,
	}
}
