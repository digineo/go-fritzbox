package fritzbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseError(t *testing.T) {
	assert.EqualError(t, &ResponseError{
		URI:        "http://",
		StatusCode: 500,
	}, "Unexpected status code 500 for http://")
}
