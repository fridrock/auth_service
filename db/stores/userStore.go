package stores

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/fridrock/auth_service/db/entities"
)

type UserStore struct {
	db *sql.DB
}

func CreateUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us UserStore) CreateUser(u entities.User) (int64, error) {
	var id int64
	hash, err := hashPassword(u.Password)
	if err != nil {
		slog.Error(fmt.Sprintf("UserStore.createUser(): %v", err))
		return 0, err
	}
	err = us.db.QueryRow(
		"INSERT INTO users (username, email, hashed_password) VALUES($1, $2, $3) RETURNING id",
		u.Username,
		u.Email,
		hash).Scan(&id)
	if err != nil {
		slog.Error(fmt.Sprintf("UserStore.createUser(): %v", err))
		return 0, err
	}
	return id, nil
}

func (us UserStore) GetUsers() ([]entities.User, error) {
	var users []entities.User
	rows, err := us.db.Query("SELECT * FROM users")
	if err != nil {
		slog.Error(fmt.Sprintf("UserStore.getUsers(): %v", err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Password); err != nil {
			slog.Error(fmt.Sprintf("UserStore.getUsers():%v", err))
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
