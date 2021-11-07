package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetRequestHeaders(t *testing.T) {
	r := requestHeaders{}

	errHeaderOk := r.Set("Content-Type:application/json")
	errHeaderWrong := r.Set("WrongHeader")

	assert.Nil(t, errHeaderOk, "Content-Type header should be parsed properly")
	assert.Error(t, errHeaderWrong, "Wrong header should cause an error")
}
