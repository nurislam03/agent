package api

import (
	"github.com/go-chi/chi"
	"github.com/nurislam03/agent/model"
	"net/http"
	"strings"
)

// @route  GET api/v1/tasks
// @desc   getTasks return all tasks based on user role
func (a *API) getTasks(w http.ResponseWriter, r *http.Request) {
	uRole := r.Header.Get("Role")
	uID := ""
	//if the request is form the customer side then we need the user id to pull only those tasks that are created by that specific user.
	if uRole == string(model.Customer) {
		uID = r.Header.Get("UserID")
	}

	// pager
	pgr := newPager(r.URL, 2)

	tskList, err := a.task.GetTasks(uID, pgr.limit(), pgr.offset())
	if err != nil {
		handleAPIError(w, newAPIError("data store error", errInternalServer, err))
		return
	}

	resp := response{
		code: http.StatusOK,
		Data: tskList,
		Meta: pgr,
	}
	resp.serveJSON(w)
}

// @route  GET api/v1/tasks/{task_id}/chats
// @desc   getChatHistory return all chat messages of a specific task
func (a *API) getChatHistory(w http.ResponseWriter, r *http.Request) {
	uRole := r.Header.Get("Role")
	uID := r.Header.Get("UserID")

	//getting the task id form the url param.
	tskID := strings.TrimSpace(chi.URLParam(r, "task_id"))

	//if the request came from a customer side then need to check (before pulling the chat history) that customer is the owner of that task.
	if uRole == string(model.Customer) {
		tsk, err := a.task.GetTaskByID(tskID)
		if err != nil {
			handleAPIError(w, newAPIError("data store error", errInternalServer, err))
			return
		}

		// checking whether the customer id associated with the task and the user who made the request is the same one or not.
		if tsk.CustomerID != uID {
			handleAPIError(w, newAPIError("unauthorized request", errInternalServer, err)) //Todo: handle error message and err tag
			return
		}
	}

	// pager
	pgr := newPager(r.URL, 10)

	chats, err := a.task.GetChatHistory(tskID, pgr.limit(), pgr.offset())
	if err != nil {
		handleAPIError(w, newAPIError("data store error", errInternalServer, err))
		return
	}

	resp := response{
		code: http.StatusOK,
		Data: chats,
		Meta: pgr,
	}
	resp.serveJSON(w)
}
