package main

import (
	"fmt"
	"log"
	"net/http"
)

type server struct {
	router *http.ServeMux
	db     database
}

func (s *server) checkJSON(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-type") != "application/json" {
			http.Error(w, "No JSON", http.StatusBadRequest)
			return
		}
		h(w, r)
	}
}

func (s *server) checkPOST(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

func (s *server) writeWeather() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws := weatherSet{}
		if err := ws.fromJSON(r.Body); err != nil {
			log.Printf("Unable to parse weather data: %v\n", err)
			http.Error(w, "Unable to parse weather data", http.StatusBadRequest)
			return
		}

		if err := s.db.add(ws); err != nil {
			log.Printf("Unable to store weather data: %v\n", err)
			http.Error(w, "Unable to store weather data", http.StatusBadRequest)
			return
		}
	}
}

func (s *server) routes() {
	s.router = http.NewServeMux()
	s.router.HandleFunc("/write/weather",
		s.checkPOST(
			s.checkJSON(
				s.writeWeather())))
}

func (s *server) run() error {
	s.routes()

	if err := s.db.init(); err != nil {
		return fmt.Errorf("unable to init database, %v", err)
	}
	defer s.db.close()

	go s.db.run()

	return http.ListenAndServe(":8080", s.router)
}
