package store

import (
	"sync"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/api"
	"github.com/DATATRONiQ/go-sparkplug-primary/internal/ssehandler"
	"github.com/sirupsen/logrus"
)

type StoreManager struct {
	mu     sync.RWMutex
	Groups map[string]*GroupManager

	GroupsSSEHandler *ssehandler.SSEHandler
}

func NewStoreManager(msgChan <-chan Message) *StoreManager {
	sm := &StoreManager{
		Groups:           make(map[string]*GroupManager),
		GroupsSSEHandler: ssehandler.NewSSEHandler(),
	}

	go sm.start(msgChan)
	return sm
}

func (sm *StoreManager) start(msgChan <-chan Message) {
	for msg := range msgChan {
		sm.ProcessMessage(msg)
	}
}

func (sm *StoreManager) ProcessMessage(msg Message) {
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
		event := groupManager.nodeBirth(msg)
		sm.sendEvent(event)
	case NodeData:
		groupManager, ok := sm.Groups[msg.GroupID]
		if !ok {
			logrus.Debugf("NDATA: Group %s is currently not in store", msg.GroupID)
			return
		}
		event := groupManager.nodeData(msg)
		sm.sendEvent(event)
	case NodeDeath:
		groupManager, ok := sm.Groups[msg.GroupID]
		if !ok {
			logrus.Debugf("NDEATH: Group %s is currently not in store", msg.GroupID)
			return
		}
		event := groupManager.nodeDeath(msg)
		sm.sendEvent(event)
	case DeviceBirth:
		groupManager, ok := sm.Groups[msg.GroupID]
		if !ok {
			logrus.Debugf("DBIRTH: Group %s is currently not in store", msg.GroupID)
			return
		}
		event := groupManager.deviceBirth(msg)
		sm.sendEvent(event)
	case DeviceData:
		groupManager, ok := sm.Groups[msg.GroupID]
		if !ok {
			logrus.Debugf("DDATA: Group %s is currently not in store", msg.GroupID)
			return
		}
		event := groupManager.deviceData(msg)
		sm.sendEvent(event)
	case DeviceDeath:
		groupManager, ok := sm.Groups[msg.GroupID]
		if !ok {
			logrus.Debugf("DDEATH: Group %s is currently not in store", msg.GroupID)
			return
		}
		event := groupManager.deviceDeath(msg)
		sm.sendEvent(event)
	default:
		logrus.Warnf("Unimplemented message type: %s", msg.Type)
	}
}

func (sm *StoreManager) sendEvent(event *api.Event) {
	if event == nil {
		logrus.Debugf("No event to send")
		return
	}
	sm.GroupsSSEHandler.Send(event)
}

func (sm *StoreManager) fetch() *[]api.FullGroup {
	fetchedGroups := make([]api.FullGroup, 0)
	for _, groupManager := range sm.Groups {
		fetchedGroups = append(fetchedGroups, *groupManager.Fetch())
	}
	return &fetchedGroups
}

func (sm *StoreManager) Fetch() *[]api.FullGroup {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.fetch()
}
