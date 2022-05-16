package store

import (
	"time"

	"github.com/DATATRONiQ/go-sparkplug-primary/third_party/sparkplugb"
)

type MessageType string

const (
	NodeBirth     MessageType = "NBIRTH"
	NodeDeath     MessageType = "NDEATH"
	NodeData      MessageType = "NDATA"
	NodeCommand   MessageType = "NCMD"
	DeviceBirth   MessageType = "DBIRTH"
	DeviceDeath   MessageType = "DDEATH"
	DeviceData    MessageType = "DDATA"
	DeviceCommand MessageType = "DCMD"
)

// Represents a sparkplug message
type Message struct {
	ReceivedAt  time.Time           `json:"receivedAt"`
	GroupID     string              `json:"groupId"`
	NodeID      string              `json:"nodeId"`
	MessageType MessageType         `json:"messageType"`
	DeviceID    string              `json:"deviceId"`
	Payload     *sparkplugb.Payload `json:"payload"`
}

// basically our in-memory database
var msgLog = make([]Message, 0)

func AddMessage(msg Message) error {
	msgLog = append(msgLog, msg)
	return nil
}

func FetchMessages() []Message {
	return msgLog
}
