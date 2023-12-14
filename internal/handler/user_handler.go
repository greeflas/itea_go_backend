package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/greeflas/itea_go_backend/pkg/server"
)

type User struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type UserHandler struct {
	logger *log.Logger

	users []*User
}

func NewUserHandler(logger *log.Logger) *UserHandler {
	return &UserHandler{
		logger: logger,
		users:  make([]*User, 0),
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
	return json.NewEncoder(w).Encode(h.users)
}

func (h *UserHandler) handlePost(w http.ResponseWriter, r *http.Request) error {
	user := new(User)

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}
	defer r.Body.Close()

	h.users = append(h.users, user)

	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *UserHandler) handlePatch(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")

	userId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	var user *User

	for _, u := range h.users {
		if u.Id == userId {
			user = u
			break
		}
	}

	if user == nil {
		server.SendNotFound(w)
		return nil
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}

func (h *UserHandler) handleDelete(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("id")

	userId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	userIdx := -1

	for index, u := range h.users {
		if u.Id == userId {
			userIdx = index
			break
		}
	}

	if userIdx == -1 {
		server.SendNotFound(w)
		return nil
	}

	h.users = append(h.users[:userIdx], h.users[userIdx+1:]...)

	w.WriteHeader(http.StatusNoContent)

	return nil
}
