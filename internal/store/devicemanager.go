package store

import (
	"sync"
	"time"
)

// Manages the state of a single sparkplug device
type DeviceManager struct {
	GroupID       string    // The group this device belongs to
	NodeID        string    // The node this device belongs to
	DeviceID      string    // The device ID
	Online        bool      // Whether the device is online
	LastMessageAt time.Time // The last time a message was received regarding this device

	mu sync.RWMutex
}

// The data structure returned by the Fetch() method
type FetchedDevice struct {
	ID            string    `json:"id"`            // The device ID
	NodeID        string    `json:"nodeId"`        // The node ID
	GroupID       string    `json:"groupId"`       // The group ID
	Online        bool      `json:"online"`        // Whether the device is online
	LastMessageAt time.Time `json:"lastMessageAt"` // The last time a message was received regarding this device
}

// Creates a new DeviceManager for the given device
func NewDeviceManager(groupID, nodeID, deviceID string) *DeviceManager {
	return &DeviceManager{
		GroupID:       groupID,
		NodeID:        nodeID,
		DeviceID:      deviceID,
		LastMessageAt: time.Now(),
	}
}

func (dm *DeviceManager) deviceBirth(msg Message) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if msg.ReceivedAt.After(dm.LastMessageAt) {
		dm.LastMessageAt = msg.ReceivedAt
	}
	dm.Online = true

	// TODO: Add Metrics
}

func (dm *DeviceManager) deviceDeath(msg Message) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if msg.ReceivedAt.After(dm.LastMessageAt) {
		dm.LastMessageAt = msg.ReceivedAt
	}
	dm.Online = false

	// TODO: Make metrics stale
}

func (dm *DeviceManager) offline() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.Online = false

	// TODO: Make metrics stale
}

// Returns the current state of the device
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
