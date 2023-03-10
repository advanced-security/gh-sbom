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

type Relationship struct {
	Element string `json:"spdxElementId"`
	Type    string `json:"relationshipType"`
	Related string `json:"relatedSpdxElement"`
}

type Doc struct {
	SPDXVersion       string `json:"spdxVersion"`
	DataLicense       string `json:"dataLicense"`
	SPDXID            string
	Name              string         `json:"name"`
	DocumentNamespace string         `json:"documentNamespace"`
	CreationInfo      CreationInfo   `json:"creationInfo"`
	Packages          []Package      `json:"packages"`
	Relationships     []Relationship `json:"relationships"`
	DocumentDescribes []string       `json:"documentDescribes"`
}

func MakeDoc(toolVersion, license, host, owner, name string, packages []Package) Doc {
	// https://spdx.github.io/spdx-spec/v2.3/
	docName := fmt.Sprintf("%s/%s/%s", host, owner, name)

	mainPackage := Package{
		Name:        name,
		SPDXID:      "SPDXRef-mainPackage",
		VersionInfo: "",
		DownloadLocation: fmt.Sprintf(
			"git+https://%s/%s/%s.git",
			host, owner, name,
		),
		FilesAnalyzed: false,
		ExternalRefs: []ExternalRef{
			{
				ReferenceCategory: "PACKAGE-MANAGER",
				ReferenceType:     "purl",
				ReferenceLocator:  fmt.Sprintf("pkg:github/%s/%s", owner, name),
			},
		},
		LicenseConcluded: "NOASSERTION",
		LicenseDeclared:  license,
		Supplier:         "NOASSERTION",
	}

	doc := Doc{
		SPDXVersion:       "SPDX-2.3",
		DataLicense:       "CC0-1.0",
		SPDXID:            "SPDXRef-DOCUMENT",
		Name:              docName,
		DocumentNamespace: "https://spdx.org/spdxdocs/" + docName + "-" + uuid.New().String(),
		CreationInfo: CreationInfo{
			Creators: []string{"Organization: GitHub, Inc", "Tool: gh-sbom-" + toolVersion},
			Created:  time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		},
		Relationships:     []Relationship{},
		Packages:          append([]Package{mainPackage}, packages...),
		DocumentDescribes: []string{mainPackage.SPDXID},
	}

	for _, p := range doc.Packages {
		doc.Relationships = append(doc.Relationships, Relationship{
			Element: mainPackage.SPDXID,
			Type:    "DEPENDS_ON",
			Related: p.SPDXID,
		})
	}

	return doc
}
