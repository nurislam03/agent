package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Messages struct {
	ID        primitive.ObjectID `json:"id" bson:"id"`
	MessageID string             `json:"messageId" bson:"messageId"`
	TaskID    string             `json:"task_id" bson:"taskId"`
	Body      string             `json:"body" bson:"body"`
	FileRefID string             `json:"file_ref_id" bson:"fileRefId"`
	ActorType string             `json:"actor_type" bson:"actorType"`
	ActorID   string             `json:"actor_id" bson:"actorId"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updatedAt"`
}
