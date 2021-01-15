package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

// RequestFunc is a function that is executed on a request and modifies
// something within that request.
type RequestFunc func(request *retryablehttp.Request)

// ResponseFunc is similar to RequestFunc, but instead is executed on a
// response.
type ResponseFunc func(response *http.Response)

// API represents an API to consume.
type API struct {
	// BaseURL is the URL that should be prepended to all endpoint requests.
	BaseURL *url.URL

	// Auth is the authentication mechanism that is used to authenticate to the API.
	Auth *Auth

	// Client is the retryablehttp.Client that is used for endpoint requests.
	Client *retryablehttp.Client

	// BeforeRequest allows a user-supplied function to be called before each request.
	BeforeRequest RequestFunc

	// AfterResponse allows a user-supplied function to be called after each response.
	AfterResponse ResponseFunc
}

// Auth contains the functions needed to authenticate to a consumable API.
type Auth interface {
	// IsAuthenticated determines if we are authenticated to the API at the current moment.
	IsAuthenticated() bool

	// Attach attaches some sort of authentication credentials to a request. It is up to
	// the developer to determine how often the credentials need to be attached, as well
	// as what attaching credentials does in the context of the implementing API.
	// For example, an API that requires a Bearer JWT to be sent upon each request might
	// write an attach function that looks like this:
	//
	//		func (a *myAuthStruct) Attach(req *retryablehttp.Request) {
	//			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.APIBearerToken))
	//		}
	//
	// It is then up to the developer to execute the Attach function on each request.
	Attach(req *retryablehttp.Request)
}

// NewAPI creates a new API to be consumed.
func NewAPI(baseURL string, auth Auth, requestFunc RequestFunc, responseFunc ResponseFunc) (*API, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("could not create new api: %w", err)
	}

	defaultClient := retryablehttp.NewClient()

	api := &API{
		BaseURL:       base,
		Auth:          &auth,
		Client:        defaultClient,
		BeforeRequest: requestFunc,
		AfterResponse: responseFunc,
	}
	return api, nil
}

// Do performs the given request and stores the response in the given body.
// Note that before each request, Do prepends the requested endpoint with
// the API's baseURL. For example, if the request contains the URL
//
//					/v1/users
//
// then the actual request that is sent will be for the endpoint
//
//					{baseURL}/v1/users
//
// (be careful about trailing /'s)
//
func (api *API) Do(request *retryablehttp.Request, body interface{}) (*http.Response, error) {
	if api.BeforeRequest != nil {
		api.BeforeRequest(request)
	}

	newURL, err := url.Parse(fmt.Sprintf("%s%s", api.BaseURL, request.URL))
	if err != nil {
		return nil, fmt.Errorf("could not concatenate base url and endpoint: %w", err)
	}
	request.URL = newURL

	response, err := api.Client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to %s returned code %d", request.URL, response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(body)

	if err != nil {
		return nil, fmt.Errorf("could not decode body: %w", err)
	}

	if api.AfterResponse != nil {
		api.AfterResponse(response)
	}

	return response, err
}
