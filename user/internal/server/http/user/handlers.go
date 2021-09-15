package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/user/internal/server/http/customerrors"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// NewUser represent new user in request.
type NewUser struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"omitempty,eqfield=Password"`
}

// UsersUsecases represent user usecases.
type Usecase interface {
	Create(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, u *entity.User) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, ul *entity.UserList) error
	GetByID(ctx context.Context, id string) (entity.User, error)
}

const idKey = "id"

// Handlers represent a user handlers.
type Handlers struct {
	usecase Usecase
	logger  *zap.Logger
}

// NewHandlers create new user handlers.
func NewUserHandlers(uc Usecase, logger *zap.Logger) *chi.Mux {
	hs := &Handlers{usecase: uc, logger: logger}
	r := chi.NewRouter()
	r.Post("/", hs.CreateUser)       // POST /api/v1/users
	r.Get("/", hs.GetUsers)          // GET /api/v1/users
	r.Get("/{id}", hs.GetUser)       // GET /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	r.Put("/{id}", hs.UpdateUser)    // PUT /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	r.Delete("/{id}", hs.DeleteUser) // DELETE /api/v1/users/6ba7b810-9dad-11d1-80b4-00c04fd430c8
	return r
}

// GetUsers Get all users.
func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	var err error
	query := r.URL.Query()
	ul := entity.UserList{
		Cursor: query.Get("cursor"),
	}
	if query.Has("limit") {
		if ul.Limit, err = strconv.Atoi(query.Get("limit")); err != nil {
			h.writeErrorResponse(http.StatusBadRequest, err.Error(), w)
			return
		}
	}

	if err := h.usecase.GetAll(r.Context(), &ul); err != nil {
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.render(w, ul)
}

// GetUser Get single user by id.
func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
	user, err := h.usecase.GetByID(r.Context(), id)
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
	var nu NewUser
	if err := json.NewDecoder(r.Body).Decode(&nu); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}

	validate := validator.New()
	if err := validate.Struct(nu); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "user data invalid", w)
		return
	}

	u := entity.User{
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: []byte(nu.Password),
	}

	if err := h.usecase.Create(r.Context(), &u); err != nil {
		h.logger.Error("can't create a user", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	h.logger.Info("user has been created", zap.String(idKey, u.ID))
	h.render(w, u)
}

// UpdateUser updates a user
func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params

	var u entity.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}

	validate := validator.New()
	if err := validate.Struct(u); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "user data invalid", w)
		return
	}
	u.ID = id
	if err := h.usecase.Update(r.Context(), &u); err != nil {
		h.logger.Error("can't update a user", zap.Error(err))
		if errors.Is(err, customerrors.ErrNotFound) {
			h.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a user with %v id", u.ID), w)
			return
		}

		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.render(w, u)
}

// DeleteUser deletes a user from storage
func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
	if err := h.usecase.Delete(r.Context(), id); err != nil {
		h.logger.Error("can't delete", zap.Error(err))
		if errors.Is(err, customerrors.ErrNotFound) {
			h.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a user with %v id", id), w)
			return
		}

		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.render(w, Message{"deleted"})
}

// Message represent error message.
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
