package handler

type CableModem struct {
	Mac                string         `json:"mac"`
	CpeMac             *string        `json:"cpeMac,omitempty"`
	MacDomain          *string        `json:"macDomain,omitempty"`
	CableModemIndex    *int32         `json:"cableModemIndex,omitempty"`
	ConfigFile         *string        `json:"configFile,omitempty"`
	Model              *string        `json:"model,omitempty"`
	FiberNode          *string        `json:"fiberNode,omitempty"`
	Ipv4               *string        `json:"ipv4,omitempty"`
	Ipv6               *string        `json:"ipv6,omitempty"`
	CpeIpv4            *string        `json:"cpeIpv4,omitempty"`
	Transponder        *string        `json:"transponder,omitempty"`
	DocsisVersion      *DocsisVersion `json:"docsisVersion,omitempty"`
	Ppod               *string        `json:"ppod,omitempty"`
	Fqdn               *string        `json:"fqdn,omitempty"`
	State              *State         `json:"state,omitempty"`
	NotFoundDate       *string        `json:"notFoundDate,omitempty"`
	RegState           *int32         `json:"regState,omitempty"`
	FnName             *string        `json:"fnName,omitempty"`
	NumberOfGenerators *int32         `json:"numberOfGenerators,omitempty"`
	RpdName            *string        `json:"rpdName,omitempty"`
	UpdatedAt          *string        `json:"updatedAt,omitempty"`
	Bootr              *string        `json:"bootr,omitempty"`
	Vendor             *string        `json:"vendor,omitempty"`
	SwRev              *string        `json:"swRev,omitempty"`
	OltName            *string        `json:"oltName,omitempty"`
	PonName            *string        `json:"ponName,omitempty"`
	UpdatedAtTs        *int32         `json:"updatedAtTs,omitempty"`
	IsCpe              *bool          `json:"isCPE,omitempty"`
	CmtsType           *string        `json:"cmtsType,omitempty"`
	// This attribute represents the current type of device. metroe(1)
	DeviceType *int32 `json:"deviceType,omitempty"`
}

type DocsisVersion string

type State string
