package fritzbox

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"gopkg.in/xmlpath.v2"

	"github.com/stretchr/testify/assert"
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
