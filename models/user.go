package models

import "errors"

type User struct {
	Name     string   `json:"name"`
	Password string   `json:"password,omitempty"`
	Token    []string `json:"token"`
	Groups   []string `json:"groups"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name empty")
	}
	return nil
}

func UserArrayValidate(users []User) error {
	for _, u := range users {
		if err := u.Validate(); err != nil {
			return err
		}
	}
	return nil
}
