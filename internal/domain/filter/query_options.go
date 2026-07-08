package filter

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

func (qo *QueryOptions) WithPagination(limit int, offset int) *QueryOptions {
	qo.Pagination = &PaginationOptions{Limit: limit, Offset: offset}
	return qo
}

func (qo *QueryOptions) WithSorting(column string, asc bool) *QueryOptions {
	qo.Sorting = &SortingOptions{Column: column, Asc: asc}
	return qo
}

func (qo *QueryOptions) HasPagination() bool {
	return qo.Pagination != nil
}

func (qo *QueryOptions) HasSorting() bool {
	return qo.Sorting != nil
}
