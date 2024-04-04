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
	CreateUserHandler(w http.ResponseWriter, r *http.Request) (status int, err error)
	LogoutUserHandler(w http.ResponseWriter, r *http.Request) (status int, err error)
	AuthUserHandler(w http.ResponseWriter, r *http.Request) (status int, err error)
}

type UserServiceImpl struct {
	store stores.UserStore
}

func CreateUserService(store stores.UserStore) *UserServiceImpl {
	return &UserServiceImpl{
		store: store,
	}
}
func (us *UserServiceImpl) CreateUserHandler(w http.ResponseWriter, r *http.Request) (status int, err error) {
	user, err := parseUser(w, r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	id, err := us.store.CreateUser(user)
	if err != nil {
		return http.StatusBadRequest, err
	}
	slog.Info(fmt.Sprintf("Created user with id: %v", id))
	w.Write([]byte(fmt.Sprintf("id is : %v", id)))
	return http.StatusOK, nil
}
func parseUser(w http.ResponseWriter, r *http.Request) (entities.User, error) {
	var usr entities.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		return usr, err
	}
	return usr, nil
}

func (us *UserServiceImpl) AuthUserHandler(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return 200, nil
}
func (us *UserServiceImpl) CheckChatHandler(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return 200, nil
}
func (us *UserServiceImpl) LogoutUserHandler(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return 200, nil
}
