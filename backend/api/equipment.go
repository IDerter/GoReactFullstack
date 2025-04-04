package api

import (
	"encoding/json"
	"net/http"
	"realtime-app/models"

	"github.com/jmoiron/sqlx"
)

func GetEquipmentList(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var equipment []models.Equipment
		err := db.Select(&equipment, "SELECT * FROM equipment")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, equipment)
	}
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
