package api

import (
	"github.com/kinbiko/jsonassert"
	"github.com/nurislam03/agent/repo/memory"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTasks(t *testing.T) {
	api := NewAPI(nil, memory.NewTaskStore())
	testdata := []struct {
		des    string
		role   string
		userID string
		code   int
		resp   string
	}{
		{
			des:    "successfully get tasks with operator",
			role:   "operator",
			userID: "1",
			code:   http.StatusOK,
			resp:   `{"data":[{"id":"000000000000000000000000","task_id":"1","name":"Assessment","start_date":"<<PRESENCE>>","category":"Test","status":"processing","customer_id":"1","created_at":"<<PRESENCE>>","updated_at":"<<PRESENCE>>"},{"id":"000000000000000000000000","task_id":"2","name":"Dentist Appointment","start_date":"<<PRESENCE>>","category":"Health","status":"processing","customer_id":"2","created_at":"<<PRESENCE>>","updated_at":"<<PRESENCE>>"}],"meta":{"current_page":1,"per_page":2,"total_page":0,"total":0}}`,
		},
		{
			des:    "successfully get tasks with customer",
			role:   "customer",
			userID: "1",
			code:   http.StatusOK,
			resp:   `{"data":[{"id":"000000000000000000000000","task_id":"1","name":"Assessment","start_date":"<<PRESENCE>>","category":"Test","status":"processing","customer_id":"1","created_at":"<<PRESENCE>>","updated_at":"<<PRESENCE>>"}],"meta":{"current_page":1,"per_page":2,"total_page":0,"total":0}}`,
		},
		{
			des:    "invalid role",
			role:   "abc",
			userID: "1",
			code:   http.StatusUnprocessableEntity,
			resp:   `{"errors":[{"id":"<<PRESENCE>>","code":"422001","detail":"invalid header","status":422,"title":"unprocessable entity"}]}`,
		},
	}

	for _, td := range testdata {
		t.Run(td.des, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
			req.Header.Add("Role", td.role)
			req.Header.Add("UserID", td.userID)
			res := httptest.NewRecorder()
			api.router.ServeHTTP(res, req)
			assert.Equal(t, td.code, res.Code)
			jsonassert.New(t).Assertf(res.Body.String(), td.resp)
		})
	}
}

func TestGetChatHistory(t *testing.T) {
	api := NewAPI(nil, memory.NewTaskStore())
	testdata := []struct {
		des    string
		role   string
		userID string
		taskID string
		code   int
		resp   string
	}{
		{
			des:    "successfully get chat history by operator",
			role:   "operator",
			userID: "1",
			taskID: "1",
			code:   http.StatusOK,
			resp:   `{"data":[{"id":"000000000000000000000000","messageId":"1","task_id":"1","body":"This is a coding challenge","file_ref_id":"https://upload.wikimedia.org/wikipedia/commons/f/f9/Flag_of_Bangladesh.svg","actor_type":"customer","actor_id":"1","created_at":"<<PRESENCE>>","updated_at":"<<PRESENCE>>"},{"id":"000000000000000000000000","messageId":"2","task_id":"1","body":"Looking into it","file_ref_id":"https://upload.wikimedia.org/wikipedia/commons/f/f9/Flag_of_Bangladesh.svg","actor_type":"operator","actor_id":"1","created_at":"<<PRESENCE>>","updated_at":"<<PRESENCE>>"}],"meta":{"current_page":1,"per_page":10,"total_page":0,"total":0}}`,
		},
		{
			des:    "successfully get chat history by customer",
			role:   "customer",
			userID: "1",
			taskID: "1",
			code:   http.StatusOK,
			resp:   `{"data":[{"id":"000000000000000000000000","messageId":"1","task_id":"1","body":"This is a coding challenge","file_ref_id":"https://upload.wikimedia.org/wikipedia/commons/f/f9/Flag_of_Bangladesh.svg","actor_type":"customer","actor_id":"1","created_at":"<<PRESENCE>>","updated_at":"<<PRESENCE>>"},{"id":"000000000000000000000000","messageId":"2","task_id":"1","body":"Looking into it","file_ref_id":"https://upload.wikimedia.org/wikipedia/commons/f/f9/Flag_of_Bangladesh.svg","actor_type":"operator","actor_id":"1","created_at":"<<PRESENCE>>","updated_at":"<<PRESENCE>>"}],"meta":{"current_page":1,"per_page":10,"total_page":0,"total":0}}`,
		},
		{
			des:    "invalid role",
			role:   "abc",
			userID: "1",
			taskID: "1",
			code:   http.StatusUnprocessableEntity,
			resp:   `{"errors":[{"id":"<<PRESENCE>>","code":"422001","detail":"invalid header","status":422,"title":"unprocessable entity"}]}`,
		},
		{
			des:    "invalid user id when the role is customer",
			role:   "customer",
			userID: "123456",
			taskID: "1",
			code:   http.StatusForbidden,
			resp:   `{"errors":[{"id":"<<PRESENCE>>","code":"403000","status":403,"title":"forbidden request"}]}`,
		},
	}

	for _, td := range testdata {
		t.Run(td.des, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+td.taskID+"/chats", nil)
			req.Header.Add("Role", td.role)
			req.Header.Add("UserID", td.userID)
			res := httptest.NewRecorder()
			api.router.ServeHTTP(res, req)
			assert.Equal(t, td.code, res.Code)
			jsonassert.New(t).Assertf(res.Body.String(), td.resp)
		})
	}
}

func TestGetFile(t *testing.T) {
	api := NewAPI(nil, memory.NewTaskStore())
	testdata := []struct {
		des       string
		role      string
		userID    string
		taskID    string
		messageID string
		code      int
		resp      string
	}{
		{
			des:       "invalid role",
			role:      "abc",
			userID:    "1",
			taskID:    "1",
			messageID: "1",
			code:      http.StatusUnprocessableEntity,
			resp:      `{"errors":[{"id":"<<PRESENCE>>","code":"422001","detail":"invalid header","status":422,"title":"unprocessable entity"}]}`,
		},
		{
			des:       "invalid user id when the role is customer",
			role:      "customer",
			userID:    "123456",
			taskID:    "1",
			messageID: "1",
			code:      http.StatusForbidden,
			resp:      `{"errors":[{"id":"<<PRESENCE>>","code":"403000","status":403,"title":"forbidden request"}]}`,
		},
	}

	for _, td := range testdata {
		t.Run(td.des, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+td.taskID+"/files/"+td.messageID, nil)
			req.Header.Add("Role", td.role)
			req.Header.Add("UserID", td.userID)
			res := httptest.NewRecorder()
			api.router.ServeHTTP(res, req)
			assert.Equal(t, td.code, res.Code)
			jsonassert.New(t).Assertf(res.Body.String(), td.resp)
		})
	}
}
