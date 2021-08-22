package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/customerrors"
	auth "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware/auth"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// NewUser represent new user in request.
type RequestUser struct {
	Name            string `json:"name,omitempty" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"omitempty,eqfield=Password"`
}

// UsersUsecases represent user usecases.
type UsersUsecases interface {
	Create(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Delete(ctx context.Context, id string) error
	Authenticate(ctx context.Context, email, password string) (entity.User, error)
}

const idKey = "id"
const tokenTime = time.Minute * 15

// Handlers represent a user handlers.
type Handlers struct {
	uc     UsersUsecases
	logger *zap.Logger
}

// NewHandlers create new user handlers.
func NewHandlers(uc UsersUsecases, log *zap.Logger) *Handlers {
	return &Handlers{uc: uc, logger: log}
}

// GetToken by basic auth.
func (h *Handlers) GetToken(w http.ResponseWriter, r *http.Request) {
	var ru RequestUser
	if err := json.NewDecoder(r.Body).Decode(&ru); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}
	if ok := isRequestValid(&ru); !ok {
		h.writeErrorResponse(http.StatusBadRequest, "user data invalid", w)
		return
	}
	u, err := h.uc.Authenticate(r.Context(), ru.Email, ru.Password)
	if err != nil {
		h.writeErrorResponse(http.StatusUnauthorized, err.Error(), w)
		return
	}
	var tkn struct {
		Token auth.JwtToken `json:"token"`
	}
	tkn.Token, err = auth.CreateToken(u.ID, u.PermissionRole, tokenTime)
	if err != nil {
		h.writeErrorResponse(http.StatusUnauthorized, err.Error(), w)
		return
	}
	h.logger.Info("ger all request succeeded")
	h.render(w, tkn)
}

// GetUsers Get all users.
func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.GetAll(r.Context())
	if err != nil {
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.logger.Info("ger all request succeeded")
	h.render(w, users)
}

// GetUser Get single user by id.
func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
	user, err := h.uc.GetByID(r.Context(), id)
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
	var ru RequestUser
	if err := json.NewDecoder(r.Body).Decode(&ru); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}
	if ok := isRequestValid(&ru); !ok {
		h.writeErrorResponse(http.StatusBadRequest, "user data invalid", w)
		return
	}
	u := entity.User{}
	u.Email = ru.Email
	u.Name = ru.Name
	u.PasswordHash = ru.Password

	if err := h.uc.Create(r.Context(), &u); err != nil {
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

	if ok := isUserRequestValid(&u); !ok {
		h.writeErrorResponse(http.StatusBadRequest, "user data invalid", w)
		return
	}
	u.ID = id
	if err := h.uc.Update(r.Context(), &u); err != nil {
		h.logger.Error("can't update a user", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
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
	if err := h.uc.Delete(r.Context(), id); err != nil {
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

func isRequestValid(nu *RequestUser) bool {
	validate := validator.New()
	err := validate.Struct(nu)
	fmt.Println(err)
	return err == nil
}

func isUserRequestValid(nu *entity.User) bool {
	validate := validator.New()
	err := validate.Struct(nu)
	fmt.Println(err)
	return err == nil
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
