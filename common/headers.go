package common

import (
	"fmt"
	"strings"
)

// RequestHeader represents a single HTTP request header (a key: value pair)
type RequestHeader struct {
	Key   string
	Value string
}

// RequestHeaders is a slice of request headers that will be added to the request
type RequestHeaders []RequestHeader

func (r *RequestHeaders) Set(value string) error {
	tokenized := strings.Split(value, ":")
	if len(tokenized) != 2 {
		return fmt.Errorf("Header '%s' doesn't match 'Key:Value' format (i.e. 'Content-Type:application/json')", value)
	}
	h := RequestHeader{tokenized[0], tokenized[1]}
	(*r) = append(*r, h)
	return nil
}

func (r *RequestHeaders) String() string {
	return fmt.Sprint(*r)
}

func (r *RequestHeaders) IsCumulative() bool {
	return true
}
