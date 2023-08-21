package models

type Group struct {
	DBInfo
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Allows      []string `json:"allows"`
	Users       []string `json:"users"`
}
