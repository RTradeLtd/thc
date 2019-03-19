package thc

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

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
