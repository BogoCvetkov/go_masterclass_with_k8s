package async

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/BogoCvetkov/go_mastercalss/mailer"
	"github.com/hibiken/asynq"
)

type EmailDeliveryPayload struct {
	UserID int64
	Email  string
	Name   string
}

type EmailProcessor struct {
	mailer *mailer.Mailer
}

func NewEmailProcessor(m *mailer.Mailer) *EmailProcessor {
	return &EmailProcessor{mailer: m}
}

// Handle tasks
func (p *EmailProcessor) ProcessTask(c context.Context, t *asynq.Task) error {
	var payload EmailDeliveryPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	fmt.Println("Processing task")

	msg := fmt.Sprintf("Welcome <b>%s</b>. Your account was just created for email <i>%s</i>!", payload.Name, payload.Email)

	err := p.mailer.NewMail(payload.Email, "Account Created", msg)
	if err != nil {
		return err
	}

	fmt.Printf("Welcome email send to ---> %s", payload.Email)

	return nil
}

// Create task
func NewEmailTask(data EmailDeliveryPayload, opts ...asynq.Option) (*asynq.Task, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeEmailDelivery, payload, opts...), nil
}
