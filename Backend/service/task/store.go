package task

import (
	"database/sql"
	"fmt"

	"github.com/Ritikjainrj18/Problem-Scheduler/Backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateTask(task types.Task) error {
	_, err := s.db.Exec("INSERT INTO tasks (userId, days, minimumRating, maximumRating, retries) VALUES (?, ?, ?, ?, ?)", task.UserID, task.Days, task.MinimumRating, task.MaximumRating, task.Retries)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAllTasksByUserID(userId int) ([]types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE userId = ?", userId)
	if err != nil {
		return nil, err
	}

	tasks := make([]types.Task, 0)
	defer rows.Close()
	for rows.Next() {
		t, err := ScanRowIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, nil
}

func (s *Store) GetTaskByID(id int) (*types.Task, error) {
	rows, err := s.db.Query("SELECT * FROM tasks WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	task := new(types.Task)
	defer rows.Close()
	for rows.Next() {
		task, err = ScanRowIntoTask(rows)
		if err != nil {
			return nil, err
		}
	}
	if task.ID == 0 {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}
func ScanRowIntoTask(rows *sql.Rows) (*types.Task, error) {
	task := new(types.Task)
	err := rows.Scan(
		&task.ID,
		&task.UserID,
		&task.Days,
		&task.MinimumRating,
		&task.MaximumRating,
		&task.Retries,
		&task.ScheduledAt,
		&task.PickedAt,
		&task.ExecutedAt,
	)

	if err != nil {
		return nil, err
	}

	return task, nil
}
