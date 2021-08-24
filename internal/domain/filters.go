package domain

type Filters interface {
	sortColumn() string
	sortDirection() string
	limit() int
	offset() int
}
