package fritzbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStatus(t *testing.T) {
	assert := assert.New(t)
	status := parseStatus(testfile("status/response"))
	assert.EqualValues("Connected", status.ConnectionStatus)
	assert.EqualValues("ERROR_NONE", status.LastConnectionError)
	assert.EqualValues(365376, status.Uptime)
}
