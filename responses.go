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

// IndexResponse is a response from a lens index call
type IndexResponse struct {
	Hash        string   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	DisplayName string   `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	MimeType    string   `protobuf:"bytes,3,opt,name=mime_type,json=mimeType,proto3" json:"mime_type,omitempty"`
	Category    string   `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	Tags        []string `protobuf:"bytes,5,rep,name=tags,proto3" json:"tags,omitempty"`
}
