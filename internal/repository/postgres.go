package repository

import (
	"strings"
)

type PostgesDB struct {
	Localhost string
	Port      string
	User      string
	Password  string
	Name      string
}

func NewPostgesDB(dataBaseDSN string) *PostgesDB {
	dataBaseDSNArray := strings.Split(dataBaseDSN, ":")
	localhost, post, user, password, name := dataBaseDSNArray[0], dataBaseDSNArray[1], dataBaseDSNArray[2], dataBaseDSNArray[3], dataBaseDSNArray[4]
	return &PostgesDB{
		Localhost: localhost,
		Port:      post,
		User:      user,
		Password:  password,
		Name:      name,
	}
}
