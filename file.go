package thc

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

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
