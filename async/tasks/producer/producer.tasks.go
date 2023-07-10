package async

import (
	"encoding/json"

	payload "github.com/BogoCvetkov/go_mastercalss/async/tasks/payload"
	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeEmailDelivery = "email:deliver"
)

// Create task
func NewEmailTask(data payload.VerifyEmailPayload, opts ...asynq.Option) (*asynq.Task, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeEmailDelivery, payload, opts...), nil
}
