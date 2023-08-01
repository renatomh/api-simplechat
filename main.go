package main

import (
	"database/sql"
	"log"

	"github.com/renatomh/api-simplechat/api"
	db "github.com/renatomh/api-simplechat/db/sqlc"
	"github.com/renatomh/api-simplechat/util"
)

func main() {
	// Here, we provide the current folder (root dir) as the path where viper should look for config files
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Connecting to the database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	// If an error is returned
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
