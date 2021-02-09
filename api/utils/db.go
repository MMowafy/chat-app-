package utils

import (
	"net/http"
	"strconv"
)

type ListRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Skip     int    `json:"skip"`
	OrderBy  string `json:"orderBy"`
	Order    string `json:"order"`
	Query    interface{}
}

func NewListRequest(r *http.Request) *ListRequest {
	var request ListRequest

	request.PageSize = 20
	request.Order = ""
	request.OrderBy = ""
	request.Page = 1
	request.FromRequest(r)
	return &request
}

func (listRequest *ListRequest) FromRequest(r *http.Request) *ListRequest {
	q := r.URL.Query()

	pageSize := q.Get("pageSize")
	page := q.Get("page")

	pageNumber, _ := strconv.Atoi(page)
	if pageNumber > 0 {
		listRequest.Page = pageNumber
	}

	pageSizeLimit, _ := strconv.Atoi(pageSize)
	if pageSizeLimit > 0 {
		listRequest.PageSize = pageSizeLimit
	}

	listRequest.Skip = (listRequest.Page - 1) * listRequest.PageSize

	return listRequest
}
