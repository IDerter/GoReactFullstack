package db

import (
	"fmt"
	"realtime-app/models"

	"github.com/jmoiron/sqlx"
)

func InitDB(db *sqlx.DB) error {
	// Создание таблиц
	if err := CreateTables(db); err != nil {
		return fmt.Errorf("ошибка создания таблиц: %v", err)
	}

	// Инициализация тестового оборудования
	if err := InitDefaultEquipment(db); err != nil {
		return fmt.Errorf("ошибка инициализации оборудования: %v", err)
	}

	// Инициализация пороговых значений
	if err := initDefaultThresholds(db); err != nil {
		return fmt.Errorf("ошибка инициализации порогов: %v", err)
	}

	return nil
}

func InitDefaultEquipment(db *sqlx.DB) error {
	defaultEquipment := []struct {
		Name   string
		Type   string
		Status string
	}{
		{"Пресс 1", "Пресс", "Рабочее"},
		{"Датчик температуры", "Датчик", "Рабочее"},
		{"Контроллер", "Контроллер", "В ремонте"},
	}

	for _, eq := range defaultEquipment {
		_, err := db.Exec(
			"INSERT INTO equipment (name, type, status) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
			eq.Name, eq.Type, eq.Status,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Инициализация пороговых значений
func initDefaultThresholds(db *sqlx.DB) error {
	defaults := []models.Threshold{
		{Type: "temperature", MinValue: 20, MaxValue: 100},
		{Type: "humidity", MinValue: 30, MaxValue: 80},
		{Type: "pressure", MinValue: 900, MaxValue: 1100},
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
			return fmt.Errorf("ошибка инициализации порога для %s: %v", t.Type, err)
		}
	}
	return nil
}
