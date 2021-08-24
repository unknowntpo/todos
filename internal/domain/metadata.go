package domain

type Metadata interface {
	Calculate(totalRecords, page, pageSize int)
}
