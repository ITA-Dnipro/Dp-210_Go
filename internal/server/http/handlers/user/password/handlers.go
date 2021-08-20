package password

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware/auth"
	"go.uber.org/zap"
)

const tokLifetime = 10 * time.Minute

type Handlers struct {
	paswCases PasswordUsecases
	logger    *zap.Logger
}

func NewHandler(paswCases PasswordUsecases, logger *zap.Logger) *Handlers {
	return &Handlers{
		paswCases: paswCases,
		logger:    logger,
	}
}

type PasswordUsecases interface {
	SendRestorePasswordCode(ctx context.Context, email string) (code string, err error)
	Authenticate(ctx context.Context, pc entity.PasswordCode) (entity.User, error)
	DeleteCode(ctx context.Context, email string) error
	ChangePassword(ctx context.Context, passw entity.UserNewPassword) error
}

func (h *Handlers) SendRestorePasswordCode(w http.ResponseWriter, r *http.Request) {
	var req entity.PasswordRestoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "is not an email", w)
		return
	}

	if _, err := h.paswCases.SendRestorePasswordCode(r.Context(), req.Email); err != nil {
		h.writeErrorResponse(http.StatusAccepted, "your request failed", w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) CheckPasswordCode(w http.ResponseWriter, r *http.Request) {
	var req entity.PasswordCode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "not a restore password data format", w)
		return
	}

	user, err := h.paswCases.Authenticate(r.Context(), req)
	if err != nil {
		h.writeErrorResponse(http.StatusForbidden, "authorization code is wrong", w)
		return
	}

	tk, err := auth.CreateToken(auth.UserAuth{Id: user.ID, Role: user.PermissionRole}, tokLifetime)
	if err != nil {
		h.writeErrorResponse(http.StatusInternalServerError, "could not generate token", w)
	}

	h.paswCases.DeleteCode(r.Context(), req.Email)

	h.render(w, tk)
}

func (h *Handlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req entity.UserNewPassword
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, "could not parse request", w)
		return
	}

	u, ok := middleware.UserFromContext(r.Context())
	if !ok {
		h.writeErrorResponse(http.StatusUnauthorized, "auth error", w)
		return
	}
	req.UserID = u.Id

	if req.Password == "" {
		h.writeErrorResponse(http.StatusBadRequest, "request does not meet needed criterium", w)
		return
	}

	if req.Password != req.PasswordConfirm {
		h.writeErrorResponse(http.StatusBadRequest, "new password and new password confirm do not match", w)
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
