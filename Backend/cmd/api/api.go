package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Ritikjainrj18/Problem-Scheduler/Backend/service/task"
	"github.com/Ritikjainrj18/Problem-Scheduler/Backend/service/user"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	taskStore := task.NewStore(s.db)
	taskHandler := task.NewHandler(taskStore, userStore)
	taskHandler.RegisterRoutes(subrouter)

	log.Println("Listining on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
