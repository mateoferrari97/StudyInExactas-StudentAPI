package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const defaultPort = ":8081"

type Server struct {
	Router *mux.Router
}

func NewServer() *Server {
	return &Server{
		Router: mux.NewRouter(),
	}
}

func (s *Server) Run(port string) error {
	port = setupPort(port)

	log.Printf("Listening on port %s", port)
	return http.ListenAndServe(port, s.Router)
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (s *Server) Wrap(method string, pattern string, handler HandlerFunc) {
	wrapH := func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		hErr := handleError(w, err)
		_ = RespondJSON(w, hErr, hErr.StatusCode)
	}

	s.Router.HandleFunc(pattern, wrapH).Methods(method)
}

func setupPort(port string) string {
	if port == "" {
		port = defaultPort
		log.Printf("Defaulting to port %s", port)
	}

	return port
}
