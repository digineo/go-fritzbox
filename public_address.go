package fritzbox

import (
	"errors"
	"net"

	"gopkg.in/xmlpath.v2"
)

var (
	requestPublicAddress = requestType{
		path:       "/igdupnp/control/WANIPConn1",
		soapAction: "urn:schemas-upnp-org:service:WANIPConnection:1#GetExternalIPAddress",
		soapBody:   `<u:GetExternalIPAddress xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1" />`,
	}
	externalAdressPath = xmlpath.MustCompile("//NewExternalIPAddress")

	// ErrEmptyAddress is returned if no public address is available
	ErrEmptyAddress = errors.New("received empty address")
)

// GetPublicAddress returns the public IP address
func (c *Client) GetPublicAddress() (address net.IP, err error) {
	result, err := c.request(requestPublicAddress)
	if err == nil {
		address, err = parseExternalAddress(result)
	}
	return
}

func parseExternalAddress(root *xmlpath.Node) (addr net.IP, err error) {
	if value, ok := externalAdressPath.String(root); ok {
		if value == "0.0.0.0" {
			// No public IP address available
			err = ErrEmptyAddress
		} else {
			// Die Adresse parsen
			addr = net.ParseIP(value)
		}
	} else {
		err = errors.New("Could not find NewExternalIPAddress in response")
	}

	return
}
