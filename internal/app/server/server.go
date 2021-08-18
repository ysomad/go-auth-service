package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/ysomad/go-auth-service/internal/app/model"
	"github.com/ysomad/go-auth-service/internal/app/store"
)

const sessionName = "session_id"

var errIncorrectEmailOrPassword = errors.New("incorrect email or password")

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) error(w http.ResponseWriter, code int, err error) {
	s.respond(w, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionCreate()).Methods("POST")
}

func (s *server) getDecodedUserRequest(w http.ResponseWriter, r *http.Request) *userRequest {
	req := &userRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		s.error(w, http.StatusBadRequest, err)
	}

	return req
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := s.getDecodedUserRequest(w, r)

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, http.StatusCreated, u)
	}
}

func (s *server) handleSessionCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := s.getDecodedUserRequest(w, r)

		u, err := s.store.User().GetByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, http.StatusOK, nil)
	}
}
