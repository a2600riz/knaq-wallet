package network

import (
	"bytes"
)

type Endpoint interface {
	GetEndpoint() string
	GetMethod() string
	Insecure() bool
	GetToken() string
	GetCustomHeader() map[string]string
}
type JSONData interface {
	GetJson() []byte
}
type MultipartData interface {
	GetMultipartFile() *bytes.Buffer
}
