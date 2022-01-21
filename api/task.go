package api

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/nurislam03/agent/model"
	"io"
	"net/http"
	"strings"
	"time"
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

	//task id that pssed throught he request param and the task id that is associated in the message file we want to download is the same. if not through an error
	if tskID != msg.TaskID {
		handleAPIError(w, newAPIError("forbidden request", errForbiddenRequest, err))
		return
	}

	// url of the the file we want to download (resource link of s3 bucket)
	url := msg.FileRefID

	//url := "https://upload.wikimedia.org/wikipedia/commons/f/f9/Flag_of_Bangladesh.svg"

	//url := "https://docs.google.com/viewer?url=https://www.computer-pdf.com/pdf/0669-introduction-to-calculus-volume-1.pdf"

	//url := "https://qtxasset.com/styles/breakpoint_xl_880px_w/s3/fiercebiotech/1638804283/BD.jpg/BD.jpg?VersionId=KG.fYNvmT7iQrKh_jdqPSGFP4a7AzDgl&itok=j1CYMpXC"

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "could not write response", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	//logger.Info(resp.Header.Get("Content-Type"))
	//logger.Info(resp.Header.Get("Content-Disposition"))

	// Setting header for the response file
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "file"))
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "could not read body", http.StatusInternalServerError)
		return
	}
}
