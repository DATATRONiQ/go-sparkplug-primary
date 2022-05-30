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

	err := newMetric.addValue(metric)
	if err != nil {
		return nil, err
	}

	return &newMetric, err
}

func (m *Metric) addValue(metric *sparkplugb.Payload_Metric) error {
	if m.IsNull {
		// metric is null so there is no value to add
		return nil
	}
	if metric == nil {
		return fmt.Errorf("metric is nil")
	}
	if metric.Value == nil {
		return fmt.Errorf("metric value is nil")
	}
	switch m.DataType {
	case sparkplugb.DataType_Boolean:
		m.Value = metric.GetBooleanValue()
	case sparkplugb.DataType_Double:
		m.Value = metric.GetDoubleValue()
	case sparkplugb.DataType_Float:
		m.Value = float32(metric.GetFloatValue())
	case sparkplugb.DataType_Int8:
		m.Value = int8(metric.GetIntValue())
	case sparkplugb.DataType_Int16:
		m.Value = int16(metric.GetIntValue())
	case sparkplugb.DataType_Int32:
		m.Value = int32(metric.GetIntValue())
	case sparkplugb.DataType_Int64:
		m.Value = int64(metric.GetLongValue())
	case sparkplugb.DataType_UInt8:
		m.Value = uint8(metric.GetIntValue())
	case sparkplugb.DataType_UInt16:
		m.Value = uint16(metric.GetIntValue())
	case sparkplugb.DataType_UInt32:
		m.Value = uint32(metric.GetIntValue())
	case sparkplugb.DataType_UInt64:
		m.Value = metric.GetLongValue()
	case sparkplugb.DataType_String, sparkplugb.DataType_UUID, sparkplugb.DataType_Text:
		m.Value = metric.GetStringValue()
	default:
		return fmt.Errorf("unsupported data type: %s", m.DataType.String())
	}
	return nil
}

func (m *Metric) Update(metric *sparkplugb.Payload_Metric) error {
	if metric == nil {
		return fmt.Errorf("metric is nil")
	}
	if *metric.Alias != m.Alias {
		return fmt.Errorf("metric alias mismatch")
	}
	if metric.Timestamp != nil {
		ts := time.UnixMilli(int64(*metric.Timestamp))
		m.LastTimeStamp = &ts
	}

	// only when IsNull exists in the payload and its value is true
	newIsNull := metric.IsNull != nil && *metric.IsNull

	if newIsNull {
		m.IsNull = true
		m.Value = nil
		return nil
	}

	m.IsNull = false
	return m.addValue(metric)
}

func (m *Metric) Fetch(isStale bool) *FetchedMetric {
	metric := FetchedMetric{
		Name:     m.Name,
		Alias:    m.Alias,
		Stale:    isStale,
		DataType: m.DataType.String(),
		IsNull:   m.IsNull,
		Value:    m.Value,
	}
	if m.LastTimeStamp != nil {
		metric.Timestamp = *m.LastTimeStamp
	}
	// TODO: add value
	return &metric
}
