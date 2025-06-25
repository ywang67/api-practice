package helpers

import (
	"fmt"
	"strings"

	"api-project/grpc-api/gen/cablemodems"
)

//
// DocsisVersion
//

func ParseDocsisVersionFromString(s string) (cablemodems.DocsisVersion, error) {
	switch strings.ToLower(s) {
	case "docsis3":
		return cablemodems.DocsisVersion_DOCSIS3, nil
	case "docsis31":
		return cablemodems.DocsisVersion_DOCSIS31, nil
	case "docsis4":
		return cablemodems.DocsisVersion_DOCSIS4, nil
	default:
		return cablemodems.DocsisVersion_DOCSIS_UNKNOWN, fmt.Errorf("invalid docsis version: %s", s)
	}
}

func DocsisVersionToString(d cablemodems.DocsisVersion) string {
	switch d {
	case cablemodems.DocsisVersion_DOCSIS3:
		return "Docsis3"
	case cablemodems.DocsisVersion_DOCSIS31:
		return "Docsis31"
	case cablemodems.DocsisVersion_DOCSIS4:
		return "Docsis4"
	default:
		return "Unknown"
	}
}

//
// State
//

func ParseStateFromString(s string) (cablemodems.State, error) {
	switch strings.ToLower(s) {
	case "online":
		return cablemodems.State_ONLINE, nil
	case "offline":
		return cablemodems.State_OFFLINE, nil
	default:
		return cablemodems.State_UNKNOWN, fmt.Errorf("invalid state: %s", s)
	}
}

func StateToString(state cablemodems.State) string {
	switch state {
	case cablemodems.State_ONLINE:
		return "Online"
	case cablemodems.State_OFFLINE:
		return "Offline"
	default:
		return "Unknown"
	}
}
