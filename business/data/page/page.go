package page

import (
	"net/http"
	"sales-api/foundation/validate"
	"strconv"
)

// Page represents the requested page and rows per page.
type Page struct {
	Page     int
	PageSize int
}

// Parse parses the request for the page and rows query string. The
// defaults are provided as well.
func Parse(r *http.Request) (Page, error) {
	values := r.URL.Query()

	number := 1
	if page := values.Get("page"); page != "" {
		var err error
		number, err = strconv.Atoi(page)
		if err != nil {
			return Page{}, validate.NewFieldsError("page", err)
		}
	}
	pageSize := 10
	if rows := values.Get("page_size"); rows != "" {
		var err error
		pageSize, err = strconv.Atoi(rows)
		if err != nil {
			return Page{}, validate.NewFieldsError("page_size", err)
		}
	}
	return Page{
		Page:     number,
		PageSize: pageSize,
	}, nil
}
