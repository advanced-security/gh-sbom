package spdx

import (
	"fmt"
	"time"
)

type Package struct {
	PackageName             string
	SPDXID                  string
	PackageVersion          string
	PackageDownloadLocation string
	FilesAnalyzed           bool
	ExternalRef             string
	PackageLicenseConcluded string
	PackageLicenseDeclared  string
}

type Doc struct {
	SPDXVersion  string
	DataLicense  string
	SPDXID       string
	DocumentName string
	Creator      string
	Created      string
	Packages     []Package
}

func MakeDoc(host, owner, name string, packages []Package) Doc {
	// https://spdx.github.io/spdx-spec/v2.3/
	return Doc{
		SPDXVersion:  "SPDX-2.3",
		DataLicense:  "CC0-1.0",
		SPDXID:       "SPDXRef-DOCUMENT",
		DocumentName: fmt.Sprintf("%s/%s/%s", host, owner, name),
		Creator:      "Tool https://github.com/steiza/gh-sbom",
		Created:      time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		Packages:     packages,
	}
}
