package async

import (
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/mailer"
)

// A list of task types.
const (
	TypeEmailDelivery = "email:deliver"
)

type AsyncManager struct {
	EmailProcessor *EmailProcessor
}

func NewAsyncManager(c *config.Config) *AsyncManager {

	mailer := mailer.NewMailer(c)

	return &AsyncManager{
		EmailProcessor: NewEmailProcessor(mailer),
	}
}
