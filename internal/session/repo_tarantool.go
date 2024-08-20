package session

import (
	"fmt"

	"github.com/tarantool/go-tarantool"
)

// AddSession добавляет новую сессию в Tarantool.
func (sr *SessionRepository) AddSession(sess *Session) error {
	_, err := sr.tConn.ReplaceAsync("sessions", []interface{}{sess.Username, sess.Token}).Get()
	if err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}
	return nil
}

// DeleteSessionByToken удаляет сессию по токену.
func (sr *SessionRepository) DeleteSessionByToken(token string) error {
	_, err := sr.tConn.DeleteAsync("sessions", "token_index", []interface{}{token}).Get()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (sr *SessionRepository) GetSessionByToken(token string) (*Session, error) {
	resp, err := sr.tConn.SelectAsync("sessions", "token_index", 0, 1, tarantool.IterEq, []interface{}{token}).Get()
	if err != nil {
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	tuples := resp.Tuples()

	if len(tuples) == 0 {
		return nil, fmt.Errorf("session not found")
	}

	session := &Session{
		Username: tuples[0][0].(string),
		Token:    tuples[0][1].(string),
	}

	return session, nil
}
