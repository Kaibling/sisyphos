package services

type filter interface {
	Filter(query string) ([]interface{}, error)
}
