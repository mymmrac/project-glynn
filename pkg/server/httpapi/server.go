package httpapi

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mymmrac/project-glynn/pkg/data/message"
	"github.com/mymmrac/project-glynn/pkg/server"
	"github.com/sirupsen/logrus"
)

type Server struct {
	service *server.Service
	router  mux.Router
	log     *logrus.Logger
}

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

func (s *Server) getMessages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID := uuid.MustParse(vars["roomID"])

		var (
			messages  []message.Message
			usernames map[uuid.UUID]string
			err       error
		)

		lastMessageIDStr := r.URL.Query().Get("lastMessageID")
		ok, err := regexp.MatchString(uuidRegx, lastMessageIDStr)
		if err != nil {
			s.log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if ok {
			lastMessageID := uuid.MustParse(lastMessageIDStr)
			messages, usernames, err = s.service.GetMessagesAfterMessage(roomID, lastMessageID)
		} else {
			messages, usernames, err = s.service.GetMessagesLatest(roomID)
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

		err = respondJSON(w, struct {
			Messages  []message.Message    `json:"messages"`
			Usernames map[uuid.UUID]string `json:"usernames"`
		}{messages, usernames}, http.StatusOK)
		if err != nil {
			s.log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (s *Server) sendMassage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID := uuid.MustParse(vars["roomID"])

		messageData := struct {
			UserID uuid.UUID `json:"userID"`
			Text   string    `json:"text"`
		}{}
		err := decodeJSON(r, &messageData)
		if err != nil {
			err := respondJSONError(w, err, http.StatusBadRequest)
			if err != nil {
				s.log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		err = s.service.SendMessage(roomID, messageData.UserID, messageData.Text)
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
