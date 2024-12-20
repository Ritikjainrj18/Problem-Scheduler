package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"ritikjainrj18/backend/config"
	"ritikjainrj18/backend/db"
	"ritikjainrj18/backend/service/task"
	"ritikjainrj18/backend/types"

	"github.com/IBM/sarama"
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
			log.Println("Processing task:", tasks)

			taskBytes, err := json.Marshal(tasks)
			if err != nil {
				log.Println("Failed to serialize task:", err)
				continue
			}

			err = PushOrdersToQueue("problem-email", taskBytes)
			if err != nil {
				log.Println("Unable to push to broker", err)
				continue
			}
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

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	return sarama.NewSyncProducer(brokers, config)
}

func PushOrdersToQueue(topic string, message []byte) error {
	brokers := []string{"localhost:9092"}

	producer, err := ConnectProducer(brokers)

	if err != nil {
		return err
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Order is stored in topic(%s)/partition(%d)/offset(%d)\n",
		topic,
		partition,
		offset)

	return nil
}

// PICKER and executer are differnet
// 1) if million mails to be send i need million concurrent connection to db thats not possible can be done in kafaka
//   so better to batch read and gave multiple connection on kafaka
