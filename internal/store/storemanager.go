package store

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type StoreManager struct {
	mu     sync.RWMutex
	Groups map[string]*GroupManager
}

func NewStoreManager(msgChan <-chan Message) *StoreManager {
	sm := &StoreManager{
		Groups: make(map[string]*GroupManager),
	}

	go sm.start(msgChan)
	return sm
}

func (sm *StoreManager) start(msgChan <-chan Message) {
	for msg := range msgChan {
		sm.processMessage(msg)
	}
}

func (sm *StoreManager) processMessage(msg Message) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	addMessage(msg)
	switch msg.Type {
	case NodeBirth:
		groupManager, ok := sm.Groups[msg.GroupID]
		if !ok {
			sm.Groups[msg.GroupID] = NewGroupManager(msg.GroupID)
			groupManager = sm.Groups[msg.GroupID]
		}
		groupManager.nodeBirth(msg)
	case NodeData:
		groupManager, ok := sm.Groups[msg.GroupID]
		if !ok {
			logrus.Debugf("NDATA: Group %s is currently not in store", msg.GroupID)
			return
		}
		groupManager.nodeData(msg)
	case NodeDeath:
		groupManager, ok := sm.Groups[msg.GroupID]
		if !ok {
			logrus.Debugf("NDEATH: Group %s is currently not in store", msg.GroupID)
			return
		}
		groupManager.nodeDeath(msg)
	case DeviceBirth:
		groupManager, ok := sm.Groups[msg.GroupID]
		if !ok {
			logrus.Debugf("DBIRTH: Group %s is currently not in store", msg.GroupID)
			return
		}
		groupManager.deviceBirth(msg)
	case DeviceData:
		groupManager, ok := sm.Groups[msg.GroupID]
		if !ok {
			logrus.Debugf("DDATA: Group %s is currently not in store", msg.GroupID)
			return
		}
		groupManager.deviceData(msg)
	case DeviceDeath:
		groupManager, ok := sm.Groups[msg.GroupID]
		if !ok {
			logrus.Debugf("DDEATH: Group %s is currently not in store", msg.GroupID)
			return
		}
		groupManager.deviceDeath(msg)
	default:
		logrus.Warnf("Unimplemented message type: %s", msg.Type)
	}
}

func (sm *StoreManager) Fetch() *[]FetchedGroup {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	fetchedGroups := make([]FetchedGroup, 0)
	for _, groupManager := range sm.Groups {
		fetchedGroups = append(fetchedGroups, *groupManager.Fetch())
	}
	return &fetchedGroups
}
