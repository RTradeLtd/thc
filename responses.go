package thc

// LoginResponse is a response from the login api call
type LoginResponse struct {
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

// Response is a general api response applicable for multiple calls
type Response struct {
	Code     int    `json:"code"`
	Response string `json:"response"`
}
