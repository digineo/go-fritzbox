package fritzbox

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"testing"

	"gopkg.in/xmlpath.v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testfile(name string) *xmlpath.Node {
	file, err := os.Open(fmt.Sprintf("testdata/%s.xml", name))
	if err != nil {
		panic(err)
	}

	root, _ := xmlpath.Parse(file)
	return root
}

func TestIsFault(t *testing.T) {
	assert.EqualError(t, checkFault(testfile("error")), "UPNP failed: Invalid Action")
}

func TestNonFault(t *testing.T) {
	assert.Nil(t, checkFault(testfile("address/valid")))
}

func checkFault(root *xmlpath.Node) error {
	if faultPath.Exists(root) {
		if s, ok := faultDescription.String(root); ok {
			return errors.New("UPNP failed: " + s)
		}
		return errors.New("UPNP failed with unknown error")
	}
	return nil
}

func TestClient(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(err)
	defer listener.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/igdupnp/control/WANCommonIFC1", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/linkspeed/valid.xml")
	})

	srv := &http.Server{Handler: mux}
	go func() {
		if e := srv.Serve(listener); !errors.Is(e, http.ErrServerClosed) {
			log.Printf("unexpected HTTP server shutdown: %v", e)
		}
	}()
	defer srv.Close()

	client := Client{
		Endpoint: listener.Addr().String(),
		Context:  context.Background(),
	}

	// should succeed
	linkSpeed, err := client.GetLinkspeed()
	assert.NoError(err)
	if assert.NotNil(linkSpeed) {
		assert.Equal(51392000, linkSpeed.DownlinkBitRate)
	}

	// should fail (no handler)
	_, err = client.GetPublicAddress()
	if assert.Error(err) {
		assert.Contains(err.Error(), "unexpected status code 404 for http://127")
	}
}

func TestClientInvalidAddress(t *testing.T) {
	assert := assert.New(t)
	client := NewClientFromIP(net.IP{127, 0, 0, 2})

	assert.Equal("127.0.0.2:49000", client.Endpoint)

	// should succeed
	_, err := client.GetLinkspeed()
	if assert.Error(err) {
		assert.Contains(err.Error(), "dial tcp 127.0.0.2:49000: connect: connection refused")
	}
}
