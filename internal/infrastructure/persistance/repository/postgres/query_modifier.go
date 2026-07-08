package postgres

import (
	"github.com/Vanak/Backend/internal/domain/filter"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func applyQueryOptions(query *gorm.DB, options *filter.QueryOptions) *gorm.DB {
	if options == nil {
		return query
	}

	if options.HasSorting() {
		if options.Sorting.Asc {
			query = query.Order(clause.OrderByColumn{
				Column: clause.Column{Name: options.Sorting.Column},
			})
		} else {
			query = query.Order(clause.OrderByColumn{
				Column: clause.Column{Name: options.Sorting.Column},
				Desc:   true,
			})
		}
	}

	if options.HasPagination() {
		query = query.Offset(options.Pagination.Offset).Limit(options.Pagination.Limit)
	}

	return query
}
