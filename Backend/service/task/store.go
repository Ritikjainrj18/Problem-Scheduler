package task

import (
	"database/sql"
	"fmt"
	"ritikjainrj18/backend/types"
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
	for rows.Next() {
		t, err := scanRowIntoTask(rows)
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
	for rows.Next() {
		task, err = scanRowIntoTask(rows)
		if err != nil {
			return nil, err
		}
	}
	if task.ID == 0 {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}
func scanRowIntoTask(rows *sql.Rows) (*types.Task, error) {
	task := new(types.Task)
	err := rows.Scan(
		&task.ID,
		&task.UserID,
		&task.Days,
		&task.MinimumRating,
		&task.MaximumRating,
		&task.Retries,
		&task.CreatedAt,
		&task.PickedAt,
		&task.ExecutedAt,
	)

	if err != nil {
		return nil, err
	}

	return task, nil
}
