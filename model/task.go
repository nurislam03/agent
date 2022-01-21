package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Role string

const (
	Customer Role = "customer"
	Operator Role = "operator"
)

// TaskStatusType ...
type TaskStatusType string

// possible status for a task
const (
	// TaskStatusDraft ...
	TaskStatusDraft TaskStatusType = "draft"
	// TaskStatusPending ...
	TaskStatusPending TaskStatusType = "pending"
	// TaskStatusApproved ...
	TaskStatusApproved TaskStatusType = "approved"
	// TaskStatusProcessing ...
	TaskStatusProcessing TaskStatusType = "processing"
	// TaskStatusDone ...
	TaskStatusDone TaskStatusType = "done"
	// TaskStatusArchive ...
	TaskStatusArchive TaskStatusType = "archive"
	// TaskStatusExpired ...
	TaskStatusExpired TaskStatusType = "expired"
	// TaskStatusFailed ...
	TaskStatusFailed TaskStatusType = "failed"
	// TaskStatusCanceled ...
	TaskStatusCanceled TaskStatusType = "canceled"
)

type Task struct {
	ID         primitive.ObjectID `json:"id" bson:"id"`
	TaskID     string             `json:"task_id" bson:"taskId"`
	Name       string             `json:"name" bson:"name"`
	StartDate  time.Time          `json:"start_date" bson:"startDate"`
	Category   string             `json:"category" bson:"category"`
	Status     TaskStatusType     `json:"status" bson:"status"`
	CustomerID string             `json:"customer_id" bson:"customerId"`
	CreatedAt  time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updatedAt"`
}
