package cyclonedx

import (
	"time"
)

type LicenseExpression struct {
	Expression string `json:"expression"`
}

type Component struct {
	Type     string              `json:"type"`
	Group    string              `json:"group,omitempty"`
	Name     string              `json:"name"`
	Version  string              `json:"version"`
	Purl     string              `json:"purl"`
	Licenses []LicenseExpression `json:"licenses,omitempty"`
}

type Tool struct {
	Vendor  string `json:"vendor"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type License struct {
	Expression string `json:"expression"`
}

type Metadata struct {
	Timestamp string    `json:"timestamp"`
	Tools     []Tool    `json:"tools"`
	Licenses  []License `json:"licenses"`
}

type Doc struct {
	BomFormat   string      `json:"bomFormat"`
	SpecVersion string      `json:"specVersion"`
	Version     int         `json:"version"`
	Metadata    Metadata    `json:"metadata"`
	Components  []Component `json:"components"`
}

func MakeDoc(toolVersion string, components []Component) Doc {
	// https://cyclonedx.org/docs/1.4/json/

	tool := Tool{
		Vendor:  "advanced-security",
		Name:    "gh-sbom",
		Version: toolVersion,
	}

	license := License{
		Expression: "CC0-1.0",
	}

	return Doc{
		BomFormat:   "CycloneDX",
		SpecVersion: "1.4",
		Version:     1,
		Metadata: Metadata{
			Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			Tools:     []Tool{tool},
			Licenses:  []License{license},
		},
		Components: components,
	}
}
