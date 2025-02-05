package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Ritikjainrj18/Problem-Scheduler/Backend/config"
	"github.com/Ritikjainrj18/Problem-Scheduler/Backend/db"
	"github.com/go-sql-driver/mysql"
)

type Response struct {
	Result struct {
		Problems []Problem `json:"problems"`
	} `json:"result"`
}

type Problem struct {
	ContestID  int      `json:"contestId"`
	QuestionId string   `json:"index"`
	Points     *float64 `json:"rating,omitempty"` // Use *int to handle missing values
}

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 "127.0.0.1",
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

	url := "https://codeforces.com/api/problemset.problems"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	_, err = db.Exec("DELETE FROM problems")
	if err != nil {
		fmt.Printf("Unable to remove problems")
	}

	// Iterate and extract required fields
	for _, problem := range data.Result.Problems {
		problemId := fmt.Sprintf("%d/%s", problem.ContestID, problem.QuestionId)
		problemPoints := 0
		if problem.Points != nil {
			problemPoints = int(*problem.Points)
		}
		_, err := db.Exec("INSERT INTO problems (points,uniqueCode) VALUES(?,?)", problemPoints, problemId)
		if err != nil {
			fmt.Printf("Unable to add problem %s\n", problemId)
		}
	}

}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Problem server: Successfully connected to DB!")
}
