package twitter

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/jake-hansen/followrs/repositories/apis"
)

// twitterAuth contains the API keys needed to authenticate to the
// Twitter API.
type twitterAuth struct {
	APIKey         string
	APISecretKey   string
	APIBearerToken string
}

// Twitter API is authenticated with pre generated tokens that don't expire.
// So, the client will always be considered authenticated.
func (a *twitterAuth) IsAuthenticated() bool {
	return true
}

// Attach attaches the API Bearer Token to a request.
func (a *twitterAuth) Attach(req *retryablehttp.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.APIBearerToken))
}

// API provides the services needed to interact with the Twitter API.
type API struct {
	Client      *apis.API
	UserService *UserService
}

// DataWrapper wraps the JSON response from the Twitter API.
type DataWrapper struct {
	Data   interface{} `json:"data"`
	Errors *[]Error    `json:"errors"`
}

// Error represents a problem that the Twitter API can return upon a
// request that fails.
type Error struct {
	Detail       string `json:"detail"`
	Title        string `json:"title"`
	ResourceType string `json:"resource_type"`
	Parameter    string `json:"parameter"`
	Value        string `json:"value"`
	Type         string `json:"type"`
}

// Endpoint represents a Twitter API endpoint. An endpoint contains information
// about the rate limits for itself.
type Endpoint struct {
	URL            string
	RemainingCalls int64
	RateLimitReset time.Time
}

func (e *Endpoint) parseRateLimitInfo(response *http.Response) error {
	headerNotFoundError := func(headerName string) error {
		return fmt.Errorf("header %s not found in response", headerName)
	}

	headerParseError := func(headerName string, err error) error {
		return fmt.Errorf("could not parse rate limit header %s: %w", headerName, err)
	}

	rateLimitRemainingHeader := "x-rate-limit-remaining"
	rateLimitResetTimeHeader := "x-rate-limit-reset"

	rateLimitRemaining := response.Header.Get(rateLimitRemainingHeader)
	if rateLimitRemaining == "" {
		return headerNotFoundError(rateLimitRemainingHeader)
	}
	rateLimitResetTime := response.Header.Get(rateLimitResetTimeHeader)
	if rateLimitResetTime == "" {
		return headerNotFoundError(rateLimitResetTimeHeader)
	}

	remaining, err := strconv.ParseUint(rateLimitRemaining, 10, 64)
	if err != nil {
		return headerParseError(rateLimitRemainingHeader, err)
	}
	resetTime, err := strconv.ParseUint(rateLimitResetTime, 10, 64)
	if err != nil {
		return headerParseError(rateLimitResetTimeHeader, err)
	}

	e.RemainingCalls = int64(remaining)
	e.RateLimitReset = time.Unix(int64(resetTime), 0)
	return nil
}

// PerformRequest is a helper function that requests a Twitter API URL on behalf of an endpoint.
// This function checks the rate limit for the endpoint before requesting the given URL. Upon
// a successful request, the endpoint is updated with the newly returned API rate limit information.
func (e *Endpoint) PerformRequest(request *retryablehttp.Request, api *apis.API, body interface{}) error {
	if e.RemainingCalls == 0 && time.Now().Before(e.RateLimitReset) {
		return fmt.Errorf("could not perform request. rate limit reached for %s", request.URL)
	}

	response, err := api.Do(request, body)
	if err != nil {
		return err
	}

	err = e.parseRateLimitInfo(response)
	if err != nil {
		return err
	}

	return nil
}

// NewTwitterAPI creates an API configured to be used with Twitter.
func NewTwitterAPI(baseURL string, apiKey string, apiSecretKey string, apiBearerToken string) (*API, error) {
	var beforeFuncs []apis.RequestFunc

	auth := &twitterAuth{
		APIKey:         apiKey,
		APISecretKey:   apiSecretKey,
		APIBearerToken: apiBearerToken,
	}

	beforeFuncs = append(beforeFuncs, auth.Attach)

	combinedFuncs := func(req *retryablehttp.Request) {
		for _, reqFunc := range beforeFuncs {
			reqFunc(req)
		}
	}

	api, err := apis.NewAPI(baseURL, auth, combinedFuncs, nil)

	twitterAPI := &API{
		Client:      api,
		UserService: NewUserService(api),
	}

	if err != nil {
		return nil, fmt.Errorf("could not create Twitter API: %w", err)
	}

	return twitterAPI, nil
}

func (a *API) GetUser(username string) (*User, error) {
	user, err := a.UserService.Show(username)
	return user, err
}
