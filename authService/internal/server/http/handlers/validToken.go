package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/authService/internal/usecase"
)

type UserToken struct {
	Token string `json:"token"`
}

func (h *Handlers) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var t UserToken
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		h.writeErrorResponse(http.StatusBadRequest, invalidRequestFormat, w)
		return
	}

	user, err := h.auth.ValidateToken(usecase.JwtToken(t.Token))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	h.render(w, user)
	w.WriteHeader(http.StatusOK)
}
