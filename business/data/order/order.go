package order

import (
	"errors"
	"fmt"
	"net/http"
	"sales-api/foundation/validate"
	"strings"
)

// Set of directions for data ordering.
const (
	ASC  = "ASC"
	DESC = "DESC"
)

var directions = map[string]string{
	ASC:  ASC,
	DESC: DESC,
}

// =============================================================================

// By represents a field used to order by and direction.
type By struct {
	Field     string
	Direction string
}

func NewBy(feild string, direction string) By {
	return By{
		Field: feild, Direction: direction,
	}
}

// =============================================================================

// Parse constructs a order.By value by parsing a string in the form
// of "field,direction".
func Parse(r *http.Request, defaultOrder By) (By, error) {
	v := r.URL.Query().Get("orderBy")

	if v == "" {
		return defaultOrder, nil
	}
	orderParts := strings.Split(v, ",")

	var by By
	switch len(orderParts) {
	case 1:
		by = NewBy(strings.TrimSpace(orderParts[0]), ASC)
	case 2:
		by = NewBy(strings.TrimSpace(orderParts[0]), strings.TrimSpace(orderParts[1]))
	default:
		return By{}, validate.NewFieldsError(v, errors.New("unknown order field"))
	}

	if _, exists := directions[by.Direction]; !exists {
		return By{}, validate.NewFieldsError(v, fmt.Errorf("unknown direction: %s", by.Direction))
	}

	return by, nil

}
