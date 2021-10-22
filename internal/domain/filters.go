package domain

import (
	"strings"
)

// Filters contains some properties about how client wants to view the data.
// including the page size, current page, the order of data, etc.
type Filters struct {
	CurrentPage  int      // CurrentPage represents the current page client wants to see.
	PageSize     int      // PageSize represents the page size of each page.
	Sort         string   // Sort represent the property that data needs to be sorted by. E.g. if Sort == "id", the data is sorted by id.
	SortSafelist []string // SortSafelist field to hold the supported sort values.
}

// SortColumn checks that the client-provided Sort field matches one of the entries in our safelist
// and if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists).
func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

// SortDirection returns the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field.
func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

// Limit returns the size of each page.
func (f Filters) Limit() int {
	return f.PageSize
}

// Offset returns the distance between first data and current data.
func (f Filters) Offset() int {
	return (f.CurrentPage - 1) * f.PageSize
}
