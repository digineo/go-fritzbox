package fritzbox

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/xmlpath.v2"
)

type requestType struct {
	path       string
	soapAction string
	soapBody   string
}

var (
	timeout          = time.Duration(5) * time.Second
	faultPath        = xmlpath.MustCompile("/Envelope/Body/Fault")
	faultDescription = xmlpath.MustCompile("/Envelope/Body/Fault//errorDescription")
)

// DefaultPort is the default UPNP port.
const DefaultPort = 49000

// Client is a Fritzbox-Client.
type Client struct {
	Endpoint string
	Context  context.Context // optional context for cancelation of outgoing requests
}

// NewClientFromIP creates a new client by its IP address
func NewClientFromIP(ip net.IP) Client {
	return Client{
		Endpoint: net.JoinHostPort(ip.String(), strconv.Itoa(DefaultPort)),
	}
}

// Schickt einen UPNP-Request an die Fritz-Box
func (c *Client) request(rType requestType) (root *xmlpath.Node, err error) {
	uri := "http://" + c.Endpoint + rType.path
	req, err := http.NewRequest(
		"POST",
		uri,
		bytes.NewBufferString(`<?xml version='1.0' encoding='utf-8'?><s:Envelope s:encodingStyle='http://schemas.xmlsoap.org/soap/encoding/' xmlns:s='http://schemas.xmlsoap.org/soap/envelope/'><s:Body>`+rType.soapBody+`</s:Body></s:Envelope>`),
	)
	if err != nil {
		return
	}

	// Add headers
	req.Header.Add("SOAPAction", rType.soapAction)
	req.Header.Add("Content-Type", `text/xml; charset="utf-8"`)

	// Add context
	if ctx := c.Context; ctx != nil {
		req = req.WithContext(ctx)
	}

	// Send request
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		err = &ResponseError{
			URI:        uri,
			StatusCode: resp.StatusCode,
		}
		return
	}

	root, err = xmlpath.Parse(resp.Body)
	return
}

func pathToInt(node *xmlpath.Node, path *xmlpath.Path) int {
	if s, ok := path.String(node); ok {
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}

	return 0
}
