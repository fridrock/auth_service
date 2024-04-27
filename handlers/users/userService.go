package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"

	"github.com/fridrock/auth_service/db/entities"
	"github.com/fridrock/auth_service/db/stores"
	"github.com/fridrock/auth_service/utils/hashing"
	mailService "github.com/fridrock/auth_service/utils/mail"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	EMAIL_CONFIRMATION = "email_confirmation"
	START_CODE         = "start_code"
)

type UserService interface {
	CreateUser(w http.ResponseWriter, r *http.Request) (status int, err error)
	LogoutUser(w http.ResponseWriter, r *http.Request) (status int, err error)
	SendConfirmation(w http.ResponseWriter, r *http.Request) (status int, err error)
	ConfirmEmail(w http.ResponseWriter, r *http.Request) (status int, err error)
	AuthUser(w http.ResponseWriter, r *http.Request) (status int, err error)
	GetUser(w http.ResponseWriter, r *http.Request) (status int, err error)
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
	err = us.cacheStore.PutUserId(EMAIL_CONFIRMATION, confirmationCode.String(), emr.Id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	mailService.Send(confirmationCode.String(), userEmail)
	return http.StatusOK, nil
}

func (us *UserServiceImpl) ConfirmEmail(w http.ResponseWriter, r *http.Request) (status int, err error) {
	code := mux.Vars(r)["code"]
	userId, err := us.cacheStore.GetUserId(EMAIL_CONFIRMATION, code)
	if err != nil {
		return http.StatusGone, fmt.Errorf("no such confirmation code")
	}
	err = us.store.UpdateUserStatus(userId, "CONFIRMED")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	startCode, err := us.setStartCode(userId)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Write([]byte(startCode))
	return http.StatusOK, nil
}

func (us *UserServiceImpl) AuthUser(w http.ResponseWriter, r *http.Request) (status int, err error) {
	var usr entities.User
	err = json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		return http.StatusBadRequest, err
	}
	userFromDb, err := us.store.GetUserByUsernameOrEmail(usr)
	if err != nil {
		return http.StatusNotFound, errors.New("no such user")
	}
	if !hashing.CheckPassword(usr.Password, userFromDb.Password) {
		return http.StatusForbidden, errors.New("wrong password")
	}
	if !us.store.CheckConfirmed(userFromDb.Id) {
		return http.StatusForbidden, errors.New("unconfirmed email")
	}
	startCode, err := us.setStartCode(userFromDb.Id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Write([]byte(startCode))
	return http.StatusOK, nil
}
func (us *UserServiceImpl) setStartCode(userId int64) (string, error) {
	startCode := uuid.New().String()
	return startCode, us.cacheStore.PutUserId(START_CODE, startCode, userId)
}

type UserIdResponse struct {
	UserId int64 `json:"user_id"`
}

func (us *UserServiceImpl) GetUser(w http.ResponseWriter, r *http.Request) (status int, err error) {
	startCode := mux.Vars(r)["startCode"]
	userId, err := us.cacheStore.GetUserId(START_CODE, startCode)
	if err != nil {
		return http.StatusGone, fmt.Errorf("no such start code")
	}
	response := UserIdResponse{
		UserId: userId,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
func (us *UserServiceImpl) LogoutUser(w http.ResponseWriter, r *http.Request) (status int, err error) {
	return 200, nil
}
