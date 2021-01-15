package services

import (
	"fmt"
	"github.com/jake-hansen/followrs/domain"
)

type TwitterService struct {
	Repo *domain.TwitterRepository
}

func NewTwitterService(repo *domain.TwitterRepository) domain.TwitterService{
	service := &TwitterService{
		Repo: repo,
	}
	return service
}

func (t *TwitterService) GetUser(username string) (*domain.TwitterUser, error) {
	user, err := (*t.Repo).GetUser(username)
	if err != nil {
		return nil, fmt.Errorf("an error ocurred retreiving the user %s from Twitter: %w", username, err)
	}

	domainUser := &domain.TwitterUser{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
	}

	return domainUser, nil
}
