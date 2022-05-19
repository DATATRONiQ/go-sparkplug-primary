package store

import (
	"sync"
	"time"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/util"
	"github.com/sirupsen/logrus"
)

// Manages the state of a single sparkplug EoN-Node
type NodeManager struct {
	GroupID       string                    // The group this node belongs to
	NodeID        string                    // The node ID
	Online        bool                      // Whether the node is online
	LastMessageAt time.Time                 // The last time a message was received regarding this node
	Devices       map[string]*DeviceManager // The device managers for each device of this node (DeviceID -> DeviceManager)
	Metrics       map[uint64]*Metric        // The metrics of this node (Alias -> Metric)

	mu sync.RWMutex
}

// The data structure returned by the Fetch() method
type FetchedNode struct {
	ID            string          `json:"id"`            // The node ID
	GroupID       string          `json:"groupId"`       // The group ID
	Online        bool            `json:"online"`        // Whether the node is online
	LastMessageAt time.Time       `json:"lastMessageAt"` // The last time a message was received regarding this node
	Devices       []FetchedDevice `json:"devices"`       // The state of the devices
	Metrics       []FetchedMetric `json:"metrics"`       // The metrics of this node
}

// Creates a new NodeManager for the given node
func NewNodeManager(groupID, nodeID string) *NodeManager {
	return &NodeManager{
		GroupID:       groupID,
		NodeID:        nodeID,
		LastMessageAt: time.Now(),
		Devices:       make(map[string]*DeviceManager),
		Metrics:       make(map[uint64]*Metric),
	}
}

func (nm *NodeManager) nodeBirth(msg Message) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if msg.Payload == nil {
		logrus.Warnf("NBIRTH: Node %s has no payload", nm.NodeID)
		return
	}

	if msg.Payload.Metrics == nil {
		logrus.Warnf("NBIRTH: Node %s has no metrics", nm.NodeID)
		return
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
				logrus.Warnf("NBIRTH: Node %s has no alias for metric with nil name", nm.NodeID)
			} else {
				logrus.Warnf("NBIRTH: Node %s has no alias for metric %s", nm.NodeID, *metric.Name)
			}
			continue
		}

		newMetric, err := NewMetric(metric)
		if err != nil {
			logrus.Warnf("NBIRTH: Node %s has an invalid metric %d: %s", nm.NodeID, metric.Name, err)
			continue
		}
		nm.Metrics[*alias] = newMetric
	}
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
		device.offline()
	}
}

func (nm *NodeManager) deviceBirth(msg Message) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	deviceManager, ok := nm.Devices[msg.DeviceID]
	if !ok {
		nm.Devices[msg.DeviceID] = NewDeviceManager(nm.GroupID, nm.NodeID, msg.DeviceID)
		deviceManager = nm.Devices[msg.DeviceID]
	}

	if msg.ReceivedAt.After(nm.LastMessageAt) {
		nm.LastMessageAt = msg.ReceivedAt
	}
	deviceManager.deviceBirth(msg)
}

func (nm *NodeManager) deviceDeath(msg Message) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	deviceManager, ok := nm.Devices[msg.DeviceID]
	if !ok {
		logrus.Debugf("DDEATH: Device %s is currently not in node %s", msg.NodeID, nm.NodeID)
		return
	}

	if msg.ReceivedAt.After(nm.LastMessageAt) {
		nm.LastMessageAt = msg.ReceivedAt
	}
	deviceManager.deviceDeath(msg)
}

// Returns the current state of the node and its devices
func (nm *NodeManager) Fetch() *FetchedNode {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	sortedDeviceIDs := util.SortedKeys(nm.Devices)
	devices := make([]FetchedDevice, 0, len(nm.Devices))
	for _, deviceID := range sortedDeviceIDs {
		fetchedDevice := nm.Devices[deviceID].Fetch()
		devices = append(devices, *fetchedDevice)
	}

	sortedAliases := util.SortedKeys(nm.Metrics)
	metrics := make([]FetchedMetric, 0, len(nm.Metrics))
	for _, alias := range sortedAliases {
		fetchedMetric := nm.Metrics[alias].Fetch(!nm.Online)
		metrics = append(metrics, *fetchedMetric)
	}

	return &FetchedNode{
		ID:            nm.NodeID,
		GroupID:       nm.GroupID,
		Online:        nm.Online,
		LastMessageAt: nm.LastMessageAt,
		Devices:       devices,
		Metrics:       metrics,
	}
}
