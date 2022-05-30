package store

import (
	"sync"
	"time"

	"github.com/DATATRONiQ/go-sparkplug-primary/third_party/sparkplugb"
)

// The sparkplug-B message type for node and device messages
type Type string

const (
	NodeBirth     Type = "NBIRTH"
	NodeDeath     Type = "NDEATH"
	NodeData      Type = "NDATA"
	NodeCommand   Type = "NCMD"
	DeviceBirth   Type = "DBIRTH"
	DeviceDeath   Type = "DDEATH"
	DeviceData    Type = "DDATA"
	DeviceCommand Type = "DCMD"
)

// Represents a sparkplug message
type Message struct {
	ReceivedAt time.Time
	GroupID    string
	NodeID     string
	Type       Type
	DeviceID   string
	Payload    *sparkplugb.Payload
}

// The data structure returned by the Fetch() method
type FetchedMessage struct {
	GroupID      string    `json:"groupId"`      // The group ID
	NodeID       string    `json:"nodeId"`       // The node ID
	DeviceID     string    `json:"deviceId"`     // The device ID
	Type         Type      `json:"type"`         // The message type
	MetricAmount int       `json:"metricAmount"` // The amount of metrics in the message
	ReceivedAt   time.Time `json:"receivedAt"`   // The time the message was received
}

// basically our in-memory database
var msgLog = make([]Message, 0)
var msgLogMutex sync.RWMutex

func addMessage(msg Message) {
	msgLogMutex.Lock()
	defer msgLogMutex.Unlock()
	msgLog = append(msgLog, msg)
}

// Returns all messages received since the start of the application
func Fetch() *[]FetchedMessage {
	msgLogMutex.RLock()
	defer msgLogMutex.RUnlock()

	fetchedMsgs := make([]FetchedMessage, 0)
	for _, msg := range msgLog {
		metricAmount := 0
		if msg.Payload != nil {
			metricAmount = len(msg.Payload.Metrics)
		}

		fetchedMsgs = append(fetchedMsgs, FetchedMessage{
			GroupID:      msg.GroupID,
			NodeID:       msg.NodeID,
			DeviceID:     msg.DeviceID,
			Type:         msg.Type,
			MetricAmount: metricAmount,
			ReceivedAt:   msg.ReceivedAt,
		})
	}

	return &fetchedMsgs
}
