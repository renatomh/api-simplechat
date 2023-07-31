package main

import (
	"database/sql"
	"log"

	"github.com/renatomh/api-simplechat/api"
	db "github.com/renatomh/api-simplechat/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_chat?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	// Connecting to the database
	conn, err := sql.Open(dbDriver, dbSource)
	// If an error is returned
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
