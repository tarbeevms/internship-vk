package logic

import (
	"fmt"
	"myapp/config"
	"myapp/internal/session"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func (ll *LogicLayer) IsAuthorized(requestToken string) (bool, error) {
	// Распарсить и проверить JWT токен
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.CFG.SecretKey), nil
	})

	// Ошибка парсинга токена
	if err != nil && err.Error() != "Token is expired" {
		return false, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid token")
	}

	// Проверка существования сессии в базе данных
	session, err := ll.SessionRepo.GetSessionByToken(requestToken)
	if err != nil || session == nil {
		return false, err
	}

	// Проверка, не истек ли срок действия токена
	exp, ok := claims["exp"].(float64) // exp в Unix формате
	if !ok {
		return false, fmt.Errorf("invalid token, missing expiration")
	}

	// Если токен истек, удаляем сессию из базы данных
	if time.Now().Unix() > int64(exp) {
		err := ll.SessionRepo.DeleteSessionByToken(requestToken)
		if err != nil {
			return false, fmt.Errorf("failed to delete expired session, %w", err)
		}
		return false, fmt.Errorf("session expired")
	}

	// Токен и сессия действительны
	return true, nil
}

// CreateSession генерирует токен и добавляет сессию в базу данных.
func (ll *LogicLayer) CreateSession(username string) (string, error) {
	// Генерация токена доступа (JWT)
	accessToken, err := CreateAccessToken(username)
	if err != nil {
		return "", fmt.Errorf("failed to create access token, %w", err)
	}

	// Создание сессии
	newSession := &session.Session{
		Username: username,
		Token:    accessToken,
	}
	// Сохранение сессии в базе данных
	err = ll.SessionRepo.AddSession(newSession)
	if err != nil {
		return "", fmt.Errorf("failed to store session, %w", err)
	}
	return accessToken, nil
}

func CreateAccessToken(username string) (accessToken string, err error) {
	exp := time.Now().Add(180 * time.Second).Unix()
	claims := &session.JwtCustomClaims{
		Name: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.CFG.SecretKey))
	if err != nil {
		return "", err
	}
	return t, err
}
