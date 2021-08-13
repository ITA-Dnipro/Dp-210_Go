package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/customerrors"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/handlers"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// UsersUsecases represent user userCases.
type UsersUsecases interface {
	Create(ctx context.Context, u entity.NewUser) (string, error)
	Update(ctx context.Context, u entity.NewUser) (entity.User, error)
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Delete(ctx context.Context, id string) error
}

const idKey = "id"

// Handlers represent a user handlers.
type Handlers struct {
	userCases UsersUsecases
	logger    *zap.Logger
}

// NewHandlers create new user handlers.
func NewHandlers(uc UsersUsecases, log *zap.Logger) *Handlers {
	return &Handlers{userCases: uc, logger: log}
}

// GetUsers Get all users.
func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userCases.GetAll(r.Context())
	if err != nil {
		handlers.WriteErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.logger.Info("ger all request succeeded")
	handlers.Render(w, users)
}

// GetUser Get single user by id.
func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
	user, err := h.userCases.GetByID(r.Context(), id)
	if err == nil {
		handlers.Render(w, user)
		return
	}

	h.logger.Error("can't get a user", zap.Error(err))

	// if you wan't you can set content type of the headers directly here
	w.Header().Set("Content-Type", "application/json")
	handlers.WriteErrorResponse(http.StatusNotFound,
		fmt.Sprintf("can't find a user with %v id", id), w)
}

// CreateUser Add new user
func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser entity.NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		handlers.WriteErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}
	if ok := handlers.IsRequestValid(&newUser); !ok {
		handlers.WriteErrorResponse(http.StatusBadRequest, "user data invalid", w)
		return
	}
	id, err := h.userCases.Create(r.Context(), newUser)
	if err != nil {
		h.logger.Error("can't create a user", zap.Error(err))
		handlers.WriteErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	newUser.ID = id
	h.logger.Info("user has been created", zap.String(idKey, id))
	handlers.Render(w, newUser)
}

// UpdateUser updates a user
func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params

	var newUser entity.NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		handlers.WriteErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}

	if ok := handlers.IsRequestValid(&newUser); !ok {
		handlers.WriteErrorResponse(http.StatusBadRequest, "user data invalid", w)
		return
	}
	newUser.ID = id
	user, err := h.userCases.Update(r.Context(), newUser)
	if err != nil {
		h.logger.Error("can't update a user", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
			handlers.WriteErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a user with %v id", newUser.ID), w)
			return
		}

		handlers.WriteErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	handlers.Render(w, user)
}

// DeleteUser deletes a user from storage
func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
	if err := h.userCases.Delete(r.Context(), id); err != nil {
		h.logger.Error("can't delete", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
			handlers.WriteErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a user with %v id", id), w)
			return
		}

		handlers.WriteErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	handlers.Render(w, handlers.Message{Msg: "deleted"})
}
