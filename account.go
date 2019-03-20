package thc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

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

// GetJWT is used to return the JWT you use to authorize this account
func (v2 *V2) GetJWT() (string, error) {
	if v2.auth.jwt == "" {
		return "", errors.New("you have not logged in yet")
	}
	return v2.auth.jwt, nil
}
