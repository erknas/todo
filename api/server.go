package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/zeze322/todo/db"
	"github.com/zeze322/todo/lib"
)

type Server struct {
	password string
	port     string
	store    db.Storage
}

func NewServer(port, password string, store db.Storage) *Server {
	return &Server{
		port:     port,
		password: password,
		store:    store,
	}
}

func (s *Server) Run() error {
	router := chi.NewMux()

	router.Handle("/*", http.FileServer(http.Dir("./web")))

	router.HandleFunc("/api/signin", lib.MakeHTTP(s.handleSign))
	router.HandleFunc("/api/tasks", withJWTAuth(lib.MakeHTTP(s.handleTask), s.password))
	router.HandleFunc("/api/task", withJWTAuth(lib.MakeHTTP(s.handleTask), s.password))
	router.Get("/api/task", withJWTAuth(lib.MakeHTTP(s.handleGetTaskByID), s.password))
	router.Post("/api/task/done", withJWTAuth(lib.MakeHTTP(s.handleTaskDone), s.password))
	router.Get("/api/nextdate", lib.MakeHTTP(s.handleNextDate))

	log.Printf("Starting server on port %s", s.port)

	if err := http.ListenAndServe(s.port, router); err != nil {
		return fmt.Errorf("failed to start server")
	}

	return nil
}
