package users

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"

	"github.com/fridrock/auth_service/db/entities"
	"github.com/fridrock/auth_service/db/stores"
	mailService "github.com/fridrock/auth_service/utils/mail"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(w http.ResponseWriter, r *http.Request) (status int, err error)
	LogoutUser(w http.ResponseWriter, r *http.Request) (status int, err error)
	SendConfirmation(w http.ResponseWriter, r *http.Request) (status int, err error)
	AuthUser(w http.ResponseWriter, r *http.Request) (status int, err error)
}

type UserServiceImpl struct {
	store      stores.UserStore
	cacheStore stores.CacheStore
}

func CreateUserService(store stores.UserStore, cacheStore stores.CacheStore) *UserServiceImpl {
	return &UserServiceImpl{
		store:      store,
		cacheStore: cacheStore,
	}
}
func (us *UserServiceImpl) CreateUser(w http.ResponseWriter, r *http.Request) (status int, err error) {
	user, err := parseUser(r)
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
func parseUser(r *http.Request) (entities.User, error) {
	var usr entities.User
	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		return usr, err
	}
	_, err = mail.ParseAddress(usr.Email)
	if err != nil {
		return usr, err
	}
	return usr, nil
}

type EmailConfirmationRequest struct {
	Id int64 `json:"id"`
}

func (us *UserServiceImpl) SendConfirmation(w http.ResponseWriter, r *http.Request) (status int, err error) {
	var emr EmailConfirmationRequest
	err = json.NewDecoder(r.Body).Decode(&emr)
	if err != nil {
		return http.StatusBadRequest, err
	}
	userEmail, err := us.store.GetUserEmailById(emr.Id)
	if err != nil {
		return http.StatusNotFound, err
	}
	confirmationCode := uuid.New()
	err = us.cacheStore.PutEmailConfirmation(confirmationCode.String(), emr.Id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	mailService.Send(confirmationCode.String(), userEmail)
	return http.StatusOK, nil
}
func (us *UserServiceImpl) AuthUser(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return 200, nil
}
func (us *UserServiceImpl) CheckChat(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return 200, nil
}
func (us *UserServiceImpl) LogoutUser(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return 200, nil
}
