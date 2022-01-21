package api

import (
	"bytes"
	"errors"
	"github.com/go-chi/chi"
	"github.com/nurislam03/agent/model"
	"io"
	"net/http"
	"strings"
)

// @route  GET api/v1/tasks
// @desc   getTasks return all tasks based on user role
func (a *API) getTasks(w http.ResponseWriter, r *http.Request) {
	uRole := r.Header.Get("Role")
	if uRole != string(model.Customer) && uRole != string(model.Operator) {
		handleAPIError(w, newAPIError("unprocessable entity", errInvalidData, errors.New("invalid header")))
		return
	}

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

	if uRole != string(model.Customer) && uRole != string(model.Operator) {
		handleAPIError(w, newAPIError("unprocessable entity", errInvalidData, errors.New("invalid header")))
		return
	}

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
			handleAPIError(w, newAPIError("forbidden request", errForbiddenRequest, err))
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

// @route:  GET api/v1/tasks/{task_id}/files/{message_id}
// @desc:   getFile downloads a file associated with a message
func (a *API) getFile(w http.ResponseWriter, r *http.Request) {
	uRole := r.Header.Get("Role")
	uID := r.Header.Get("UserID")

	if uRole != string(model.Customer) && uRole != string(model.Operator) {
		handleAPIError(w, newAPIError("unprocessable entity", errInvalidData, errors.New("invalid header")))
		return
	}

	tskID := strings.TrimSpace(chi.URLParam(r, "task_id"))

	//if the request came from a customer side then need to check (before pulling the chat message) that customer is the owner of that task.
	if uRole == string(model.Customer) {
		tsk, err := a.task.GetTaskByID(tskID)
		if err != nil {
			handleAPIError(w, newAPIError("data store error", errInternalServer, err))
			return
		}

		// checking whether the customer id associated with the task and the user who made the request is the same one or not.
		if tsk.CustomerID != uID {
			handleAPIError(w, newAPIError("forbidden request", errForbiddenRequest, err))
			return
		}
	}

	//file id could be used directly to the url params instead of message id, but I did not do that to protect the file id/url (id that linked to s3 bucket) from being exposed to outside
	msgID := strings.TrimSpace(chi.URLParam(r, "message_id"))

	msg, err := a.task.GetMessageByID(msgID)
	if err != nil {
		handleAPIError(w, newAPIError("data store error", errInternalServer, err))
		return
	}

	//task id that passed throughout he requests param and the task id that is associated in the message file we want to download is the same. if not through an error
	if tskID != msg.TaskID {
		handleAPIError(w, newAPIError("forbidden request", errForbiddenRequest, err))
		return
	}

	// url of the file we want to download (resource link of s3 bucket)
	obj, err := a.objectStore.GetObject(msg.FileRefID)
	if err != nil {
		handleAPIError(w, newAPIError("file store error", errInternalServer, err))
		return
	}

	// Setting header for the response file
	w.Header().Set("Content-Type", http.DetectContentType(obj))

	_, err = io.Copy(w, bytes.NewReader(obj))
	if err != nil {
		http.Error(w, "could not read body", http.StatusInternalServerError)
		return
	}
}
