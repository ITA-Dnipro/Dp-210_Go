package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/auth"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/entity"
	md "github.com/ITA-Dnipro/Dp-210_Go/auth/internal/server/http/middleware"

	"go.uber.org/zap"
)

type Auth interface {
	CreateToken(user auth.UserAuth) (auth.JwtToken, error)
	InvalidateToken(userId string) error
}

type Handlers struct {
	paswCases PasswordUsecases
	logger    *zap.Logger
	auth      Auth
}

func NewHandler(paswCases PasswordUsecases, logger *zap.Logger, auth Auth) *Handlers {
	return &Handlers{
		paswCases: paswCases,
		logger:    logger,
		auth:      auth,
	}
}

type PasswordUsecases interface {
	SendRestorePasswordCode(ctx context.Context, email string) (code string, err error)
	Authenticate(ctx context.Context, pc entity.PasswordCode) (entity.User, error)
	Auth(ctx context.Context, email, password string) (u entity.User, err error)
	DeleteCode(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, passw entity.UserNewPassword) error
	SetNewPassword(ctx context.Context, password string, user *entity.User) error
}

func (h *Handlers) LogIn(w http.ResponseWriter, r *http.Request) {
	var newUser entity.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, invalidRequestFormat, w)
		return
	}

	user, err := h.paswCases.Auth(r.Context(), newUser.Email, newUser.Password)
	if err != nil {
		h.writeErrorResponse(http.StatusUnauthorized, incorrectEmailOrPassword, w)
		return
	}

	var tkn struct {
		Token auth.JwtToken `json:"token"`
	}
	tkn.Token, err = h.auth.CreateToken(auth.UserAuth{Id: user.ID, Role: user.PermissionRole})
	if err != nil {
		h.writeErrorResponse(http.StatusUnauthorized, incorrectEmailOrPassword, w)
		return
	}

	h.render(w, tkn)
}

func (h *Handlers) LogOut(w http.ResponseWriter, r *http.Request) {
	u, ok := md.UserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(http.StatusUnauthorized, "no such session", w)
		return
	}

	if err := h.auth.InvalidateToken(u.Id); err != nil {
		h.logger.Warn(fmt.Sprintf("log out: user %v; err: %v", u.Id, err))
		h.writeErrorResponse(http.StatusInternalServerError, requestFailed, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) SendRestorePasswordCode(w http.ResponseWriter, r *http.Request) {
	var req entity.PasswordRestoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, invalidRequestFormat, w)
		return
	}

	if _, err := h.paswCases.SendRestorePasswordCode(r.Context(), req.Email); err != nil {
		h.writeErrorResponse(http.StatusAccepted, requestFailed, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) RestorePassword(w http.ResponseWriter, r *http.Request) {
	var req entity.PasswordCode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, invalidRequestFormat, w)
		return
	}

	user, err := h.paswCases.Authenticate(r.Context(), req)
	if err != nil {
		h.writeErrorResponse(http.StatusForbidden, incorrectEmailOrAuthCode, w)
		return
	}

	defer h.paswCases.DeleteCode(r.Context(), req.Email)

	if err = h.paswCases.SetNewPassword(r.Context(), req.NewPassword, &user); err != nil {
		h.writeErrorResponse(http.StatusInternalServerError, requestFailed, w)
		return
	}

	tk, err := h.auth.CreateToken(auth.UserAuth{Id: user.ID, Role: user.PermissionRole})
	if err != nil {
		h.writeErrorResponse(http.StatusInternalServerError, requestFailed, w)
		return
	}

	h.render(w, tk)
}

func (h *Handlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req entity.UserNewPassword
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, invalidRequestFormat, w)
		return
	}

	u, ok := md.UserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(http.StatusUnauthorized, requestFailed, w)
		return
	}
	req.UserID = u.Id

	if req.Password == "" {
		h.writeErrorResponse(http.StatusBadRequest, "request does not meet needed criteria", w)
		return
	}

	if req.Password != req.PasswordConfirm {
		h.writeErrorResponse(http.StatusBadRequest, "new password and new password confirm confirm do not match", w)
		return
	}

	if err := h.paswCases.ChangePassword(r.Context(), req); err != nil {
		h.writeErrorResponse(http.StatusForbidden, "wrong password", w)
	}

	w.WriteHeader(http.StatusOK)
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
