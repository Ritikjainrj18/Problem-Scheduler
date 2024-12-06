package main

import (
	"database/sql"
	"fmt"
	"log"
	"ritikjainrj18/backend/config"
	"ritikjainrj18/backend/db"
	"ritikjainrj18/backend/service/task"
	"ritikjainrj18/backend/types"

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

	for {
		tx, err := db.Begin()
		if err != nil {
			log.Println("Failed to begin transaction:", err)
			continue
		}

		rows, err := tx.Query("SELECT * FROM tasks WHERE scheduledAt < DATE_ADD(NOW(), INTERVAL 30 SECOND) AND pickedAt IS NULL ORDER BY scheduledAt LIMIT 10 FOR UPDATE SKIP LOCKED;")
		if err != nil {
			log.Println("Failed to query tasks: ", err)
			tx.Rollback()
			continue
		}
		tasks := make([]types.Task, 0)
		for rows.Next() {
			t, err := task.ScanRowIntoTask(rows)
			if err != nil {
				log.Println("Failed to scan tasks: ", err)
				tx.Rollback()
				continue
			}
			tasks = append(tasks, *t)
		}

		if err := rows.Err(); err != nil {
			log.Println("Error iterating rows: ", err)
			continue
		}
		rows.Close()

		for _, tasks := range tasks {
			fmt.Println("Processing task:", tasks)
			// push to kafka
			_, err = tx.Exec("UPDATE tasks SET pickedAt = NOW() WHERE id = ?", tasks.ID)
			if err != nil {
				log.Println("Failed to update task status: ", err)
				tx.Rollback()
				continue
			}
		}

		err = tx.Commit()
		if err != nil {
			log.Println("Failed to commit transaction: ", err)
			continue
		}
	}

}
func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Producer: Successfully connected to DB!")
}

// PICKER and executer are differnet
// 1) if million mails to be send i need million concurrent connection to db thats not possible can be done in kafaka
//   so better to batch read and gave multiple connection on kafaka
