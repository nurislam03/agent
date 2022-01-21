package memory

import (
	"errors"
	"github.com/nurislam03/agent/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskStore struct {
	tstore map[string]*model.Task
	mstore map[string]*model.Messages
}

func NewTaskStore() *TaskStore {
	id, _ := primitive.ObjectIDFromHex("1001")
	return &TaskStore{
		tstore: map[string]*model.Task{
			"1": &model.Task{
				ID:         id,
				TaskID:     "1",
				Name:       "Assessment",
				StartDate:  time.Now(),
				Category:   "Test",
				Status:     model.TaskStatusProcessing,
				CustomerID: "1",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			"2": &model.Task{
				ID:         id,
				TaskID:     "2",
				Name:       "Dentist Appointment",
				StartDate:  time.Now(),
				Category:   "Health",
				Status:     model.TaskStatusProcessing,
				CustomerID: "2",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			//"3": &model.Task{
			//	ID:         id,
			//	TaskID:     "3",
			//	Name:       "Birthday Party",
			//	StartDate:  time.Now(),
			//	Category:   "Event Management",
			//	Status:     model.TaskStatusDraft,
			//	CustomerID: "3",
			//	CreatedAt:  time.Now(),
			//	UpdatedAt:  time.Now(),
			//},
			//"4": &model.Task{
			//	ID:         id,
			//	TaskID:     "4",
			//	Name:       "Assessment",
			//	StartDate:  time.Now(),
			//	Category:   "Test",
			//	Status:     model.TaskStatusProcessing,
			//	CustomerID: "4",
			//	CreatedAt:  time.Now(),
			//	UpdatedAt:  time.Now(),
			//},
		},
		mstore: map[string]*model.Messages{
			"1": &model.Messages{
				ID:        id,
				MessageID: "1",
				TaskID:    "1",
				Body:      "This is a coding challenge",
				FileRefID: "12345",
				ActorType: string(model.Customer),
				ActorID:   "1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			"2": &model.Messages{
				ID:        id,
				MessageID: "2",
				TaskID:    "1",
				Body:      "Looking into it",
				FileRefID: "12345",
				ActorType: string(model.Operator),
				ActorID:   "1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			"3": &model.Messages{
				ID:        id,
				MessageID: "3",
				TaskID:    "2",
				Body:      "Arrange an event",
				FileRefID: "12345",
				ActorType: string(model.Customer),
				ActorID:   "1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			"4": &model.Messages{
				ID:        id,
				MessageID: "4",
				TaskID:    "2",
				Body:      "This is a coding challenge",
				FileRefID: "12345",
				ActorType: string(model.Operator),
				ActorID:   "2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
}

// GetTaskByID return a specific task by its id
func (s TaskStore) GetTaskByID(tskID string) (*model.Task, error) {
	if s.tstore[tskID] == nil {
		return nil, errors.New("invalid request")
	}
	return s.tstore[tskID], nil
}

//GetTasks returns all data for operator or all related data of an user
func (s TaskStore) GetTasks(uId string, plimt, poffset int) ([]*model.Task, error) {
	ts := []*model.Task{}

	//when uId is empty (operator call) -> return all tasks. Do necessary calculation for pager
	if uId == "" {
		count := 0
		for _, val := range s.tstore {
			if count < poffset {
				count++
				continue
			}
			ts = append(ts, val)
			plimt--
			count++
			if plimt == 0 {
				break
			}
		}
	} else {
		//uId is not empty (customer call) -> return only tasks related to the user.
		count := 0
		for _, val := range s.tstore {
			if val.CustomerID != uId {
				continue
			}

			if count < poffset {
				count++
				continue
			}
			ts = append(ts, val)
			plimt--
			count++
			if plimt == 0 {
				break
			}
		}
	}

	//check whether the request has any data available
	if len(ts) <= 0 {
		return nil, errors.New("no data found")
	}

	return ts, nil
}

//GetMessageByID returns a single message
func (s TaskStore) GetMessageByID(msgID string) (*model.Messages, error) {
	for _, val := range s.mstore {
		if val.MessageID == msgID {
			return val, nil
		}
	}
	return nil, errors.New("invalid message id")
}

//GetChatHistory returns chat history of a specific task
func (s TaskStore) GetChatHistory(tskID string, plimt, poffset int) ([]*model.Messages, error) {
	msglist := []*model.Messages{}

	count := 0
	for _, val := range s.mstore {
		if val.TaskID != tskID {
			continue
		}

		if count < poffset {
			count++
			continue
		}
		msglist = append(msglist, val)
		plimt--
		count++
		if plimt == 0 {
			break
		}
	}
	if len(msglist) <= 0 {
		return nil, errors.New("invalid request")
	}
	return msglist, nil
}

