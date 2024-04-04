package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/fridrock/auth_service/db/core"
	"github.com/fridrock/auth_service/db/stores"
	"github.com/fridrock/auth_service/handlers"
	"github.com/fridrock/auth_service/handlers/users"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type App struct {
	server      *http.Server
	db          *sqlx.DB
	userStore   stores.UserStore
	userService users.UserService
}

func startApp() {
	a := App{}
	a.setup()
}

func (a App) setup() {
	a.db = core.CreateConnection()
	defer a.db.Close()
	a.userStore = *stores.CreateUserStore(a.db)
	a.userService = users.CreateUserService(a.userStore)
	a.setupServer()
}
func (a App) setupServer() {
	a.server = &http.Server{
		Addr:         ":9000",
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		Handler:      a.getRouter(),
	}
	slog.Info("Starting server on port 9000")
	a.server.ListenAndServe()
}
func (a App) getRouter() http.Handler {
	mainRouter := mux.NewRouter()
	mainRouter.Handle("/users/", a.getUsersRouter(mainRouter))
	return mainRouter
}

func (a App) getUsersRouter(r *mux.Router) *mux.Router {
	usersRouter := r.PathPrefix("/users").Subrouter()
	usersRouter.Handle("/signup", handlers.HandleErrorMiddleware(a.userService.CreateUserHandler)).Methods("POST")
	usersRouter.Handle("/signin", handlers.HandleErrorMiddleware(a.userService.AuthUserHandler)).Methods("POST")
	usersRouter.Handle("/logout", handlers.HandleErrorMiddleware(a.userService.LogoutUserHandler)).Methods("POST")
	return usersRouter
}
