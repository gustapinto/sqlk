package main

import (
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connections, err := GenerateConnectionMap(GetConfigFilePath())
	if err != nil || len(connections) == 0 {
		log.Fatalf("Failed to get connections, got error: %s", err.Error())
	}

	conn := connections[GetConnectionName()]
	db, err := OpenConnection(conn)
	if err != nil {
		log.Fatalf("Failed to connect to database, got error: %s", err.Error())
	}
	defer db.Close()

	query, err := GetQuery(GetQueryFilePath(conn))
	if err != nil {
		log.Fatalf("Failed to find query")
	}

	results, err := ExecuteQuery(db, query)
	if err != nil {
		log.Fatalf("Failed to execute query, got error: %s", err.Error())
	}

	var buffer []byte
	if len(results) == 1 {
		buffer, err = json.Marshal(results[0])
	} else {
		buffer, err = json.Marshal(results)
	}

	if err != nil {
		log.Fatalf("Failed to marshal result, got error: %s", err.Error())
	}

	fmt.Println(string(buffer))
}
