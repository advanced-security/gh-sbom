package spdx

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ExternalRef struct {
	ReferenceCategory string `json:"referenceCategory"`
	ReferenceType     string `json:"referenceType"`
	ReferenceLocator  string `json:"referenceLocator"`
}

type Package struct {
	Name             string `json:"name"`
	SPDXID           string
	VersionInfo      string        `json:"versionInfo"`
	DownloadLocation string        `json:"downloadLocation"`
	FilesAnalyzed    bool          `json:"filesAnalyzed"`
	ExternalRefs     []ExternalRef `json:"externalRefs"`
	LicenseConcluded string        `json:"licenseConcluded"`
	LicenseDeclared  string        `json:"licenseDeclared"`
	Supplier         string        `json:"supplier"`
}

type CreationInfo struct {
	Creators []string `json:"creators"`
	Created  string   `json:"created"`
}

type Doc struct {
	SPDXVersion       string `json:"spdxVersion"`
	DataLicense       string `json:"dataLicense"`
	SPDXID            string
	Name              string       `json:"name"`
	DocumentNamespace string       `json:"documentNamespace"`
	CreationInfo      CreationInfo `json:"creationInfo"`
	Packages          []Package    `json:"packages"`
}

func MakeDoc(host, owner, name string, packages []Package) Doc {
	// https://spdx.github.io/spdx-spec/v2.3/
	docName := fmt.Sprintf("%s/%s/%s", host, owner, name)

	return Doc{
		SPDXVersion:       "SPDX-2.3",
		DataLicense:       "CC0-1.0",
		SPDXID:            "SPDXRef-DOCUMENT",
		Name:              docName,
		DocumentNamespace: "https://spdx.org/spdxdocs/" + docName + "-" + uuid.New().String(),
		CreationInfo: CreationInfo{
			Creators: []string{"Organization: GitHub, Inc", "Tool: gh-sbom"},
			Created:  time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		},
		Packages: packages,
	}
}
