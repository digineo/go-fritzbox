package fritzbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressValid(t *testing.T) {
	assert := assert.New(t)
	addr, err := parseExternalAddress(testfile("address/valid"))
	assert.Nil(err)
	assert.Equal("91.57.217.71", addr.String())
}

func TestAddressEmpty(t *testing.T) {
	assert := assert.New(t)
	addr, err := parseExternalAddress(testfile("address/zero"))
	assert.EqualError(err, ErrEmptyAddress.Error())
	assert.Nil(addr)
}

func TestAddressInvalid(t *testing.T) {
	assert := assert.New(t)
	addr, err := parseExternalAddress(testfile("address/invalid"))
	assert.EqualError(err, "Could not find NewExternalIPAddress in response")
	assert.Nil(addr)
}
