package store

import (
	"sync"
	"time"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/util"
	"github.com/sirupsen/logrus"
)

// Manages the state of a single sparkplug device
type DeviceManager struct {
	GroupID       string             // The group this device belongs to
	NodeID        string             // The node this device belongs to
	DeviceID      string             // The device ID
	Online        bool               // Whether the device is online
	LastMessageAt time.Time          // The last time a message was received regarding this device
	Metrics       map[uint64]*Metric // The metrics of this device (Alias -> Metric)

	mu sync.RWMutex
}

// The data structure returned by the Fetch() method
type FetchedDevice struct {
	ID            string          `json:"id"`            // The device ID
	NodeID        string          `json:"nodeId"`        // The node ID
	GroupID       string          `json:"groupId"`       // The group ID
	Online        bool            `json:"online"`        // Whether the device is online
	LastMessageAt time.Time       `json:"lastMessageAt"` // The last time a message was received regarding this device
	Metrics       []FetchedMetric `json:"metrics"`       // The metrics of this device
}

// Creates a new DeviceManager for the given device
func NewDeviceManager(groupID, nodeID, deviceID string) *DeviceManager {
	return &DeviceManager{
		GroupID:       groupID,
		NodeID:        nodeID,
		DeviceID:      deviceID,
		LastMessageAt: time.Now(),
		Metrics:       make(map[uint64]*Metric),
	}
}

func (dm *DeviceManager) deviceBirth(msg Message) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if msg.ReceivedAt.After(dm.LastMessageAt) {
		dm.LastMessageAt = msg.ReceivedAt
	}
	dm.Online = true

	dm.Metrics = make(map[uint64]*Metric)

	for _, metric := range msg.Payload.Metrics {
		alias := metric.Alias
		if alias == nil {
			if metric.Name == nil {
				logrus.Warnf("DBIRTH: Device %s has no alias for metric with nil name", dm.DeviceID)
			} else {
				logrus.Warnf("DBIRTH: Device %s has no alias for metric %s", dm.DeviceID, *metric.Name)
			}
			continue
		}

		newMetric, err := NewMetric(metric)
		if err != nil {
			logrus.Warnf("DBIRTH: Device %s has an invalid metric %d: %s", dm.DeviceID, metric.Name, err)
			continue
		}
		dm.Metrics[*alias] = newMetric
	}
}

func (dm *DeviceManager) deviceDeath(msg Message) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if msg.ReceivedAt.After(dm.LastMessageAt) {
		dm.LastMessageAt = msg.ReceivedAt
	}
	dm.Online = false
}

func (dm *DeviceManager) offline() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.Online = false
}

// Returns the current state of the device
func (dm *DeviceManager) Fetch() *FetchedDevice {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	sortedAliases := util.SortedKeys(dm.Metrics)
	metrics := make([]FetchedMetric, 0, len(dm.Metrics))
	for _, alias := range sortedAliases {
		fetchedMetric := dm.Metrics[alias].Fetch(!dm.Online)
		metrics = append(metrics, *fetchedMetric)
	}

	return &FetchedDevice{
		ID:            dm.DeviceID,
		NodeID:        dm.NodeID,
		GroupID:       dm.GroupID,
		Online:        dm.Online,
		LastMessageAt: dm.LastMessageAt,
		Metrics:       metrics,
	}
}
