package twitter

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jake-hansen/followrs/repositories/apis"

	"github.com/hashicorp/go-retryablehttp"
)

// User represents a Twitter user.
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// UserService provides methods for accessing Twitter users via the API.
type UserService struct {
	baseURL            string
	twitterAPI         *apis.API
	userLookupEndpoint *Endpoint
}

// NewUserService creates a UserService with the default configuration.
func NewUserService(api *apis.API) *UserService {
	return &UserService{
		baseURL:            "/users",
		twitterAPI:         api,
		userLookupEndpoint: newUserLookupEndpoint(),
	}
}

func newUserLookupEndpoint() *Endpoint {
	return &Endpoint{
		URL: "/by/username/",
	}
}

func (u *UserService) parseError(wrapper *DataWrapper) error {
	if wrapper.Errors != nil {
		apiErrors := *wrapper.Errors
		if len(apiErrors) > 0 {
			firstError := apiErrors[0]
			if firstError.Title == "Not Found Error" {
				return errors.New("user not found")
			}
			return errors.New("unknown error from Twitter API")
		}
	}

	return nil
}

// Show returns the requested User.
func (u *UserService) Show(username string) (*User, error) {
	wrapper := &DataWrapper{
		Data:   new(User),
		Errors: new([]Error),
	}
	req, err := retryablehttp.NewRequest(http.MethodGet, fmt.Sprintf("%s%s%s", u.baseURL, u.userLookupEndpoint.URL, username), nil)
	if err != nil {
		return nil, err
	}

	err = u.userLookupEndpoint.PerformRequest(req, u.twitterAPI, wrapper)
	if err != nil {
		return nil, err
	}

	err = u.parseError(wrapper)
	if err != nil {
		return nil, err
	}

	return wrapper.Data.(*User), err
}
