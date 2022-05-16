package store

import "time"

// Represents a sparkplug message
type Message struct {
	ReceivedAt  time.Time `json:"receivedAt"`
	GroupID     string    `json:"groupId"`
	NodeID      string    `json:"nodeId"`
	MessageType string    `json:"messageType"`
	DeviceID    string    `json:"deviceId"`
	Payload     any       `json:"payload"`
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
