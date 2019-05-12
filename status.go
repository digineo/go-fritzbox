package fritzbox

import (
	"gopkg.in/xmlpath.v2"
)

// Status represents the status of the WAN connection
type Status struct {
	Uptime              int    `json:"uptime"`
	ConnectionStatus    string `json:"connection_status"`
	LastConnectionError string `json:"last_connection_error"`
}

var (
	requestStatusInfo = requestType{
		path:       "/igdupnp/control/WANIPConn1",
		soapAction: "urn:schemas-upnp-org:service:WANIPConnection:1#GetStatusInfo",
		soapBody:   `<u:GetStatusInfo xmlns:u="urn:schemas-upnp-org:service:WANIPConnection:1" />`,
	}

	statusUptimePath          = xmlpath.MustCompile("//GetStatusInfoResponse/NewUptime")
	statusConnectionStatus    = xmlpath.MustCompile("//GetStatusInfoResponse/NewConnectionStatus")
	statusLastConnectionError = xmlpath.MustCompile("//GetStatusInfoResponse/NewLastConnectionError")
)

// GetStatusInfo liest den Verbindungsstatus aus
func (c *Client) GetStatusInfo() (status Status, err error) {
	result, err := c.request(requestStatusInfo)
	if err != nil {
		return
	}
	status = parseStatus(result)
	return
}

func parseStatus(root *xmlpath.Node) (status Status) {
	status.Uptime = pathToInt(root, statusUptimePath)
	status.ConnectionStatus, _ = statusConnectionStatus.String(root)
	status.LastConnectionError, _ = statusLastConnectionError.String(root)

	return
}
