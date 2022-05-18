package store

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Manages the state of a single sparkplug group
type GroupManager struct {
	GroupID       string                  // The group ID
	LastMessageAt time.Time               // The last time a message was received regarding this group
	Nodes         map[string]*NodeManager // The node managers for each node in the group

	mu sync.RWMutex
}

// The data structure returned by the Fetch() method
type FetchedGroup struct {
	ID            string        `json:"id"`            // The group ID
	LastMessageAt time.Time     `json:"lastMessageAt"` // The last time a message was received regarding this group
	Nodes         []FetchedNode `json:"nodes"`         // The state of the nodes in the group
}

// Creates a new group manager for the given group ID
func NewGroupManager(groupID string) *GroupManager {
	return &GroupManager{
		GroupID:       groupID,
		LastMessageAt: time.Now(),
		Nodes:         make(map[string]*NodeManager),
	}
}

func (gm *GroupManager) nodeBirth(msg Message) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	nodeManager, ok := gm.Nodes[msg.NodeID]
	if !ok {
		gm.Nodes[msg.NodeID] = NewNodeManager(gm.GroupID, msg.NodeID)
		nodeManager = gm.Nodes[msg.NodeID]
	}

	if msg.ReceivedAt.After(gm.LastMessageAt) {
		gm.LastMessageAt = msg.ReceivedAt
	}
	nodeManager.nodeBirth(msg)
}

func (gm *GroupManager) nodeDeath(msg Message) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	nodeManager, ok := gm.Nodes[msg.NodeID]
	if !ok {
		logrus.Debugf("NDEATH: Node %s is currently not in group %s", msg.NodeID, gm.GroupID)
		return
	}

	if msg.ReceivedAt.After(gm.LastMessageAt) {
		gm.LastMessageAt = msg.ReceivedAt
	}
	nodeManager.nodeDeath(msg)
}

func (gm *GroupManager) deviceBirth(msg Message) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	nodeManager, ok := gm.Nodes[msg.NodeID]
	if !ok {
		logrus.Debugf("DBIRTH: Node %s is currently not in group %s", msg.NodeID, gm.GroupID)
		return
	}

	if msg.ReceivedAt.After(gm.LastMessageAt) {
		gm.LastMessageAt = msg.ReceivedAt
	}
	nodeManager.deviceBirth(msg)
}

func (gm *GroupManager) deviceDeath(msg Message) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	nodeManager, ok := gm.Nodes[msg.NodeID]
	if !ok {
		logrus.Debugf("DDEATH: Node %s is currently not in group %s", msg.NodeID, gm.GroupID)
		return
	}

	if msg.ReceivedAt.After(gm.LastMessageAt) {
		gm.LastMessageAt = msg.ReceivedAt
	}
	nodeManager.deviceDeath(msg)
}

// Returns the current state of the group and its nodes
func (gm *GroupManager) Fetch() *FetchedGroup {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	nodes := make([]FetchedNode, 0, len(gm.Nodes))
	for _, nodeManager := range gm.Nodes {
		nodes = append(nodes, *nodeManager.Fetch())
	}

	return &FetchedGroup{
		ID:            gm.GroupID,
		LastMessageAt: gm.LastMessageAt,
		Nodes:         nodes,
	}
}
