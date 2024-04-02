package stores

import (
	"fmt"
	"log/slog"

	"github.com/fridrock/auth_service/db/entities"
	"github.com/fridrock/auth_service/utils/hashing"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

func CreateUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us UserStore) CreateUser(u entities.User) (int64, error) {
	hash, err := hashing.HashPassword(u.Password)
	if err != nil {
		slog.Error(fmt.Sprintf("UserStore.createUser(): %v", err))
		return 0, err
	}
	u.Password = hash
	id, err := us.createUserQuery(u)
	if err != nil {
		slog.Error(fmt.Sprintf("UserStore.createUser(): %v", err))
		return 0, err
	}
	return id, nil
}

func (us UserStore) createUserQuery(u entities.User) (int64, error) {
	var id int64
	if us.checkIfUserExist(u) {
		return 0, fmt.Errorf("user with username and email:%v, %v already exists", u.Username, u.Email)
	}
	q := `INSERT INTO users (username, email, hashed_password) VALUES($1, $2, $3) RETURNING id`
	err := us.db.QueryRow(
		q,
		u.Username,
		u.Email,
		u.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (us UserStore) checkIfUserExist(u entities.User) bool {
	var user entities.User
	q := `SELECT * FROM users WHERE username=$1 OR email=$2`
	row := us.db.QueryRowx(q, u.Username, u.Email)
	err := row.StructScan(&user)
	return err == nil && user.Username != ""
}

func (us UserStore) GetUsers() ([]entities.User, error) {
	var users []entities.User
	err := us.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		slog.Error(fmt.Sprintf("UserStore.getUsers(): %v", err))
		return nil, err
	}
	return users, nil
}
