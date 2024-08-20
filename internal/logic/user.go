package logic

import (
	"myapp/internal/users"
)

func (ll *LogicLayer) VerifyUserCredentials(username, password string) (*users.LoginRequest, error) {
	// Проверяем учетные данные пользователя
	user, err := ll.UserRepo.GetUserByUsername(username)
	if err != nil || user.Password != password {
		return nil, err
	}
	return user, nil
}
