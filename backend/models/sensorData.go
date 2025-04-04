// backend/models/thresholds.go
package models

import (
	"fmt"
	"time"
)

type SensorData struct {
	ID        int       `db:"id" json:"id"`
	Value     float64   `db:"value" json:"value"`
	Type      string    `db:"type" json:"type"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
}

type SensorType int

// Константы SensorType
const (
	Temperature SensorType = iota // 0 (начинаем с большой буквы для экспорта)
	Humidity                      // 1
	Pressure                      // 2
)

// String() преобразует SensorType в строку
func (s SensorType) String() string {
	switch s {
	case Temperature:
		return "temperature"
	case Humidity:
		return "humidity"
	case Pressure:
		return "pressure"
	default:
		return "unknown"
	}
}

// ParseSensorType преобразует строку в SensorType
func ParseSensorType(str string) (SensorType, error) {
	switch str {
	case "temperature":
		return Temperature, nil
	case "humidity":
		return Humidity, nil
	case "pressure":
		return Pressure, nil
	default:
		return -1, fmt.Errorf("unknown sensor type: %s", str)
	}
}
