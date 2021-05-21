package httpapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mymmrac/project-glynn/pkg/data/chat"
	"github.com/mymmrac/project-glynn/pkg/server"
	"github.com/mymmrac/project-glynn/pkg/uuid"
	"github.com/sirupsen/logrus"
)

const (
	roomIDParameter        = "roomID"
	LastMessageIDParameter = "lastMessageID"
)

// Server http api
type Server struct {
	service *server.Service
	router  mux.Router
	log     *logrus.Logger
}

// NewServer creates new server and initializes routes
func NewServer(service *server.Service, log *logrus.Logger) *Server {
	srv := &Server{
		service: service,
		log:     log,
	}
	srv.routes()
	return srv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	handlers.CORS(origins, methods)(&s.router).ServeHTTP(w, r)
}

func (s *Server) routes() {
	api := s.router.PathPrefix("/api").Subrouter()
	roomMessagesAPI := api.PathPrefix(fmt.Sprintf("/rooms/{%s:%s}/messages", roomIDParameter, uuid.Regex)).Subrouter()

	roomMessagesAPI.HandleFunc("", s.getMessages()).
		Methods(http.MethodGet)
	roomMessagesAPI.HandleFunc("", s.getMessages()).
		Queries(LastMessageIDParameter, fmt.Sprintf("{%s:%s}", LastMessageIDParameter, uuid.Regex)).
		Methods(http.MethodGet)
	roomMessagesAPI.HandleFunc("", s.sendMassage()).
		Methods(http.MethodPost)
}

func (s *Server) getMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// TODO move to func
		roomIDStr, ok := vars[roomIDParameter]
		if !ok {
			s.log.Error("roomID is required")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		roomID, err := uuid.Parse(roomIDStr)
		if err != nil {
			s.log.Error("invalid roomID: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var messages *chat.Messages

		lastMessageIDStr := r.URL.Query().Get(LastMessageIDParameter)
		if lastMessageIDStr == "" {
			messages, err = s.service.GetMessagesLatest(roomID)
		} else {
			var lastMessageID uuid.UUID
			if lastMessageID, err = uuid.Parse(lastMessageIDStr); err != nil {
				if err = respondJSONError(w, err, http.StatusBadRequest); err != nil {
					s.log.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			messages, err = s.service.GetMessagesAfterMessage(roomID, lastMessageID)
		}

		if err != nil {
			status := http.StatusBadRequest
			if errors.Is(err, server.ErrorRoomNotFound) {
				status = http.StatusNotFound
			}

			err := respondJSONError(w, err, status)
			if err != nil {
				s.log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		err = respondJSON(w, messages, http.StatusOK)
		if err != nil {
			s.log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (s *Server) sendMassage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// TODO move to func
		roomIDStr, ok := vars[roomIDParameter]
		if !ok {
			s.log.Error("roomID is required")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		roomID, err := uuid.Parse(roomIDStr)
		if err != nil {
			s.log.Error("invalid roomID: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var newMessage chat.NewMessage
		err = decodeJSON(r, &newMessage)
		if err != nil {
			err := respondJSONError(w, err, http.StatusBadRequest)
			if err != nil {
				s.log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		err = s.service.SendMessage(roomID, newMessage)
		if err != nil {
			err := respondJSONError(w, err, http.StatusBadRequest)
			if err != nil {
				s.log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
