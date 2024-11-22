package task

import (
	"fmt"
	"net/http"
	"ritikjainrj18/backend/service/auth"
	"ritikjainrj18/backend/types"
	"ritikjainrj18/backend/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.TaskStore
	userStore types.UserStore
}

func NewHandler(store types.TaskStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/task", auth.WithJWTAuth(h.handleCreateTask, h.userStore)).Methods("POST")
	router.HandleFunc("/tasks", auth.WithJWTAuth(h.handleGetAllTasksByUserID, h.userStore)).Methods("GET")
	router.HandleFunc("/task/{taskID}", auth.WithJWTAuth(h.handleGetTaskByID, h.userStore)).Methods("GET")
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateTaskPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload &v", errors))
		return
	}

	err := h.store.CreateTask(types.Task{
		UserID:        auth.GetUserIDFromContext(r.Context()),
		Days:          payload.Days,
		MinimumRating: payload.MinimumRating,
		MaximumRating: payload.MaximumRating,
		Retries:       payload.Retries,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleGetAllTasksByUserID(w http.ResponseWriter, r *http.Request) {

	UserID := auth.GetUserIDFromContext(r.Context())

	tasks, err := h.store.GetAllTasksByUserID(UserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tasks)
}

func (h *Handler) handleGetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["taskID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing task ID"))
		return
	}

	taskID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	task, err := h.store.GetTaskByID(taskID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, task)
}
