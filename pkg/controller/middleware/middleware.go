package middleware

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/deviceplane/deviceplane/pkg/utils"
)

const (
	TotalPagesHeader     = "Total-Pages"
	TotalItemCountHeader = "Total-Item-Count"
)

const MaxPageSize = int(100)
const DefaultOrderByParam = "id"

const (
	PageSizeParam = "page_size"
	AfterParam    = "after"
	OrderParam    = "order"
	OrderByParam  = "order_by"
)

var (
	ErrInvalidPageSizeParameter = fmt.Errorf("invalid %s parameter", PageSizeParam)
	ErrInvalidOrderByParameter  = fmt.Errorf("invalid %s parameter", OrderByParam)
	ErrInvalidOrderParameter    = fmt.Errorf("invalid %s parameter", OrderParam)
)

type orderDirection string

const (
	OrderAscending  = orderDirection("asc")
	OrderDescending = orderDirection("desc")
)

func SortAndPaginateAndRespond(r http.Request, w http.ResponseWriter, arr []interface{}) {
	values := r.URL.Query()
	after := values.Get(AfterParam)

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

	// Set total count header
	w.Header().Set(TotalItemCountHeader, strconv.Itoa(len(arr)))

	arr, err := paginateAfter(after, "id", *pageSize, arr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.Respond(w, arr)
}
