package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"time"
)

func GenerateConnectionMap(configPath string) (map[string]Connection, error) {
	buffer, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(buffer, &config); err != nil {
		return nil, err
	}

	connections := make(map[string]Connection)
	for _, conn := range config.Connections {
		if conn.Default {
			connections["default"] = conn
		}

		connections[conn.Name] = conn
	}

	return connections, nil
}

func OpenConnection(conn Connection) (*sql.DB, error) {
	db, err := sql.Open(conn.Driver, conn.DSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func ExecuteQuery(db *sql.DB, query string) ([]map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	results := make([]map[string]any, 0)

	i := 0
	for i < GetLimit() && rows.Next() {
		results = append(results, make(map[string]any))
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		for j, colName := range cols {
			val := columnPointers[j].(*any)
			results[i][colName] = *val
		}

		i++
	}

	return results, nil
}
