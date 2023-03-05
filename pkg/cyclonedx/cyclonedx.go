package cyclonedx

type LicenseExpression struct {
	Expression string `json:expression`
}

type Component struct {
	Type     string              `json:"type"`
	Group    string              `json:"group,omitempty"`
	Name     string              `json:"name"`
	Version  string              `json:"version"`
	Purl     string              `json:"purl"`
	Licenses []LicenseExpression `json:"licenses,omitempty"`
}

type Doc struct {
	BomFormat    string      `json:"bomFormat"`
	SpecVersion  string      `json:"specVersion"`
	SerialNumber string      `json:"serialNumber"`
	Version      int         `json:"version"`
	Components   []Component `json:"components"`
}

func MakeDoc(components []Component) Doc {
	// https://cyclonedx.org/docs/1.4/json/
	return Doc{
		BomFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Version:     1,
		Components:  components,
	}
}
