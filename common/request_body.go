package common

import "fmt"

type RawRequestBody []byte

func (b *RawRequestBody) Set(value string) error {
	(*b) = []byte(value)
	return nil
}

func (b *RawRequestBody) String() string {
	return fmt.Sprint(*b)
}

func (b *RawRequestBody) IsCumulative() bool {
	return false
}
