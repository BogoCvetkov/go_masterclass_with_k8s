package db

import (
	"context"
	"errors"
	"time"

	async_pay "github.com/BogoCvetkov/go_mastercalss/async/tasks/payload"
	async_prod "github.com/BogoCvetkov/go_mastercalss/async/tasks/producer"
	db "github.com/BogoCvetkov/go_mastercalss/db/generated"
	"github.com/hibiken/asynq"
)

func (s *Store) CreateUserTrx(c context.Context, data db.CreateUserParams, asq *asynq.Client) (*db.User, error) {
	// Begin Trx --->
	tx, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	qtx := s.WithTx(tx)

	// Create new user
	user, err := qtx.CreateUser(c, data)
	if err != nil {

		if ErrorCode(err) == UniqueViolation {
			return nil, errors.New("email already exists")
		}

		return nil, err
	}

	// Send email-verification link
	mailInfo := async_pay.VerifyEmailPayload{
		UserID: user.ID,
		Email:  user.Email,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.ProcessIn(5 * time.Second),
		asynq.Queue("critical"),
	}

	mailTask, err := async_prod.NewEmailTask(mailInfo, opts...)
	if err != nil {
		return nil, errors.New("failed to create email-task")
	}
	_, err = asq.Enqueue(mailTask)
	if err != nil {
		return nil, errors.New("failed to send email")
	}

	// Commit Trx <---
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &user, nil
}
