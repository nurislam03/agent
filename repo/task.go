package repo

import "github.com/nurislam03/agent/model"

type Task interface {
	GetTasks(uId string, plimt, poffset int) ([]*model.Task, error)
	GetTaskByID(tskID string) (*model.Task, error)
	GetChatHistory(tskID string, plimt, poffset int) ([]*model.Messages, error)
	GetMessageByID(msgID string) (*model.Messages, error)
}
