package api

import (
	"encoding/json"
	"net/http"
	"realtime-app/models"

	"github.com/jmoiron/sqlx"
)

func GetCurrentParameters(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params []models.CurrentParameter
		query := `SELECT * FROM current_parameters 
                 ORDER BY timestamp DESC LIMIT 100`
		err := db.Select(&params, query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, params)
	}
}

func UpdateReferenceParameter(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")

		var refParam models.ReferenceParameter
		if err := json.NewDecoder(r.Body).Decode(&refParam); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(`INSERT INTO reference_parameters 
                          (id_param, min_value, max_value)
                          VALUES ($1, $2, $3)
                          ON CONFLICT (id_param) DO UPDATE SET
                          min_value = EXCLUDED.min_value,
                          max_value = EXCLUDED.max_value`,
			refParam.ParamID, refParam.Min, refParam.Max)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
