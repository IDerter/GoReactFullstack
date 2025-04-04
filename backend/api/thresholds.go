package api

import (
	"encoding/json"
	"log"
	"net/http"
	"realtime-app/models"

	"github.com/jmoiron/sqlx"
)

func GetThresholds(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		log.Println("Received thresholds request") // Добавьте это
		var thresholds []models.Threshold
		err := db.Select(&thresholds, "SELECT * FROM thresholds ORDER BY type")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(thresholds)
	}
}

type UpdateThresholdCallback func(threshold models.Threshold)

func UpdateThresholdWrapper(db *sqlx.DB, callback UpdateThresholdCallback) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var threshold models.Threshold
		if err := json.NewDecoder(r.Body).Decode(&threshold); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Обновление в БД
		if _, err := db.Exec(`
            INSERT INTO thresholds (type, min_value, max_value)
            VALUES ($1, $2, $3)
            ON CONFLICT (type) DO UPDATE SET 
                min_value = EXCLUDED.min_value, 
                max_value = EXCLUDED.max_value,
                updated_at = CURRENT_TIMESTAMP`,
			threshold.Type, threshold.MinValue, threshold.MaxValue); err != nil {
			http.Error(w, "Failed to update threshold", http.StatusInternalServerError)
			return
		}

		// Вызов callback для обновления в памяти
		callback(threshold)

		// Явно устанавливаем Content-Type перед отправкой ответа
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}
