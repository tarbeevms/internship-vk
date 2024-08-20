package logic

import (
	"myapp/internal/data"
	"myapp/internal/session"
	"myapp/internal/users"
)

// LogicLayer предоставляет бизнес-логику для работы с данными.
type LogicLayer struct {
	DataRepo    *data.DataRepository
	SessionRepo *session.SessionRepository
	UserRepo    *users.UserRepository
}

// NewLogicLayer создает новый экземпляр слоя логики.
func NewLogicLayer(dr *data.DataRepository, sr *session.SessionRepository, ur *users.UserRepository) *LogicLayer {
	return &LogicLayer{
		DataRepo:    dr,
		SessionRepo: sr,
		UserRepo:    ur,
	}
}
