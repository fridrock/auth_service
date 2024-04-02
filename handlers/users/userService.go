package users

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fridrock/auth_service/db/entities"
	"github.com/fridrock/auth_service/db/stores"
)

type UserService interface {
	CreateUserHandler(w http.ResponseWriter, r *http.Request)
	LogoutUserHandler(w http.ResponseWriter, r *http.Request)
	AuthUserHandler(w http.ResponseWriter, r *http.Request)
}

type UserServiceImpl struct {
	store stores.UserStore
}

func CreateUserService(store stores.UserStore) *UserServiceImpl {
	return &UserServiceImpl{
		store: store,
	}
}
func (us *UserServiceImpl) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := parseUser(r)
	if err != nil {
		writeError(w, err)
		return
	}
	id, err := us.store.CreateUser(user)
	if err != nil {
		writeError(w, err)
		return
	}
	slog.Info(fmt.Sprintf("Created user with id: %v", id))
	w.Write([]byte(fmt.Sprintf("id is : %v", id)))
}
func parseUser(r *http.Request) (entities.User, error) {
	var usr entities.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		return usr, err
	}
	return usr, err
}
func writeError(w http.ResponseWriter, err error) {
	slog.Error(fmt.Sprintf("CreateUserHandler.handleHttp(): %v", err))
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprint(err)))
}
func (us *UserServiceImpl) AuthUserHandler(w http.ResponseWriter, r *http.Request) {

}
func (us *UserServiceImpl) CheckChatHandler(w http.ResponseWriter, r *http.Request) {

}
func (us *UserServiceImpl) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {

}
