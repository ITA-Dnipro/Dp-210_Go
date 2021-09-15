package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/authService/internal/entity"
	md "github.com/ITA-Dnipro/Dp-210_Go/authService/internal/server/http/middleware"
	"github.com/ITA-Dnipro/Dp-210_Go/authService/internal/usecase"
)

func (h *Handlers) SendRestorePasswordCode(w http.ResponseWriter, r *http.Request) {
	var req entity.PasswordRestoreReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, invalidRequestFormat, w)
		return
	}

	if _, err := h.paswCases.SendRestorePasswordCode(r.Context(), req.Email); err != nil {
		h.logger.Warn(fmt.Sprintf("http: send restore code: %v", err))
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

	tk, err := h.auth.CreateToken(usecase.UserAuth{Id: user.ID, Role: user.PermissionRole})
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
