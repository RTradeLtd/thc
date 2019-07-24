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

// IndexHash is used to index a hash with lens
func (v2 *V2) IndexHash(hash string, reindex bool) (string, error) {
	bodyBuff := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuff)
	idWriter, err := bodyWriter.CreateFormField("object_identifier")
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(idWriter, strings.NewReader(hash)); err != nil {
		return "", err
	}
	typeWriter, err := bodyWriter.CreateFormField("object_type")
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(typeWriter, strings.NewReader("ipld")); err != nil {
		return "", err
	}
	if reindex {
		reindexWriter, err := bodyWriter.CreateFormField("reindex")
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(reindexWriter, strings.NewReader("true")); err != nil {
			return "", err
		}
	}
	if err := bodyWriter.Close(); err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		"POST",
		v2.formatURL(LensIndex),
		bodyBuff,
	)
	if err != nil {
		return "", err
	}
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
	response := &IndexResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		return "", err
	}
	return response.Hash, nil
}

func (v2 *V2) SearchLens(query string) (*SearchResponse, error) {
	bodyBuff := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuff)
	queryWriter, err := bodyWriter.CreateFormField("query")
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(queryWriter, strings.NewReader(query)); err != nil {
		return nil, err
	}
	if err := bodyWriter.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"POST",
		v2.formatURL(LensSearch),
		bodyBuff,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	res, err := v2.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, v2.handleError(body)
	}
	var response SearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
