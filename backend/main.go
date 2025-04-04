package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"realtime-app/api"
	"realtime-app/db"
	"realtime-app/models"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const timeRefresh time.Duration = 1 * time.Second

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Глобальная переменная для хранения текущих порогов
var currentThresholds = map[string]models.Threshold{
	"temperature": {Type: "temperature", MinValue: 20, MaxValue: 35}, // Дефолтные значения
	"humidity":    {Type: "humidity", MinValue: 30, MaxValue: 80},
	"pressure":    {Type: "pressure", MinValue: 900, MaxValue: 1100},
}

func main() {
	// Подключение к PostgreSQL
	dbConn, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Инициализация базы данных
	if err := db.InitDB(dbConn); err != nil {
		log.Fatal(err)
	}

	// Загрузка порогов из БД при старте
	if err := loadThresholds(dbConn); err != nil {
		log.Printf("Warning: couldn't load thresholds: %v", err)
	}

	// Настройка HTTP маршрутов
	setupRoutes(dbConn)

	// Запуск сервера
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Загрузка порогов из БД
func loadThresholds(db *sqlx.DB) error {
	var thresholds []models.Threshold
	if err := db.Select(&thresholds, "SELECT * FROM thresholds"); err != nil {
		return err
	}

	for _, t := range thresholds {
		currentThresholds[t.Type] = t
	}
	return nil
}

// Функция подключения к базе данных
func connectDB() (*sqlx.DB, error) {
	connStr := "user=postgres password=postgres host=postgres port=5432 dbname=realtime sslmode=disable connect_timeout=5"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %v\nСтрока подключения: %s", err, connStr)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// Настройка маршрутов HTTP
func setupRoutes(db *sqlx.DB) {
	// WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, db)
	})

	// API endpoints
	http.HandleFunc("/api/thresholds", api.GetThresholds(db))
	http.HandleFunc("/api/thresholds/update", api.UpdateThresholdWrapper(db, updateThresholdCallback))

}

func updateThresholdCallback(updatedThreshold models.Threshold) {
	currentThresholds[updatedThreshold.Type] = updatedThreshold
	log.Printf("Thresholds updated: %+v", updatedThreshold)
}

// WebSocket handler
func wsHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(timeRefresh)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Генерация данных с учетом текущих порогов
			sensorData, err := generateSensorData(db)
			if err != nil {
				log.Printf("Error generating sensor data: %v", err)
				continue
			}

			// Получение текущих порогов
			var thresholds []models.Threshold
			for _, t := range currentThresholds {
				thresholds = append(thresholds, t)
			}

			// Отправка данных клиенту
			if err := conn.WriteJSON(map[string]interface{}{
				"data":       sensorData,
				"thresholds": thresholds,
			}); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

// Генерация данных датчиков
func generateSensorData(db *sqlx.DB) ([]models.SensorData, error) {
	var allData []models.SensorData
	randSrc := rand.New(rand.NewSource(time.Now().UnixNano()))

	for sensorType, threshold := range currentThresholds {
		// Генерация значения в диапазоне [min, max]
		value := threshold.MinValue + randSrc.Float64()*(threshold.MaxValue-threshold.MinValue)

		// Сохранение в базу данных
		if _, err := db.Exec(
			"INSERT INTO sensor_data (value, type) VALUES ($1, $2)",
			value, sensorType,
		); err != nil {
			return nil, fmt.Errorf("DB insert error: %v", err)
		}

		// Получение последних 10 записей
		var records []models.SensorData
		if err := db.Select(
			&records,
			"SELECT * FROM sensor_data WHERE type=$1 ORDER BY timestamp DESC LIMIT 10",
			sensorType,
		); err != nil {
			return nil, fmt.Errorf("DB select error: %v", err)
		}

		allData = append(allData, records...)
	}

	return allData, nil
}
