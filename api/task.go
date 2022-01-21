package api

import (
	"github.com/nurislam03/agent/model"
	"net/http"
)

// @route  GET api/v1/tasks
// @desc   getTasks return all tasks based on user role
func (a *API) getTasks(w http.ResponseWriter, r *http.Request) {
	uRole := r.Header.Get("Role")
	uID := ""
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