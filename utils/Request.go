package utils

import (
	"net/http"
	"encoding/json"
	"time"
	"net/url"
	"fmt"
	"errors"
)
type JSONResponse interface {

}

type Response struct {
	Error error
	Message JSONResponse
}

var client = &http.Client{Timeout: 3 * time.Second}

type RequestOptions struct {
	Scheme string
	Host   string
	Path   string
	Q      map[string]string
	Body   interface{}
}



func Get(options RequestOptions) Response {
	resp := Response{};

	var u = &url.URL{}
	u.Scheme = options.Scheme
	u.Host = options.Host
	u.Path = options.Path

	if options.Scheme == "" {
		u.Scheme = "http"
	}
	if options.Host == "" {
		resp.Error = errors.New("Host must be provided")
		return resp
	}
	if options.Path == "" {
		u.Path = "/"
	}

	if options.Q != nil {
		q := u.Query()
		for key, value := range options.Q {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}

	fmt.Println(u.String())
	res, err := client.Get(u.String())
	if err != nil {
		resp.Error = err
		return resp
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&resp.Message)

	return resp
}