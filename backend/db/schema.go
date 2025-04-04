package db

import "github.com/jmoiron/sqlx"

func CreateTables(db *sqlx.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS equipment (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        type VARCHAR(50) NOT NULL,
        status VARCHAR(20) CHECK (status IN ('Рабочее', 'Неисправное', 'В ремонте'))
    );
    
    CREATE TABLE IF NOT EXISTS process_parameters (
        id SERIAL PRIMARY KEY,
        id_equipment INT REFERENCES equipment(id),
        name VARCHAR(60) NOT NULL,
        units VARCHAR(60) NOT NULL
    );
    
    CREATE TABLE IF NOT EXISTS reference_parameters (
        id SERIAL PRIMARY KEY,
        id_param INT REFERENCES process_parameters(id),
        min_value FLOAT NOT NULL,
        max_value FLOAT NOT NULL
    );
    
    CREATE TABLE IF NOT EXISTS current_parameters (
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        id_param INT REFERENCES process_parameters(id),
        value FLOAT NOT NULL,
        PRIMARY KEY (timestamp, id_param)
    );
    
    CREATE INDEX IF NOT EXISTS idx_current_params_timestamp ON current_parameters (timestamp);
    CREATE INDEX IF NOT EXISTS idx_current_params_id ON current_parameters (id_param);

	CREATE TABLE IF NOT EXISTS sensor_data (
		id SERIAL PRIMARY KEY,
		value DOUBLE PRECISION NOT NULL,
		type TEXT NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS thresholds (
		id SERIAL PRIMARY KEY,
		type TEXT NOT NULL UNIQUE,
		min_value DOUBLE PRECISION NOT NULL,
		max_value DOUBLE PRECISION NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_sensor_data_type ON sensor_data (type);
	CREATE INDEX IF NOT EXISTS idx_sensor_data_timestamp ON sensor_data (timestamp);`

	_, err := db.Exec(schema)
	return err
}
