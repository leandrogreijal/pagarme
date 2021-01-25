package transactions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalError(t *testing.T) {
	err := InternalError{Path: "/test"}
	assertTest := assert.New(t)
	assertTest.Equal("Mundipagg internal error. Path: /test", err.Error())
}

func TestInvalidValueError(t *testing.T) {
	err := InvalidValueError{"param", "value"}
	assertTest := assert.New(t)
	assertTest.Equal("param is invalid. Value: value", err.Error())
}
