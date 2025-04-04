package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"realtime-app/api"
	"realtime-app/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetThresholds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Настройка ожидаемого запроса и результата
	expectedThresholds := []models.Threshold{
		{Type: "temperature", MinValue: 20, MaxValue: 35},
		{Type: "humidity", MinValue: 30, MaxValue: 80},
	}

	rows := sqlmock.NewRows([]string{"type", "min_value", "max_value"}).
		AddRow("temperature", 20, 35).
		AddRow("humidity", 30, 80)

	mock.ExpectQuery("SELECT \\* FROM thresholds ORDER BY type").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/thresholds", nil)
	w := httptest.NewRecorder()

	handler := api.GetThresholds(sqlxDB)
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result []models.Threshold
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, expectedThresholds, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateThresholdWrapper(t *testing.T) {
	// Создаем mock для sql.DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer db.Close()

	// Оборачиваем в sqlx.DB
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Тестовые данные
	testThreshold := models.Threshold{
		Type:     "temperature",
		MinValue: 25,
		MaxValue: 40,
	}

	// Настройка ожидаемого запроса
	mock.ExpectExec("INSERT INTO thresholds (.+) VALUES (.+)").
		WithArgs(testThreshold.Type, testThreshold.MinValue, testThreshold.MaxValue).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Создание тестового запроса
	body, _ := json.Marshal(testThreshold)
	req := httptest.NewRequest("POST", "/api/thresholds/update", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Mock для callback
	called := false
	callback := func(threshold models.Threshold) {
		called = true
		assert.Equal(t, testThreshold, threshold)
	}

	// Вызов тестируемого метода
	handler := api.UpdateThresholdWrapper(sqlxDB, callback)
	handler(w, req)

	// Проверки
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var result map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "success", result["status"])
	assert.True(t, called, "Callback should be called")

	// Проверяем, что все ожидания по mock выполнены
	assert.NoError(t, mock.ExpectationsWereMet())
}
