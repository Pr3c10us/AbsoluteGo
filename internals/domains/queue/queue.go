package queue

import (
	"github.com/google/uuid"
)

type TaskQueueMessage struct {
	UserSlackId      string
	WorkspaceSlackId string
	WorkflowId       uuid.UUID
	CronJobID        uuid.UUID
}

type SendSlackMessage struct {
	Message string
	Token   string
	SlackID string
}

type MessageParams struct {
	Queue   string
	Message string
	Key     string
}
