package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"ritikjainrj18/backend/config"
	"ritikjainrj18/backend/db"
	"ritikjainrj18/backend/types"
	"syscall"

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

	topic := "problem-email"
	group := "email-consumer-group"

	consumerGroup, err := ConnectConsumerGroup([]string{"localhost:9092"}, group)
	if err != nil {
		panic(err)
	}
	defer consumerGroup.Close()

	log.Println("Consumer Started")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := &ConsumerGroupHandler{
		db:     db,
		topic:  topic,
		msgCnt: 0,
	}

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
				log.Fatalf("Error during consumption: %v", err)
			}
		}
	}()

	// Wait for SIGINT or SIGTERM to exit gracefully
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan
}

type ConsumerGroupHandler struct {
	db     *sql.DB
	topic  string
	msgCnt int
}

func ConnectConsumerGroup(brokers []string, group string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewConsumerGroup(brokers, group, config)
}

func (h *ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session setup")
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session cleanup")
	log.Println("Shutdown signal received")
	log.Println("Processed ", h.msgCnt, " messages")
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Received order | Topic(%s) | Message(%s) \n", string(msg.Topic), string(msg.Value))

		var taskMsg types.Task
		if err := json.Unmarshal(msg.Value, &taskMsg); err != nil {
			log.Printf("Failed to deserialize message: %v", err)
			continue
		}

		// Process the message and update the database
		_, err := h.db.Exec("UPDATE tasks SET executedAt = NOW() WHERE id = ?", taskMsg.ID)
		if err != nil {
			log.Println("executedAt not updated", err)
		} else {
			// Mark message as processed
			h.msgCnt++
			session.MarkMessage(msg, "")
			session.Commit()
		}
	}
	return nil
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Consumer: Successfully connected to DB!")
}
