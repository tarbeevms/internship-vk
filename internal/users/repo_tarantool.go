package users

import (
	"errors"

	"github.com/tarantool/go-tarantool"
)

// GetUserByUsername ищет пользователя по username в Tarantool
func (ur *UserRepository) GetUserByUsername(username string) (*LoginRequest, error) {
	resp, err := ur.tConn.Select("users", "primary", 0, 1, tarantool.IterEq, []interface{}{username})
	if err != nil {
		return nil, err
	}

	if len(resp.Tuples()) == 0 {
		return nil, errors.New("user not found")
	}
	userTuple := resp.Tuples()[0]
	user := &LoginRequest{
		Username: userTuple[0].(string),
		Password: userTuple[1].(string),
	}

	return user, nil
}
