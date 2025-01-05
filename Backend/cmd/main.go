package main

import (
	"database/sql"
	"log"

	"github.com/Ritikjainrj18/Problem-Scheduler/Backend/cmd/api"
	"github.com/Ritikjainrj18/Problem-Scheduler/Backend/config"
	"github.com/Ritikjainrj18/Problem-Scheduler/Backend/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)
	defer db.Close()

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("API server: Successfully connected to DB!")
}
