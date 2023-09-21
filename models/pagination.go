package models

type Pagination struct {
	Limit  int
	Order  string
	Before *string
	After  *string
}
