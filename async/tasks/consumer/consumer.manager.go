package async

import (
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
	"github.com/BogoCvetkov/go_mastercalss/mailer"
)

type TaskConsumerManager struct {
	EmailProcessor *EmailProcessor
}

func NewTaskConsumerManager(srv interfaces.IAsyncServer) *TaskConsumerManager {

	mailer := mailer.NewMailer(srv.GetConfig())

	return &TaskConsumerManager{
		EmailProcessor: NewEmailProcessor(mailer, srv),
	}
}
