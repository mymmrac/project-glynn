package httpapi

import "net/http"

const uuidRegx = "\\b[0-9a-f]{8}\\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\\b[0-9a-f]{12}\\b"

func (s *Server) routes() {
	api := s.router.PathPrefix("/api/rooms").Subrouter()

	api.HandleFunc("/{roomID:"+uuidRegx+"}/messages", s.getMessages()).
		Methods(http.MethodGet)
	api.HandleFunc("/{roomID:"+uuidRegx+"}/messages", s.getMessages()).
		Queries("lastMessageID", "{lastMessageID:"+uuidRegx+"}").Methods(http.MethodGet)
	api.HandleFunc("/{roomID:"+uuidRegx+"}/messages", s.sendMassage()).
		Methods(http.MethodPost)
}
