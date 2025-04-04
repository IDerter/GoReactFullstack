// backend/models/thresholds.go
package models

import "time"

type Threshold struct {
	ID        int       `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"`
	MinValue  float64   `json:"min_value" db:"min_value"`
	MaxValue  float64   `json:"max_value" db:"max_value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
