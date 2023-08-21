package models

type Tag struct {
	DBInfo
	Name        string `json:"name"`
	Description string `json:"description"`
}
