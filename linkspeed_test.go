package fritzbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkspeedValid(t *testing.T) {
	assert := assert.New(t)
	linkspeed := parseLinkspeed(testfile("linkspeed/valid"))
	assert.Equal("Up", linkspeed.Status)
	assert.Equal(10048000, linkspeed.UplinkBitRate)
	assert.Equal(51392000, linkspeed.DownlinkBitRate)
}

func TestLinkspeedInitializing(t *testing.T) {
	assert := assert.New(t)
	linkspeed := parseLinkspeed(testfile("linkspeed/initializing"))
	assert.Equal("Initializing", linkspeed.Status)
	assert.Equal(0, linkspeed.UplinkBitRate)
	assert.Equal(0, linkspeed.DownlinkBitRate)
}
