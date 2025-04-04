package db

import (
	"github.com/jmoiron/sqlx"
)

// InitDefaultThresholds устанавливает начальные пороговые значения
func InitDefaultThresholds(db *sqlx.DB) error {
	defaults := []struct {
		Type     string
		MinValue float64
		MaxValue float64
	}{
		{"temperature", 20, 100},
		{"humidity", 30, 80},
		{"pressure", 900, 1100},
	}

	for _, t := range defaults {
		_, err := db.Exec(`
            INSERT INTO thresholds (type, min_value, max_value)
            VALUES ($1, $2, $3)
            ON CONFLICT (type) DO UPDATE SET 
                min_value = EXCLUDED.min_value, 
                max_value = EXCLUDED.max_value,
                updated_at = CURRENT_TIMESTAMP`,
			t.Type, t.MinValue, t.MaxValue)
		if err != nil {
			return err
		}
	}
	return nil
}
