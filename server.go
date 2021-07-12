package main

import (
	"fmt"
	"net/http"
)

type server struct {
	db     DB
	router http.ServeMux
}

func NewServer() *server {
	s := &server{
		db:     NewInMemoryDB(),
		router: http.ServeMux{},
	}
	s.router.HandleFunc("/postschedule", s.handlePostSchedule())
	s.router.HandleFunc("/getnextconflict", s.handleGetNextConflict())
	return s
}

func (s *server) handlePostSchedule() http.HandlerFunc {
	type request struct {
		Name string
	}
	type response struct {
		Greeting string `json:"greeting"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello handlePostSchedule")
	}
}

func (s *server) handleGetNextConflict() http.HandlerFunc {
	type response struct {
		Greeting string `json:"greeting"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello handleGetNextConflict")
	}
}

// ServeHTTP ... This is used to satisfy http.Handler interface.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
