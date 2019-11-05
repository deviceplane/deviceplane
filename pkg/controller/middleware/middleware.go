package middleware

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/deviceplane/deviceplane/pkg/utils"
)

const (
	TotalPagesHeader = "Total-Pages"
)

const MaxPageSize = int(100)

const (
	PageNumParam  = "page"
	PageSizeParam = "page_size"
	OrderByParam  = "order_by"
	OrderParam    = "order"
)

var (
	ErrPageNotFound = errors.New("the requested page was not found")

	ErrInvalidPageNumParameter  = errors.New(fmt.Sprintf("invalid %s parameter", PageNumParam))
	ErrInvalidPageSizeParameter = errors.New(fmt.Sprintf("invalid %s parameter", PageSizeParam))
	ErrInvalidOrderByParameter  = errors.New(fmt.Sprintf("invalid %s parameter", OrderByParam))
	ErrInvalidOrderParameter    = errors.New(fmt.Sprintf("invalid %s parameter", OrderParam))
)

type orderDirection string

const (
	OrderAscending  = orderDirection("asc")
	OrderDescending = orderDirection("desc")
)

func SortAndPaginateAndRespond(r http.Request, w http.ResponseWriter, arr []interface{}) {
	values := r.URL.Query()
	var pageNum *int
	if pageNumStr := values.Get(PageNumParam); pageNumStr != "" {
		p, err := strconv.Atoi(pageNumStr)
		if err != nil {
			http.Error(w, ErrInvalidPageNumParameter.Error(), http.StatusBadRequest)
			return
		}
		pageNum = &p
	}
	if pageNum == nil {
		p := 0
		pageNum = &p
	}

	var pageSize *int
	if pageSizeStr := values.Get(PageSizeParam); pageSizeStr != "" {
		p, err := strconv.Atoi(pageSizeStr)
		if err != nil || p <= 0 || p > MaxPageSize {
			http.Error(w, ErrInvalidPageSizeParameter.Error(), http.StatusBadRequest)
			return
		}
		pageSize = &p
	}
	if pageSize == nil {
		m := MaxPageSize
		pageSize = &m
	}

	var direction *orderDirection
	if orderStr := values.Get(OrderParam); orderStr != "" {
		switch orderStr {
		case string(OrderAscending):
			d := OrderAscending
			direction = &d
		case string(OrderDescending):
			d := OrderDescending
			direction = &d
		default:
			http.Error(w, ErrInvalidOrderParameter.Error(), http.StatusBadRequest)
			return
		}
	}
	if direction == nil {
		d := OrderAscending
		direction = &d
	}

	var orderBy string = values.Get(OrderByParam)
	if orderBy == "" {
		// Do nothing, no order
	} else {
		err := order(orderBy, *direction, arr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Set total pages header, as pages are required
	totalPages := int(math.Ceil(float64(len(arr)) / float64(*pageSize)))
	w.Header().Set(TotalPagesHeader, strconv.Itoa(totalPages))

	// If there are no items, and only initial page is requested, don't 404
	if len(arr) == 0 && *pageNum == 0 {
		utils.Respond(w, make([]interface{}, 0))
		return
	}

	var exists bool
	arr, exists = paginate(*pageNum, *pageSize, arr)
	if !exists {
		http.Error(w, ErrPageNotFound.Error(), http.StatusNotFound)
		return
	}

	utils.Respond(w, arr)
}
