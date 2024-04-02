package entities

type User struct {
	Id       uint64 `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"hashed_password"`
}
