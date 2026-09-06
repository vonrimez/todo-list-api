package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vonrimez/TaskAPI/internal/handlers"
)

type TAPIServer struct {
	s   *http.Server
	mux *chi.Mux
}

func NewServer(addr string) *TAPIServer {
	mux := chi.NewRouter()
	return &TAPIServer{&http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}, mux}
}

func (s *TAPIServer) InitRouting(hdl *handlers.Handler) {

	api := s.mux

	api.Use(middleware.Recoverer)
	api.Use(middleware.Logger)

	api.Route("/api/v1", func(v1 chi.Router) {
		v1.Route("/user", func(user chi.Router) {
			user.Post("/login", hdl.UserLogin)
			user.Post("/register", hdl.UserRegister)
		})
		v1.Route("/task", func(task chi.Router) {
			task.Use(hdl.JWTAuth)
			task.Get("/list", hdl.GetTasks)
			task.Get("/{id}", hdl.GetTaskById)
			task.Delete("/{id}", hdl.DeleteTask)
			task.Put("/{id}", hdl.UpdateTask)
			task.Post("/create", hdl.CreateTask)
		})
	})
}

func (s *TAPIServer) Run() error {
	go func() {
		err := s.s.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-stopCtx.Done()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.s.Shutdown(timeoutCtx); err != nil {
		return err
	}

	return nil
}
