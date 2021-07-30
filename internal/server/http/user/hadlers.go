package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/customerrors"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/user/entity"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type UsersUsecases interface {
	Create(ctx context.Context, u entity.User) (string, error)
	Update(ctx context.Context, u entity.User) error
	GetByID(ctx context.Context, id string) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Delete(ctx context.Context, id string) error
}

const idKey = "id"

type Server struct {
	usecases UsersUsecases
	logger   *zap.Logger
}

func NewServer(uc UsersUsecases, log *zap.Logger) *Server {
	return &Server{usecases: uc, logger: log}
}

// GetUsers Get all users.
func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.usecases.GetAll(r.Context())
	if err != nil {
		s.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	s.logger.Info("ger all request succeeded")
	s.render(w, users)
}

// GetUser Get single user
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // Gets params
	// Loop through user and find one with the id from the params
	id := params[idKey]

	user, err := s.usecases.GetByID(r.Context(), id)
	if err == nil {
		s.render(w, user)
		return
	}

	s.logger.Error("can't get a user", zap.Error(err))

	// if you wan't you can set content type of the headers directly here
	w.Header().Set("Content-Type", "application/json")
	s.writeErrorResponse(http.StatusNotFound,
		fmt.Sprintf("can't find a user with %v id", id), w)
}

// CreateUser Add new user
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.writeErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}

	id, err := s.usecases.Create(r.Context(), user)
	if err != nil {
		s.logger.Error("can't create a user", zap.Error(err))
		s.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}
	user.ID = id
	s.logger.Info("user has been created", zap.String(idKey, id))
	s.render(w, user)
}

// UpdateUser updates a user
func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.writeErrorResponse(http.StatusBadRequest, "can't parse a user", w)
		return
	}

	user.ID = params["id"]
	if err := s.usecases.Update(r.Context(), user); err != nil {
		s.logger.Error("can't update a user", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
			s.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a user with %v id", user.ID), w)
			return
		}

		s.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	s.render(w, user)
}

// DeleteUser deletes a user from storage
func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params[idKey]
	if err := s.usecases.Delete(r.Context(), id); err != nil {
		s.logger.Error("can't delete", zap.Error(err))
		if errors.Is(err, customerrors.NotFound) {
			s.writeErrorResponse(http.StatusNotFound,
				fmt.Sprintf("can't find a user with %v id", id), w)
			return
		}

		s.writeErrorResponse(http.StatusInternalServerError, err.Error(), w)
		return
	}

	s.render(w, Message{"deleted"})
}

type Message struct {
	Msg string
}

func (*Server) writeErrorResponse(code int, msg string, w http.ResponseWriter) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Message{msg})
}

func (s *Server) render(w http.ResponseWriter, data interface{}) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(data)
}
