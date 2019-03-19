package thc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Call is a typed string representing an API call
type Call string

// String returns a string type of the call
func (c Call) String() string {
	return string(c)
}

// FillParams is used to fill parameter placements for api calls
func (c Call) FillParams(args ...interface{}) Call {
	return Call(
		fmt.Sprintf(
			c.String(), args...,
		),
	)
}

const (
	// DevURL is the api url for development temporal environment
	DevURL = "https://dev.api.temporal.cloud/v2"
	// ProdURL is the api url for production temporal environment
	ProdURL = "https://api.temporal.cloud/v2"
	// Login is the login api call
	Login Call = "/auth/login"
	// FileAddPublic is a file upload api call for public ipfs network
	FileAddPublic Call = "/ipfs/public/file/add"
	// PinAddPublic is a pin add api call for public ipfs network
	PinAddPublic Call = "/ipfs/public/pin/%s"
	// LensIndex is used to index content against lens
	LensIndex Call = "/lens/index"
)

// V2 is our interface with temporal's v2 api
type V2 struct {
	c   *http.Client
	url string
	auth
}

// auth is used to handle authentication with the api
type auth struct {
	user, pass, jwt string
}

// NewV2 instantiates our V2 client
func NewV2(user, pass, url string) *V2 {
	v2 := &V2{
		auth: auth{
			user: user,
			pass: pass,
		},
		c:   &http.Client{},
		url: url,
	}
	return v2
}

// formatURL is used to format the api call url
func (v2 *V2) formatURL(call Call) string {
	return v2.url + call.String()
}

// addAuthHeader is used to update the http request
// with the jwt used for authentication
func (v2 *V2) addAuthHeader(req *http.Request) {
	bearer := fmt.Sprintf("Bearer %s", v2.auth.jwt)
	req.Header.Add("Authorization", bearer)
}

// handleError is used to return a formatted error
// indicating the code that was given with the failure
// and the reason for the failure.
func (v2 *V2) handleError(body []byte) error {
	out := &Response{}
	if err := json.Unmarshal(body, out); err != nil {
		return err
	}
	return fmt.Errorf(
		"http status code %v returned with error %v",
		out.Code, out.Response,
	)
}
