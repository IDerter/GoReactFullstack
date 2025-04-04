package main

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// main.go
var (
	dbConnectionString = "user=postgres password=postgres host=postgres port=5433 dbname=realtime sslmode=disable connect_timeout=5"
)

func connectDBTest() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", dbConnectionString)
}

func TestConnectDB(t *testing.T) {
	originalConnStr := dbConnectionString
	defer func() { dbConnectionString = originalConnStr }()

	t.Run("Successful connection", func(t *testing.T) {
		dbConnectionString = "user=postgres password=postgres host=localhost port=5433 dbname=realtime sslmode=disable"

		db, err := connectDBTest()
		if err != nil {
			t.Fatalf("Failed to connect to test database: %v\n"+
				"Make sure PostgreSQL is running and accessible at localhost:5432\n"+
				"Test database can be created with: createdb -U postgres realtime", err)
		}
		defer db.Close()

		assert.NoError(t, err)
		assert.NotNil(t, db)

		err = db.Ping()
		assert.NoError(t, err)
	})

	t.Run("Invalid connection string", func(t *testing.T) {
		// Подменяем строку подключения на невалидную
		dbConnectionString = "invalid_connection_string"

		db, err := connectDBTest()
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}
