package services

import (
	"sisyphos/lib/pagination"
	"sisyphos/models"
)

type paginationRepo interface {
	ReadPaginationIds(sqlQuery string) ([]pagination.DBResult, error)
}

type PaginationService struct {
	repo   paginationRepo
	dbInfo pagination.DatabaseInfo
}

func NewPaginationService(repo paginationRepo) *PaginationService {
	dbInfo := pagination.DatabaseInfo{}
	return &PaginationService{repo: repo, dbInfo: dbInfo}
}

func (s *PaginationService) Paginate(cursorString string, md models.MetaData, table string) ([]string, *pagination.Cursor, error) {
	ci := pagination.ParseCursor(cursorString)
	sortInfo := pagination.SortInfo{
		Field:      md.SortField,
		Order:      md.Order,
		Limit:      md.Limit,
		Table:      table,
		CursorInfo: ci,
	}
	return pagination.New(sortInfo, s.repo, s.dbInfo)
}
