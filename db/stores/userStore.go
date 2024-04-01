package stores

import (
	"database/sql"
	"fmt"
	"log/slog"
	"user_service/db/entities"
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
	err := us.db.QueryRow("INSERT INTO users (username, hashed_password) VALUES($1, $2) RETURNING id", u.Username, u.Password).Scan(&id)
	if err != nil {
		slog.Error(fmt.Sprintf("createUser(): %v", err))
		return 0, err
	}
	return id, nil
}
func (us UserStore) GetUsers() ([]entities.User, error) {
	var users []entities.User
	rows, err := us.db.Query("SELECT * FROM users")
	if err != nil {
		slog.Error(fmt.Sprintf("getUsers(): %v", err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Password); err != nil {
			slog.Error(fmt.Sprintf("getUsers():%v", err))
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
