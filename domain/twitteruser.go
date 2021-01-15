package domain

import "github.com/jake-hansen/followrs/repositories/apis/twitter"

type TwitterUser struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`
}

type TwitterService interface {
	GetUser(username string) (*TwitterUser, error)
}

type TwitterRepository interface {
	GetUser(username string) (*twitter.User, error)
}
