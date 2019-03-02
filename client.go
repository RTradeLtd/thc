package thc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
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

// Login is used to authenticate with the API and generate a JWT
func (v2 *V2) Login() error {
	payload := fmt.Sprintf(
		"{\n  \"username\": \"%s\",\n  \"password\": \"%s\"\n}",
		v2.auth.user, v2.auth.pass,
	)
	req, err := http.NewRequest(
		"POST",
		v2.formatURL(Login),
		strings.NewReader(payload),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := v2.c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return v2.handleError(body)
	}
	out := &LoginResponse{}
	if err := json.Unmarshal(body, out); err != nil {
		return err
	}
	v2.auth.jwt = out.Token
	return nil
}

// FileAddOpts are options used to configure
// file uploads
type FileAddOpts struct {
	Encrypted  bool
	Passphrase string
	HoldTime   string
}

// FileAdd is used to add a file to ipfs
// it returns the hash of the file that was uploaded
func (v2 *V2) FileAdd(filePath string, opts FileAddOpts) (string, error) {
	fh, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fh.Close()
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", filePath)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(fileWriter, fh); err != nil {
		return "", err
	}
	holdWriter, err := bodyWriter.CreateFormField("hold_time")
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(holdWriter, strings.NewReader(opts.HoldTime)); err != nil {
		return "", err
	}
	if opts.Encrypted {
		passphraseWriter, err := bodyWriter.CreateFormField("passphrase")
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(passphraseWriter, strings.NewReader(opts.Passphrase)); err != nil {
			return "", err
		}
	}
	if err := bodyWriter.Close(); err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		"POST",
		v2.formatURL(FileAddPublic),
		bodyBuf,
	)
	if err != nil {
		return "", err
	}
	v2.addAuthHeader(req)
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	res, err := v2.c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", v2.handleError(body)
	}
	response := &Response{}
	if err := json.Unmarshal(body, response); err != nil {
		return "", err
	}
	return response.Response, nil
}

// PinAdd is used to pin an ipfs hash
func (v2 *V2) PinAdd(hash, holdTime string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	holdWriter, err := bodyWriter.CreateFormField("hold_time")
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(holdWriter, strings.NewReader(holdTime)); err != nil {
		return "", err
	}
	if err := bodyWriter.Close(); err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		"POST",
		v2.formatURL(PinAddPublic.FillParams(hash)),
		bodyBuf,
	)
	if err != nil {
		return "", err
	}
	v2.addAuthHeader(req)
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	res, err := v2.c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", v2.handleError(body)
	}
	response := &Response{}
	if err := json.Unmarshal(body, response); err != nil {
		return "", err
	}
	return response.Response, nil
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
