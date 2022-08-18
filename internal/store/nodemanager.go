package store

import (
	"sync"
	"time"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/api"
	"github.com/DATATRONiQ/go-sparkplug-primary/internal/util"
	"github.com/sirupsen/logrus"
)

// Manages the state of a single sparkplug EoN-Node
type NodeManager struct {
	MetricContainer
	GroupID       string                    // The group this node belongs to
	NodeID        string                    // The node ID
	Online        bool                      // Whether the node is online
	LastMessageAt time.Time                 // The last time a message was received regarding this node
	Devices       map[string]*DeviceManager // The device managers for each device of this node (DeviceID -> DeviceManager)

	mu sync.RWMutex
}

// Creates a new NodeManager for the given node
func NewNodeManager(groupID, nodeID string) *NodeManager {
	return &NodeManager{
		MetricContainer: *NewMetricContainer(),
		GroupID:         groupID,
		NodeID:          nodeID,
		LastMessageAt:   time.Now(),
		Devices:         make(map[string]*DeviceManager),
	}
}

func (nm *NodeManager) nodeBirth(msg Message) *api.Event {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if msg.Payload == nil {
		logrus.Warnf("NBIRTH: Node %s got message with nil payload", nm.NodeID)
		return nil
	}

	if msg.Payload.Metrics == nil || len(msg.Payload.Metrics) == 0 {
		logrus.Warnf("NBIRTH: Node %s got message with no metrics (not even bdSeq)", nm.NodeID)
		return nil
	}

	// TODO: Check for bdSeq

	if msg.ReceivedAt.After(nm.LastMessageAt) {
		nm.LastMessageAt = msg.ReceivedAt
	}
	nm.Online = true

	nm.Metrics = make(map[uint64]*Metric)

	for _, metric := range msg.Payload.Metrics {
		alias := metric.Alias
		if alias == nil {
			if metric.Name == nil {
				logrus.Warnf("NBIRTH: Node %s got metric with nil name and alias", nm.NodeID)
			} else {
				logrus.Warnf("NBIRTH: Node %s got metric with nil alias and name: %s", nm.NodeID, *metric.Name)
			}
			continue
		}

		newMetric, err := NewMetric(metric)
		if err != nil {
			if metric.Name == nil {
				logrus.Warnf("NBIRTH: Node %s got an invalid metric with alias %d and name %s: %v", nm.NodeID, *metric.Alias, *metric.Name, err)
			} else {
				logrus.Warnf("NBIRTH: Node %s got an invalid metric with alias %d: %v", nm.NodeID, *metric.Name, err)
			}
			continue
		}
		nm.Metrics[*alias] = newMetric
	}

	return &api.Event{
		Type:      string(NodeBirth),
		Timestamp: nm.LastMessageAt,
		Data: api.NodeBirthEvent{
			Node:        *nm.toApiNode(),
			NodeMetrics: *nm.getMetrics(),
		},
	}
}

func (nm *NodeManager) nodeData(msg Message) *api.Event {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	// TODO: Check seq number

	if msg.Payload == nil {
		logrus.Warnf("NDATA: Node %s got message with nil payload", nm.NodeID)
		return nil
	}

	if msg.Payload.Metrics == nil || len(msg.Payload.Metrics) == 0 {
		logrus.Warnf("NDATA: Node %s got message with no metrics", nm.NodeID)
		return nil
	}

	if msg.ReceivedAt.After(nm.LastMessageAt) {
		nm.LastMessageAt = msg.ReceivedAt
	}

	for _, metric := range msg.Payload.Metrics {
		alias := metric.Alias
		if alias == nil {
			if metric.Name == nil {
				logrus.Warnf("NDATA: Node %s got metric with nil name and alias", nm.NodeID)
			} else {
				logrus.Warnf("NDATA: Node %s got metric with nil alias and name: %s", nm.NodeID, *metric.Name)
			}
			continue
		}

		currMetric, ok := nm.Metrics[*alias]
		if !ok {
			logrus.Warnf("NDATA: Node %s got metric with unknown alias %d", nm.NodeID, *alias)
			continue
		}

		err := currMetric.Update(metric)
		if err != nil {
			logrus.Warnf("NDATA: Node %s got an invalid metric with name %s: %v", nm.NodeID, currMetric.Name, err)
		}
	}

	return &api.Event{
		Type:      string(NodeData),
		Timestamp: nm.LastMessageAt,
		Data: api.NodeDataEvent{
			Node:        *nm.toApiNode(),
			NodeMetrics: *nm.getMetrics(),
		},
	}
}

func (nm *NodeManager) nodeDeath(msg Message) *api.Event {
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

	return &api.Event{
		Type:      string(NodeDeath),
		Timestamp: nm.LastMessageAt,
		Data: &api.NodeDeathEvent{
			Node: *nm.toApiNode(),
		},
	}
}

func (nm *NodeManager) deviceBirth(msg Message) *api.Event {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	deviceManager, ok := nm.Devices[msg.DeviceID]
	if !ok {
		nm.Devices[msg.DeviceID] = NewDeviceManager(nm.GroupID, nm.NodeID, msg.DeviceID)
		deviceManager = nm.Devices[msg.DeviceID]
	}

	fullDevice := deviceManager.deviceBirth(msg)

	if fullDevice == nil {
		return nil
	}

	if msg.ReceivedAt.After(nm.LastMessageAt) {
		nm.LastMessageAt = msg.ReceivedAt
	}

	return &api.Event{
		Type:      string(DeviceBirth),
		Timestamp: nm.LastMessageAt,
		Data: api.DeviceBirthEvent{
			Node:          *nm.toApiNode(),
			Device:        fullDevice.Device,
			DeviceMetrics: fullDevice.Metrics,
		},
	}
}

func (nm *NodeManager) deviceData(msg Message) *api.Event {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	deviceManager, ok := nm.Devices[msg.DeviceID]
	if !ok {
		logrus.Debugf("DDATA: Device %s is currently not in node %s", msg.NodeID, nm.NodeID)
		return nil
	}

	if msg.ReceivedAt.After(nm.LastMessageAt) {
		nm.LastMessageAt = msg.ReceivedAt
	}
	fullDevice := deviceManager.deviceData(msg)
	if fullDevice == nil {
		return nil
	}

	return &api.Event{
		Type:      string(DeviceData),
		Timestamp: nm.LastMessageAt,
		Data: api.DeviceDataEvent{
			Node:          *nm.toApiNode(),
			Device:        fullDevice.Device,
			DeviceMetrics: fullDevice.Metrics,
		},
	}
}

func (nm *NodeManager) deviceDeath(msg Message) *api.Event {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	deviceManager, ok := nm.Devices[msg.DeviceID]
	if !ok {
		logrus.Debugf("DDEATH: Device %s is currently not in node %s", msg.NodeID, nm.NodeID)
		return nil
	}

	if msg.ReceivedAt.After(nm.LastMessageAt) {
		nm.LastMessageAt = msg.ReceivedAt
	}
	device := deviceManager.deviceDeath(msg)

	if device == nil {
		return nil
	}

	return &api.Event{
		Type:      string(DeviceDeath),
		Timestamp: nm.LastMessageAt,
		Data: api.DeviceDeathEvent{
			Node:   *nm.toApiNode(),
			Device: *device,
		},
	}
}

func (nm *NodeManager) getMetrics() *[]api.Metric {
	sortedAliases := util.SortedKeys(nm.Metrics)
	return util.MapSlice(sortedAliases, func(alias uint64) api.Metric {
		return *nm.Metrics[alias].Fetch(!nm.Online)
	})
}

func (nm *NodeManager) toApiNode() *api.Node {
	return &api.Node{
		ID:            nm.NodeID,
		GroupID:       nm.GroupID,
		Online:        nm.Online,
		LastMessageAt: nm.LastMessageAt,
	}
}

// Returns the current state of the node and its devices
func (nm *NodeManager) FetchFull() *api.FullNode {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	sortedDeviceIDs := util.SortedKeys(nm.Devices)
	devices := util.MapSlice(sortedDeviceIDs, func(deviceID string) api.FullDevice {
		return *nm.Devices[deviceID].FetchFull()
	})

	return &api.FullNode{
		Node:    *nm.toApiNode(),
		Devices: *devices,
		Metrics: *nm.fetchMetrics(!nm.Online),
	}
}
