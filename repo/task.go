package repo

import "github.com/nurislam03/agent/model"

type Task interface {
	GetTasks(uId string, plimt, poffset int) ([]*model.Task, error)
}
