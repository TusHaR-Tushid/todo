package server

import (
	"github.com/go-chi/chi/v5"
	"go-todo/handler"
	"net/http"
)

type Server struct {
	chi.Router
}

func SetupRoutes() *Server {
	router := chi.NewRouter()

	router.Route("/todo", func(todo chi.Router) {
		todo.Post("/register", handler.Register)
		todo.Post("/sign-in", handler.Login)
		todo.Route("/home", func(api chi.Router) {
			api.Use(handler.Middleware)
			api.Post("/todo", handler.CreateTodo)
			api.Get("/all-todo", handler.GetAll)
			api.Get("/completed-todo", handler.GetCompleted)
			api.Get("/upcoming-todo", handler.GetUpcoming)
			api.Get("/expired-todo", handler.GetExpired)
			api.Put("/update-todo", handler.UpdateTodo)
			api.Delete("/delete-todo", handler.DeleteTodo)
		})

	})
	return &Server{router}
}

func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
