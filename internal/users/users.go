package users

import "github.com/tarantool/go-tarantool"

type UserRepository struct {
	tConn *tarantool.Connection
}

func NewUserRepo(conn *tarantool.Connection) *UserRepository {
	return &UserRepository{
		tConn: conn,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
