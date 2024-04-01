package users

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"user_service/db/entities"
	"user_service/db/stores"
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
	var usr entities.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		slog.Error(fmt.Sprintf("handle create user: %v", err))
	}
	id, err := us.store.CreateUser(usr)
	if err != nil {
		slog.Error(fmt.Sprintf("handle create user: %v", err))
	}
	w.Write([]byte(fmt.Sprintf("id is : %v", id)))
}
func (us *UserServiceImpl) AuthUserHandler(w http.ResponseWriter, r *http.Request) {

}
func (us *UserServiceImpl) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {

}
