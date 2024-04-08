package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/fridrock/auth_service/db/cache"
	"github.com/fridrock/auth_service/db/core"
	"github.com/fridrock/auth_service/db/stores"
	"github.com/fridrock/auth_service/handlers"
	"github.com/fridrock/auth_service/handlers/users"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type App struct {
	server      *http.Server
	db          *sqlx.DB
	redisClient *redis.Client
	userStore   stores.UserStore
	cacheStore  stores.CacheStore
	userService users.UserService
}

func startApp() {
	a := App{}
	a.setup()
}

func (a App) setup() {
	a.db = core.CreateConnection()
	defer a.db.Close()
	a.redisClient = cache.CreateRedisClient()
	defer a.redisClient.Close()
	a.userStore = *stores.CreateUserStore(a.db)
	a.cacheStore = *stores.CreateCacheStore(a.redisClient)
	a.userService = users.CreateUserService(a.userStore, a.cacheStore)
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
	usersRouter.Handle("/signup", handlers.HandleErrorMiddleware(a.userService.CreateUser)).Methods("POST")
	usersRouter.Handle("/send-confirmation", handlers.HandleErrorMiddleware(a.userService.SendConfirmation)).Methods("POST")
	usersRouter.Handle("/signin", handlers.HandleErrorMiddleware(a.userService.AuthUser)).Methods("POST")
	usersRouter.Handle("/logout", handlers.HandleErrorMiddleware(a.userService.LogoutUser)).Methods("POST")
	return usersRouter
}
