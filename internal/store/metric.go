package store

import (
	"fmt"
	"time"

	"github.com/DATATRONiQ/go-sparkplug-primary/third_party/sparkplugb"
)

type Metric struct {
	Name          string
	Alias         uint64
	DataType      sparkplugb.DataType
	LastTimeStamp *time.Time
	IsNull        bool
	Value         any
}

type FetchedMetric struct {
	Name      string    `json:"name"`
	Alias     uint64    `json:"alias"`
	Stale     bool      `json:"stale"`
	DataType  string    `json:"dataType"`
	Timestamp time.Time `json:"timestamp"`
	IsNull    bool      `json:"isNull"`
	Value     any       `json:"value"`
}

func NewMetric(metric *sparkplugb.Payload_Metric) (*Metric, error) {
	if metric == nil {
		return nil, fmt.Errorf("metric is nil")
	}

	if metric.Alias == nil {
		return nil, fmt.Errorf("metric alias is nil")
	}

	if metric.Name == nil {
		return nil, fmt.Errorf("metric name is nil")
	}

	if metric.Datatype == nil {
		return nil, fmt.Errorf("metric data type is nil")
	}

	newMetric := Metric{
		Name:     *metric.Name,
		Alias:    *metric.Alias,
		DataType: sparkplugb.DataType(*metric.Datatype),
	}

	if metric.Timestamp != nil {
		ts := time.UnixMilli(int64(*metric.Timestamp))
		newMetric.LastTimeStamp = &ts
	}

	if metric.IsNull != nil {
		newMetric.IsNull = *metric.IsNull
	}

	// TODO: add value

	return &newMetric, nil
}

func (m *Metric) Fetch(isStale bool) *FetchedMetric {
	metric := FetchedMetric{
		Name:     m.Name,
		Alias:    m.Alias,
		Stale:    isStale,
		DataType: m.DataType.String(),
		IsNull:   m.IsNull,
	}
	if m.LastTimeStamp != nil {
		metric.Timestamp = *m.LastTimeStamp
	}
	// TODO: add value
	return &metric
}
