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
	if us.checkIfUserExist(u) {
		return 0, fmt.Errorf("user with username and email:%v, %v already exists", u.Username, u.Email)
	}
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
	q := `INSERT INTO users (username, email, hashed_password) VALUES($1, $2, $3) RETURNING id`
	err := us.db.QueryRow(
		q,
		u.Username,
		u.Email,
		u.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	if err = us.setStatusToUser(id, "UNCONFIRMED"); err != nil {
		us.deleteUserById(id)
		return 0, err
	}
	return id, nil
}
func (us UserStore) setStatusToUser(userId int64, statusName string) error {
	statusId, err := us.getStatusByName(statusName)
	if err != nil {
		return err
	}
	var resultId int64
	q := `INSERT INTO users_statuses (user_id, status_id) VALUES ($1, $2) RETURNING user_id`
	err = us.db.QueryRow(q, userId, statusId).Scan(&resultId)
	if err != nil {
		return err
	}
	return nil
}
func (us UserStore) UpdateUserStatus(userId int64, statusName string) error {
	statusId, err := us.getStatusByName(statusName)
	if err != nil {
		return err
	}
	var resultId int64
	q := `UPDATE users_statuses SET status_id = $2 WHERE user_id=$1;`
	err = us.db.QueryRow(q, userId, statusId).Scan(&resultId)
	if err != nil {
		return err
	}
	return nil
}
func (us UserStore) getStatusByName(name string) (int64, error) {
	var id int64
	q := `SELECT (id) FROM user_statuses WHERE val=$1`
	row := us.db.QueryRow(q, name)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, fmt.Errorf("no such status with name: %v", name)
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
func (us UserStore) deleteUserById(id int64) error {
	var deletedId int64
	q := `DELETE FROM users WHERE id=$1 RETURNING id`
	row := us.db.QueryRow(q, id)
	err := row.Scan(&deletedId)
	if err != nil {
		return err
	}
	return nil
}

func (us UserStore) GetUserEmailById(id int64) (string, error) {
	var userEmail string
	q := `SELECT (email) FROM users WHERE id=$1`
	row := us.db.QueryRow(q, id)
	err := row.Scan(&userEmail)
	if err != nil {
		return "", fmt.Errorf("no such user with id: %v", id)
	}
	return userEmail, nil
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
