package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/customerrors"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type UsersUsecases interface {
	Create(ctx context.Context, u entity.User) (string, error)
	Update(ctx context.Context, u entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Delete(ctx context.Context, id string) error
}

const idKey = "id"

type Handlers struct {
	usecases UsersUsecases
	logger   *zap.Logger
}

func NewHandlers(uc UsersUsecases, log *zap.Logger) *Handlers {
	return &Handlers{usecases: uc, logger: log}
}

// GetUsers Get all users.
func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.usecases.GetAll(r.Context())
	if err != nil {
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.logger.Info("ger all request succeeded")
	h.render(w, users)
}

// GetUser Get single user
func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // Gets params
	// Loop through user and find one with the id from the params
	id := params[idKey]

	user, err := h.usecases.GetByID(r.Context(), id)
	if err == nil {
		h.render(w, user)
		return
	}

	h.logger.Error("can't get a user", zap.Error(err))

	// if you wan't you can set content type of the headers directly here
	w.Header().Set("Content-Type", "application/json")
	h.writeErrorResponse(http.StatusNotFound,
		fmt.Sprintf("can't find a user with %v id", id), w)
}

// CreateUser Add new user
func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}

	id, err := h.usecases.Create(r.Context(), user)
	if err != nil {
		h.logger.Error("can't create a user", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	user.ID = id
	h.logger.Info("user has been created", zap.String(idKey, id))
	h.render(w, user)
}

// UpdateUser updates a user
func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}

	user.ID = params["id"]
	if err := h.usecases.Update(r.Context(), user); err != nil {
		h.logger.Error("can't update a user", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
			h.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a user with %v id", user.ID), w)
			return
		}

		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.render(w, user)
}

// DeleteUser deletes a user from storage
func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params[idKey]
	if err := h.usecases.Delete(r.Context(), id); err != nil {
		h.logger.Error("can't delete", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
			h.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a user with %v id", id), w)
			return
		}

		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.render(w, Message{"deleted"})
}

type Message struct {
	Msg string
}

func (*Handlers) writeErrorResponse(code int, msg string, w http.ResponseWriter) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Message{msg})
}

func (h *Handlers) render(w http.ResponseWriter, data interface{}) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(data)
}
