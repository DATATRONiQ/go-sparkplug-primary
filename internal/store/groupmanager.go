package store

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type GroupManager struct {
	GroupID       string
	mu            sync.RWMutex
	LastMessageAt time.Time
	Nodes         map[string]*NodeManager
}

type FetchedGroup struct {
	ID            string        `json:"id"`
	LastMessageAt time.Time     `json:"lastMessageAt"`
	Nodes         []FetchedNode `json:"nodes"`
}

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
