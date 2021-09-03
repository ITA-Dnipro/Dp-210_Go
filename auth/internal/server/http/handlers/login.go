package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/auth"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/entity"
	md "github.com/ITA-Dnipro/Dp-210_Go/auth/internal/server/http/middleware"
	"net/http"
)

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
