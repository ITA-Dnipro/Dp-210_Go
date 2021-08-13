package password

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/handlers"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware/auth"
	"go.uber.org/zap"
)

const tokenTime = time.Minute * 15

type PasswordUsecases interface {
	SendRestorePasswordCode(ctx context.Context, email string) (string, error)
	DeleteCode(ctx context.Context, email string) error
	Authenticate(ctx context.Context, pc entity.PasswordCode) (string, error)
}

type UsersUsecases interface {
	Authenticate(ctx context.Context, email, password string) (id string, err error)
	ChangePassword(ctx context.Context, passw entity.UserNewPassword) error
}

type Handlers struct {
	userCases UsersUsecases
	paswCases PasswordUsecases
	logger    *zap.Logger
}

func (h *Handlers) GetToken(w http.ResponseWriter, r *http.Request) {
	var newUser entity.NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		handlers.WriteErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}
	if ok := handlers.IsRequestValid(&newUser); !ok {
		handlers.WriteErrorResponse(http.StatusBadRequest, "user data invalid", w)
		return
	}
	id, err := h.userCases.Authenticate(r.Context(), newUser.Email, newUser.Password)
	if err != nil {
		handlers.WriteErrorResponse(http.StatusUnauthorized, err.Error(), w)
		return
	}
	var tkn struct {
		Token auth.JwtToken `json:"token"`
	}
	tkn.Token, err = auth.CreateToken(id, tokenTime)
	if err != nil {
		handlers.WriteErrorResponse(http.StatusUnauthorized, err.Error(), w)
		return
	}
	h.logger.Info("ger all request succeeded")
	handlers.Render(w, tkn)
}

func (h *Handlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req entity.UserNewPassword
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteErrorResponse(http.StatusBadRequest, "could not parse request", w)
		return
	}

	var ok bool
	if req.UserID, ok = middleware.FromContext(r.Context()); !ok {
		handlers.WriteErrorResponse(http.StatusUnauthorized, "auth error", w)
		return
	}

	if !handlers.IsRequestValid(&req) {
		handlers.WriteErrorResponse(http.StatusBadRequest, "request does not meet needed criterium", w)
		return
	}

	if req.Password != req.PasswordConfirm {
		handlers.WriteErrorResponse(http.StatusBadRequest, "new password and new password confirm do not match", w)
		return
	}

	if err := h.userCases.ChangePassword(r.Context(), req); err != nil {
		handlers.WriteErrorResponse(http.StatusForbidden, "wrong password", w)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) SendRestorePasswordCode(w http.ResponseWriter, r *http.Request) {
	var req entity.PasswordRestoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteErrorResponse(http.StatusBadRequest, "is not an email", w)
		return
	}

	if _, err := h.paswCases.SendRestorePasswordCode(r.Context(), req.Email); err != nil {
		handlers.WriteErrorResponse(http.StatusAccepted, "your request failed", w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) CheckPasswordCode(w http.ResponseWriter, r *http.Request) {
	var req entity.PasswordCode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteErrorResponse(http.StatusBadRequest, "not a restore password data format", w)
		return
	}

	uId, err := h.paswCases.Authenticate(r.Context(), req)
	if err != nil {
		handlers.WriteErrorResponse(http.StatusForbidden, "authorization code is wrong", w)
		return
	}

	tk, err := auth.CreateToken(uId, 10*time.Minute)
	if err != nil {
		handlers.WriteErrorResponse(http.StatusInternalServerError, "could not generate token", w)
	}

	h.paswCases.DeleteCode(r.Context(), req.Email)

	handlers.Render(w, tk)
}
