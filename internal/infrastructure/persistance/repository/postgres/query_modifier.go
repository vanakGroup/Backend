package postgres

type PaginationOptions struct {
	Limit  int
	Offset int
}

type SortingOptions struct {
	Column string
	Asc    bool
}

type QueryOptions struct {
	Pagination *PaginationOptions
	Sorting    *SortingOptions
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{}
}

func (qp *QueryOptions) WithPagination(limit int, offset int) *QueryOptions {
	qp.Pagination = &PaginationOptions{
		Offset: offset,
		Limit:  limit,
	}

	return qp
}

func (qp *QueryOptions) WithSorting(column string, asc bool) *QueryOptions {
	qp.Sorting = &SortingOptions{
		Column: column,
		Asc:    asc,
	}

	return qp
}

func (q *QueryOptions) HasPagination() bool {
	return q.Pagination != nil
}

func (q *QueryOptions) HasSorting() bool {
	return q.Sorting != nil
}
