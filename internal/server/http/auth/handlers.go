package password

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/entity"
	authPkg "github.com/ITA-Dnipro/Dp-210_Go/internal/service/auth"
	"go.uber.org/zap"
)

type Auth interface {
	CreateToken(user authPkg.UserAuth) (authPkg.JwtToken, error)
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

	tk, err := h.auth.CreateToken(authPkg.UserAuth{Id: user.ID, Role: user.PermissionRole})
	if err != nil {
		h.writeErrorResponse(http.StatusInternalServerError, "could not generate token", w)
	}

	h.paswCases.DeleteCode(r.Context(), req.Email)

	h.render(w, tk)
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
		Token authPkg.JwtToken `json:"token"`
	}
	tkn.Token, err = h.auth.CreateToken(authPkg.UserAuth{Id: user.ID, Role: user.PermissionRole})
	if err != nil {
		h.writeErrorResponse(http.StatusUnauthorized, err.Error(), w)
		return
	}
	h.logger.Info("ger all request succeeded")
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
		h.writeErrorResponse(http.StatusInternalServerError, "could not log out", w)
		return
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
