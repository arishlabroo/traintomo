package internal

import (
	"net/http"
)

// Server represents the http server.
type Server struct {
	cc     commandController
	qc     queryController
	router http.ServeMux
}

// NewServer returns a http server where the application routes have already been configured.
func NewServer(db DB) *Server {
	s := &Server{
		cc:     newCommandController(db),
		qc:     newQueryController(db),
		router: http.ServeMux{},
	}
	s.router.HandleFunc("/postschedule", s.cc.handlePostSchedule())
	s.router.HandleFunc("/getnextconflict", s.qc.handleGetNextConflict())
	return s
}

// ServeHTTP ... This is used to satisfy http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
