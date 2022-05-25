package store

import (
	"sync"
	"time"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/api"
	"github.com/sirupsen/logrus"
)

// Manages the state of a single sparkplug device
type DeviceManager struct {
	MetricContainer
	GroupID       string    // The group this device belongs to
	NodeID        string    // The node this device belongs to
	DeviceID      string    // The device ID
	Online        bool      // Whether the device is online
	LastMessageAt time.Time // The last time a message was received regarding this device

	mu sync.RWMutex
}

// Creates a new DeviceManager for the given device
func NewDeviceManager(groupID, nodeID, deviceID string) *DeviceManager {
	return &DeviceManager{
		MetricContainer: *NewMetricContainer(),
		GroupID:         groupID,
		NodeID:          nodeID,
		DeviceID:        deviceID,
		LastMessageAt:   time.Now(),
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

func (dm *DeviceManager) deviceData(msg Message) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if msg.Payload == nil {
		logrus.Warnf("DDATA: Device %s got message with nil payload", dm.DeviceID)
		return
	}

	if msg.Payload.Metrics == nil || len(msg.Payload.Metrics) == 0 {
		logrus.Warnf("DDATA: Device %s got message with no metrics", dm.DeviceID)
		return
	}

	if msg.ReceivedAt.After(dm.LastMessageAt) {
		dm.LastMessageAt = msg.ReceivedAt
	}

	for _, metric := range msg.Payload.Metrics {
		alias := metric.Alias
		if alias == nil {
			if metric.Name == nil {
				logrus.Warnf("DDATA: Device %s got metric with nil name and alias", dm.DeviceID)
			} else {
				logrus.Warnf("DDATA: Device %s got metric with nil alias and name: %s", dm.DeviceID, *metric.Name)
			}
			continue
		}

		currMetric, ok := dm.Metrics[*alias]
		if !ok {
			logrus.Warnf("DDATA: Device %s got metric with unknown alias %d", dm.DeviceID, *alias)
			continue
		}

		err := currMetric.Update(metric)
		if err != nil {
			logrus.Warnf("DDATA: Device %s got an invalid metric with name %s: %v", dm.DeviceID, currMetric.Name, err)
		}
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
func (dm *DeviceManager) Fetch() *api.FullDevice {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	return &api.FullDevice{
		Device: api.Device{
			ID:            dm.DeviceID,
			NodeID:        dm.NodeID,
			GroupID:       dm.GroupID,
			Online:        dm.Online,
			LastMessageAt: dm.LastMessageAt,
		},
		Metrics: *dm.fetchMetrics(!dm.Online),
	}
}
