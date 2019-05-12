package fritzbox

import (
	xmlpath "gopkg.in/xmlpath.v2"
)

// Linkspeed represents the synchronized speed
type Linkspeed struct {
	Status          string `json:"status"`
	UplinkBitRate   int    `json:"uplink_bitrate"`
	DownlinkBitRate int    `json:"downlink_bitrate"`
}

var (
	requestLinkspeed = requestType{
		path:       "/igdupnp/control/WANCommonIFC1",
		soapAction: "urn:schemas-upnp-org:service:WANCommonInterfaceConfig:1#GetCommonLinkProperties",
		soapBody:   `<u:GetCommonLinkProperties xmlns:u="urn:schemas-upnp-org:service:WANCommonInterfaceConfig:1" />`,
	}

	uplinkPath   = xmlpath.MustCompile("//NewLayer1UpstreamMaxBitRate")
	downlinkPath = xmlpath.MustCompile("//NewLayer1DownstreamMaxBitRate")
	statusPath   = xmlpath.MustCompile("//NewPhysicalLinkStatus")
)

// GetLinkspeed returns the synchronized speed
func (c *Client) GetLinkspeed() (linkspeed Linkspeed, err error) {
	result, err := c.request(requestLinkspeed)
	if err == nil {
		linkspeed = parseLinkspeed(result)
	}
	return
}

func parseLinkspeed(root *xmlpath.Node) Linkspeed {
	status, _ := statusPath.String(root)

	return Linkspeed{
		Status:          status,
		UplinkBitRate:   pathToInt(root, uplinkPath),
		DownlinkBitRate: pathToInt(root, downlinkPath),
	}
}
