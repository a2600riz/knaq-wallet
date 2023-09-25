package network

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	TypeText      = "text/plain"
	TypeJSON      = "application/json"
	TypeForm      = "application/x-www-form-urlencoded"
	TypeJwt       = "application/jwt"
	TypeMultipart = "multipart/form-data"
)

type Sender interface {
	Send() (f []byte, status int, err error)
	SendWithCA(certs TransportCerts) (f []byte, status int, err error)
}

type Request struct {
	URL                string
	Method             string
	Accept             string
	ContentType        string
	Body               []byte
	FormData           url.Values
	MultipartData      *bytes.Buffer
	Bearer             string
	Basic              string
	InsecureSkipVerify bool
	CustomerHeader     map[string]string
}

func (r *Request) Send() (f []byte, status int, err error) {
	client := &http.Client{}
	if r.InsecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: r.InsecureSkipVerify},
		}
		client = &http.Client{Transport: tr}
	}

	var req *http.Request

	if r.ContentType == TypeForm {
		req, err = http.NewRequest(r.Method, r.URL, strings.NewReader(r.FormData.Encode()))
	} else if strings.Contains(r.ContentType, TypeMultipart) {
		req, err = http.NewRequest(r.Method, r.URL, r.MultipartData)
	} else {
		req, err = http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Body))
	}
	if err != nil {
		return nil, 0, err
	}

	if r.CustomerHeader != nil {
		keys := reflect.ValueOf(r.CustomerHeader).MapKeys()
		for _, k := range keys {
			req.Header.Set(k.String(), r.CustomerHeader[k.String()])
		}
	}
	if r.Accept != "" {
		req.Header.Add("Accept", r.Accept)
	}
	if r.ContentType != "" {
		req.Header.Add("Content-Type", r.ContentType)
	}
	if r.Bearer != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.Bearer))
	}
	if r.Basic != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Basic %s", r.Basic))
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	f, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, 0, err
	}

	return f, resp.StatusCode, err
}
func (r *Request) SendWithCA(certs TransportCerts) (f []byte, status int, err error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{certs.Cert},
				RootCAs:            certs.CertPool,
				InsecureSkipVerify: r.InsecureSkipVerify,
			},
		},
	}

	var req *http.Request

	if r.ContentType == TypeForm {
		req, err = http.NewRequest(r.Method, r.URL, strings.NewReader(r.FormData.Encode()))
	} else {
		req, err = http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Body))
	}
	if err != nil {
		return nil, 0, err
	}

	if r.CustomerHeader != nil {
		keys := reflect.ValueOf(r.CustomerHeader).MapKeys()
		for _, k := range keys {
			req.Header.Set(k.String(), r.CustomerHeader[k.String()])
		}
	}
	if r.Accept != "" {
		req.Header.Set("Accept", r.Accept)
	}
	if r.ContentType != "" {
		req.Header.Set("Content-Type", r.ContentType)
	}
	if r.Bearer != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.Bearer))
	}
	if r.Basic != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Basic %s", r.Basic))
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	f, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, 0, err
	}

	return f, resp.StatusCode, err
}
