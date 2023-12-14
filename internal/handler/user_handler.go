package handler

import (
	"encoding/json"
	"github.com/greeflas/itea_go_backend/internal/repository"
	"github.com/greeflas/itea_go_backend/internal/service"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/greeflas/itea_go_backend/pkg/server"
)

type UserListItem struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type UserHandler struct {
	logger         *log.Logger
	userRepository *repository.UserInMemoryRepository
	userService    *service.UserService
}

func NewUserHandler(
	logger *log.Logger,
	userRepository *repository.UserInMemoryRepository,
	userService *service.UserService,
) *UserHandler {
	return &UserHandler{
		logger:         logger,
		userRepository: userRepository,
		userService:    userService,
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := h.handleGet(w, r); err != nil {
			server.SendInternalServerError(w)
		}
	case http.MethodPost:
		if err := h.handlePost(w, r); err != nil {
			server.SendInternalServerError(w)
		}
	case http.MethodPatch:
		if err := h.handlePatch(w, r); err != nil {
			server.SendInternalServerError(w)
		}
	case http.MethodDelete:
		if err := h.handleDelete(w, r); err != nil {
			server.SendInternalServerError(w)
		}
	default:
		server.SendError(w, "invalid HTTP method", http.StatusBadRequest)
	}
}

func (h *UserHandler) handleGet(w http.ResponseWriter, r *http.Request) error {
	users := h.userRepository.GetAll()

	listItems := make([]*UserListItem, len(users))
	for _, user := range users {
		listItem := &UserListItem{
			Id:    user.Id.String(),
			Email: user.Email,
		}
		listItems = append(listItems, listItem)
	}

	return json.NewEncoder(w).Encode(listItems)
}

func (h *UserHandler) handlePost(w http.ResponseWriter, r *http.Request) error {
	args := new(service.NewUserArgs)

	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		return err
	}
	defer r.Body.Close()

	// TODO: validate args
	if err := h.userService.Create(args); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *UserHandler) handlePatch(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")

	userId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	args := new(service.UpdatedUserArgs)

	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		return err
	}
	defer r.Body.Close()

	// TODO: validate args
	if err := h.userService.Update(userId, args); err != nil {
		return err
	}

	return nil
}

func (h *UserHandler) handleDelete(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")

	userId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if err := h.userService.Delete(userId); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
