package helpers

import (
	"math"
	"strings"
)

// Filters is use by task layer: GetAll api.
type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string // SortSafelist field to hold the supported sort values.
}

// sortColumn Check that the client-provided Sort field matches one of the entries in our safelist
// and if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists).
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

// sortDirection Return the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// Metadata holds the pagination metadata.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// Calculate calculates the appropriate pagination metadata
// values given the total number of records, current page, and page size values. Note
// that the last page value is calculated using the math.Ceil() function, which rounds
// up a float to the nearest integer. So, for example, if there were 12 records in total
// and a page size of 5, the last page value would be math.Ceil(12/5) = 3.
func (m *Metadata) Calculate(totalRecords, page, pageSize int) {
	if totalRecords == 0 {
		// empty Metadata struct if there are no records.
		return
	}

	m.CurrentPage = page
	m.PageSize = pageSize
	m.FirstPage = 1
	m.LastPage = int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	m.TotalRecords = totalRecords
}