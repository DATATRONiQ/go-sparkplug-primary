package api

import "time"

type Event struct {
	Type      string    `json:"type"`
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

type InitialEvent struct {
	Groups []FullGroup `json:"groups"`
}

// Server-Sent-Event when a valid NBIRTH message was received and handled
type NodeBirthEvent struct {
	Node        Node     `json:"node"`
	NodeMetrics []Metric `json:"nodeMetrics"`
}

// Server-Sent-Event when a valid NDATA message was received and handled
type NodeDataEvent struct {
	Node        Node     `json:"node"`
	NodeMetrics []Metric `json:"nodeMetrics"`
}

// Server-Sent-Event when a valid NDEATH message was received and handled
type NodeDeathEvent struct {
	Node Node `json:"node"`
}

// Server-Sent-Event when a valid DBIRTH message was received and handled
type DeviceBirthEvent struct {
	Node          Node     `json:"node"`
	Device        Device   `json:"device"`
	DeviceMetrics []Metric `json:"deviceMetrics"`
}

// Server-Sent-Event when a valid DDATA message was received and handled
type DeviceDataEvent struct {
	Node          Node     `json:"node"`
	Device        Device   `json:"device"`
	DeviceMetrics []Metric `json:"deviceMetrics"`
}

// Server-Sent-Event when a valid DDEATH message was received and handled
type DeviceDeathEvent struct {
	Node   Node   `json:"node"`
	Device Device `json:"device"`
}
