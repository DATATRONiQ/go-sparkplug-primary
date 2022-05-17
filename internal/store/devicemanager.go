package store

import (
	"sync"
	"time"
)

type DeviceManager struct {
	GroupID       string
	NodeID        string
	DeviceID      string
	mu            sync.RWMutex
	Online        bool
	LastMessageAt time.Time
}

type FetchedDevice struct {
	ID            string    `json:"id"`
	NodeID        string    `json:"nodeId"`
	GroupID       string    `json:"groupId"`
	Online        bool      `json:"online"`
	LastMessageAt time.Time `json:"lastMessageAt"`
}

func NewDeviceManager(groupID, nodeID, deviceID string) *DeviceManager {
	return &DeviceManager{
		GroupID:       groupID,
		NodeID:        nodeID,
		DeviceID:      deviceID,
		LastMessageAt: time.Now(),
	}
}

func (dm *DeviceManager) Offline() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.Online = false

	// TODO: Make metrics stale
}

func (dm *DeviceManager) Fetch() *FetchedDevice {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	return &FetchedDevice{
		ID:            dm.DeviceID,
		NodeID:        dm.NodeID,
		GroupID:       dm.GroupID,
		Online:        dm.Online,
		LastMessageAt: dm.LastMessageAt,
	}
}
