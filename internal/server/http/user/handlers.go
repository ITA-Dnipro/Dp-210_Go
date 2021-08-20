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

// UsersUsecases represent user userCases.
type UsersUsecases interface {
	Create(ctx context.Context, u entity.NewUser) (string, error)
	Update(ctx context.Context, u entity.NewUser) (entity.User, error)
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Delete(ctx context.Context, id string) error
	Authenticate(ctx context.Context, email, password string) (u entity.User, err error)
}

const idKey = "id"
const tokenTime = time.Minute * 15

// Handlers represent a user handlers.
type Handlers struct {
	userCases UsersUsecases
	logger    *zap.Logger
}

// NewHandlers create new user handlers.
func NewHandlers(uc UsersUsecases, log *zap.Logger) *Handlers {
	return &Handlers{userCases: uc, logger: log}
}

// GetToken by basic auth.
func (h *Handlers) GetToken(w http.ResponseWriter, r *http.Request) {
	var newUser entity.NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}
	if ok := isRequestValid(&newUser); !ok {
		h.writeErrorResponse(http.StatusBadRequest, "user data invalid", w)
		return
	}
	user, err := h.userCases.Authenticate(r.Context(), newUser.Email, newUser.Password)
	if err != nil {
		h.writeErrorResponse(http.StatusUnauthorized, err.Error(), w)
		return
	}
	var tkn struct {
		Token auth.JwtToken `json:"token"`
	}
	tkn.Token, err = auth.CreateToken(auth.UserAuth{Id: user.ID, Role: user.PermissionRole}, tokenTime)
	if err != nil {
		h.writeErrorResponse(http.StatusUnauthorized, err.Error(), w)
		return
	}
	h.logger.Info("ger all request succeeded")
	h.render(w, tkn)
}

// GetUsers Get all users.
func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userCases.GetAll(r.Context())
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
	user, err := h.userCases.GetByID(r.Context(), id)
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
	var newUser entity.NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}
	if ok := isRequestValid(&newUser); !ok {
		h.writeErrorResponse(http.StatusBadRequest, "user data invalid", w)
		return
	}
	id, err := h.userCases.Create(r.Context(), newUser)
	if err != nil {
		h.logger.Error("can't create a user", zap.Error(err))
		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	newUser.ID = id
	h.logger.Info("user has been created", zap.String(idKey, id))
	h.render(w, newUser)
}

// UpdateUser updates a user
func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params

	var newUser entity.NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}

	if ok := isRequestValid(&newUser); !ok {
		h.writeErrorResponse(http.StatusBadRequest, "user data invalid", w)
		return
	}
	newUser.ID = id
	user, err := h.userCases.Update(r.Context(), newUser)
	if err != nil {
		h.logger.Error("can't update a user", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
			h.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a user with %v id", newUser.ID), w)
			return
		}

		h.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	h.render(w, user)
}

// DeleteUser deletes a user from storage
func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, idKey) // Gets params
	if err := h.userCases.Delete(r.Context(), id); err != nil {
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

func isRequestValid(nu interface{}) bool {
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
