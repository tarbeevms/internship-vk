package session

import (
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/tarantool/go-tarantool"
)

type Session struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type JwtCustomClaims struct {
	Name string
	jwt.StandardClaims
}

type SessionRepository struct {
	tConn *tarantool.Connection
}

func NewSessionRepo(conn *tarantool.Connection) *SessionRepository {
	return &SessionRepository{
		tConn: conn,
	}
}
