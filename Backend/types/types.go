package types

import "time"

type UserStore interface {
	CreateUser(User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
}

type TaskStore interface {
	CreateTask(Task) error
	GetTaskByID(id int) (*Task, error)
	GetAllTasksByUserID(userId int) ([]Task, error)
}

type Task struct {
	ID            int       `json:"id"`
	UserID        int       `json:"userId"`
	Days          int       `json:"days"`
	MinimumRating int       `json:"minimumRating"`
	MaximumRating int       `json:"maximumRating"`
	Retries       int       `json:"retries" `
	CreatedAt     time.Time `json:"createdAt"`
	PickedAt      time.Time `json:"pickedAt" `
	ExecutedAt    time.Time `json:"executedAt" `
	//Name          string    `json:"name" validate:"required"`
	//Email         string    `json:"email" validate:"required"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateTaskPayload struct {
	UserID        int `json:"userId" validate:"required"`
	Days          int `json:"days" validate:"required"`
	MinimumRating int `json:"minimumRating" validate:"required"`
	MaximumRating int `json:"maximumRating" validate:"required"`
	Retries       int `json:"retries" `
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" `
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
