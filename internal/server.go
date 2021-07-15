package internal

import (
	"net/http"
)

type server struct {
	cc     commandController
	qc     queryController
	router http.ServeMux
}

func NewServer(db DB) *server {
	s := &server{
		cc:     newCommandController(db),
		qc:     newQueryController(db),
		router: http.ServeMux{},
	}
	s.router.HandleFunc("/postschedule", s.cc.handlePostSchedule())
	s.router.HandleFunc("/getnextconflict", s.qc.handleGetNextConflict())
	return s
}

// ServeHTTP ... This is used to satisfy http.Handler interface.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
