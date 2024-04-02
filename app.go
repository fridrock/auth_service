package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/fridrock/auth_service/db/core"
	"github.com/fridrock/auth_service/db/stores"
	"github.com/fridrock/auth_service/handlers/users"
	"github.com/fridrock/auth_service/utils"
	"github.com/gorilla/mux"
)

//Service funcionality:
// create account
// auth account => get special unique code => save in database unique_code -> username
// send link to bot with this special code: start=unique_code
// find user by unique_code, save message.chat.id->username
// on next query compare message.chat.id ->TODO add this to cache

// Main facade
type App struct {
	server      *http.Server
	db          *sql.DB
	userStore   stores.UserStore
	userService users.UserService
}

func startApp() {
	a := App{}
	utils.Send("Privet")
	a.setup()
}

func (a App) setup() {
	a.db = core.CreateConnection()
	defer a.db.Close()
	a.userStore = *stores.CreateUserStore(a.db)
	a.userService = users.CreateUserService(a.userStore)
	a.server = &http.Server{
		Addr:         ":9000",
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		Handler:      a.getRouter(),
	}
	a.server.ListenAndServe()
}
func (a App) getRouter() http.Handler {
	mainRouter := mux.NewRouter()
	mainRouter.Handle("/users/", a.getUsersRouter(mainRouter))
	return mainRouter
}
func (a App) getUsersRouter(r *mux.Router) *mux.Router {
	usersRouter := r.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/signup", a.userService.CreateUserHandler).Methods("POST")
	usersRouter.HandleFunc("/signin", a.userService.AuthUserHandler).Methods("POST")
	usersRouter.HandleFunc("/logout", a.userService.LogoutUserHandler).Methods("POST")
	return usersRouter
}
