package async

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	payload "github.com/BogoCvetkov/go_mastercalss/async/tasks/payload"
	gen "github.com/BogoCvetkov/go_mastercalss/db/generated"
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
	"github.com/BogoCvetkov/go_mastercalss/mailer"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type EmailProcessor struct {
	mailer *mailer.Mailer
	srv    interfaces.IAsyncServer
}

func NewEmailProcessor(m *mailer.Mailer, srv interfaces.IAsyncServer) *EmailProcessor {
	return &EmailProcessor{mailer: m, srv: srv}
}

// Handle tasks
func (p *EmailProcessor) ProcessTask(c context.Context, t *asynq.Task) error {
	var payload payload.VerifyEmailPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	fmt.Println("Processing task")

	// Find user
	user, err := p.srv.GetStore().GetUser(c, payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}

		return fmt.Errorf(err.Error())
	}

	// Generate verification code
	verification, err := p.srv.GetStore().CreateVerifyEmail(c, gen.CreateVerifyEmailParams{
		UserID:     user.ID,
		Email:      user.Email,
		SecretCode: fmt.Sprintf("%d_%s", time.Now().Unix(), uuid.NewString()),
	})

	if err != nil {
		return err
	}

	// Build verification link
	link := fmt.Sprintf("localhost:%s/api/user/verify?email=%s&code=%s", p.srv.GetConfig().Port, verification.Email, verification.SecretCode)

	// HTML message
	msg := fmt.Sprintf("Welcome <b>%s</b>. Your account was just created for email <i>%s</i>!.You can verify at this link - <a href=%s> Verify </href>  ", user.FullName, user.Email, link)

	err = p.mailer.NewMail(user.Email, "Verify your email", msg)
	if err != nil {
		return err
	}

	fmt.Printf("Verify email send to ---> %s", user.Email)

	return nil
}
