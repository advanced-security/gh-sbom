package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	gh "github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/spf13/pflag"

	"github.com/advanced-security/gh-sbom/pkg/cyclonedx"
	dg "github.com/advanced-security/gh-sbom/pkg/dependency-graph"
	"github.com/advanced-security/gh-sbom/pkg/spdx"
)

type ClearlyDefinedDefinition struct {
	Licensed struct {
		Declared string
		Facets   struct {
			Core struct {
				Discovered struct {
					Expressions []string
				}
			}
		}
	}
}

type Purl struct {
	cdType     string
	cdProvider string
	Provider   string
	Namespace  string
	Name       string
	Version    string
}

func getPurl(packageManager, packageName, version string) Purl {
	p := Purl{}

	p.Namespace = ""
	p.Name = packageName
	p.Version = version

	if packageManager == "actions" {
		p.cdType = "git"
		p.cdProvider = "github"
		p.Provider = "github"
		packageParts := strings.SplitN(packageName, "/", 2)
		if len(packageParts) >= 2 {
			p.Namespace = packageParts[0]
			p.Name = packageParts[1]
		}
	} else if packageManager == "go" {
		p.cdType = "go"
		p.cdProvider = "golang"
		p.Provider = "golang"
		packageParts := strings.SplitN(packageName, "/", 3)
		if len(packageParts) == 2 {
			p.Namespace = packageParts[0]
			p.Name = packageParts[1]
		} else if len(packageParts) == 3 {
			p.Namespace = packageParts[0] + "/" + packageParts[1]
			p.Name = packageParts[2]
		}
	} else if packageManager == "rubygems" {
		p.cdType = "gem"
		p.cdProvider = "rubygems"
		p.Provider = "gem"
	} else if packageManager == "maven" {
		p.cdType = "maven"
		p.cdProvider = "mavenCentral"
		p.Provider = "maven"
		packageParts := strings.SplitN(packageName, ":", 2)
		if len(packageParts) >= 2 {
			p.Namespace = packageParts[0]
			p.Name = packageParts[1]
		}
	} else if packageManager == "npm" {
		p.cdType = "npm"
		p.cdProvider = "npmjs"
		p.Provider = "npm"
		if strings.HasPrefix(packageName, "@") {
			packageParts := strings.SplitN(packageName, "/", 2)
			if len(packageParts) >= 2 {
				p.Namespace = packageParts[0]
				p.Name = packageParts[1]
			}
		}
	} else if packageManager == "pip" {
		p.cdType = "pypi"
		p.cdProvider = "pypi"
		p.Provider = "pypi"
	}

	return p
}

func (p Purl) String() string {
	// https://github.com/package-url/purl-spec/blob/master/PURL-SPECIFICATION.rst
	prefix := "pkg:" + p.Provider + "/"
	if p.Namespace != "" {
		namespaceParts := strings.Split(p.Namespace, "/")
		for _, part := range namespaceParts {
			prefix = prefix + url.QueryEscape(part) + "/"
		}
	}

	suffix := ""
	if p.Version != "" {
		suffix = "@" + url.QueryEscape(p.Version)
	}

	return prefix + url.QueryEscape(p.Name) + suffix
}

func getLicense(p *Purl) (string, string, error) {
	client := &http.Client{}

	if p.Namespace == "" {
		p.Namespace = "-"
	}

	version := "v" + p.Version
	if p.cdType == "gem" || p.cdType == "pypi" || p.cdType == "maven" || p.cdType == "npm" {
		version = p.Version
	}

	req, err := http.NewRequest("GET", "https://api.clearlydefined.io/definitions/"+url.PathEscape(p.cdType)+"/"+url.PathEscape(p.cdProvider)+"/"+url.PathEscape(p.Namespace)+"/"+url.PathEscape(p.Name)+"/"+url.PathEscape(version), nil)
	if err != nil {
		log.Print(err)
		return "", "", err
	}

	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return "", "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return "", "", err
	}

	def := ClearlyDefinedDefinition{}
	err = json.Unmarshal(body, &def)
	if err != nil {
		log.Print(err)
		return "", "", err
	}

	declared := def.Licensed.Declared
	discovered := ""

	if len(def.Licensed.Facets.Core.Discovered.Expressions) > 0 {
		discovered = def.Licensed.Facets.Core.Discovered.Expressions[0]
	}

	return declared, discovered, nil
}

type RepoResp struct {
	License struct {
		SPDXId string `json:"spdx_id"`
	}
}

func getRepoLicense(owner, repo string) string {
	client, err := gh.RESTClient(nil)
	if err != nil {
		return ""
	}

	repoResp := RepoResp{}

	err = client.Get("repos/"+owner+"/"+repo, &repoResp)
	if err != nil {
		return ""
	}

	return repoResp.License.SPDXId
}

func main() {
	version := "0.0.8"

	repoOverride := pflag.StringP("repository", "r", "", "Repository to query. Current directory used by default.")
	cdx := pflag.BoolP("cyclonedx", "c", false, "Use CycloneDX SBOM format. Default is to use SPDX.")
	includeLicense := pflag.BoolP("license", "l", false, "Include license information from clearlydefined.io in SBOM.")
	pflag.Parse()

	var repo repository.Repository
	var err error

	if *repoOverride == "" {
		repo, err = gh.CurrentRepository()
	} else {
		repo, err = repository.Parse(*repoOverride)
	}

	if err != nil {
		log.Fatal(err)
	}

	dependencies := dg.GetDependencies(repo.Owner(), repo.Name())

	if len(dependencies) == 0 {
		log.Fatal("No dependencies found\n\nIf you own this repository, check if Dependency Graph is enabled:\nhttps://" + repo.Host() + "/" + repo.Owner() + "/" + repo.Name() + "/settings/security_analysis\n\n")
	}

	i := 0

	if *cdx {
		components := []cyclonedx.Component{}

		for packageManager, packageManagerMap := range dependencies {
			for packageName, requirementsMap := range packageManagerMap {
				for requirements, _ := range requirementsMap {
					p := getPurl(packageManager, packageName, requirements)

					c := cyclonedx.Component{
						Type:    "library",
						Name:    p.Name,
						Version: p.Version,
						Purl:    p.String(),
					}

					if p.Namespace != "" {
						c.Group = p.Namespace
					}

					if *includeLicense {
						l := []cyclonedx.LicenseExpression{}
						declared, discovered, err := getLicense(&p)
						if err == nil && len(declared) > 0 {
							le := cyclonedx.LicenseExpression{
								Expression: declared,
							}
							c.Licenses = append(l, le)
						} else if err == nil && len(discovered) > 0 {
							le := cyclonedx.LicenseExpression{
								Expression: discovered,
							}
							c.Licenses = append(l, le)
						}
					}

					components = append(components, c)
				}
			}
		}

		doc := cyclonedx.MakeDoc(version, components)
		jsonBinary, err := json.Marshal(&doc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonBinary))

	} else {
		packages := []spdx.Package{}

		for packageManager, packageManagerMap := range dependencies {
			for packageName, requirementsMap := range packageManagerMap {
				for requirements, _ := range requirementsMap {
					purl := getPurl(packageManager, packageName, requirements)

					externalRef := spdx.ExternalRef{
						ReferenceCategory: "PACKAGE-MANAGER",
						ReferenceType:     "purl",
						ReferenceLocator:  purl.String(),
					}

					pkg := spdx.Package{
						Name:             purl.Name,
						SPDXID:           fmt.Sprintf("SPDXRef-%d", i),
						VersionInfo:      purl.Version,
						DownloadLocation: "NOASSERTION",
						FilesAnalyzed:    false,
						ExternalRefs:     []spdx.ExternalRef{externalRef},
						LicenseDeclared:  "NOASSERTION",
						LicenseConcluded: "NOASSERTION",
						Supplier:         "NOASSERTION",
					}

					if *includeLicense {
						declared, discovered, err := getLicense(&purl)
						if err == nil && len(declared) > 0 {
							pkg.LicenseDeclared = declared
						} else if err == nil && len(discovered) > 0 {
							pkg.LicenseConcluded = discovered
						}
					}

					packages = append(packages, pkg)
					i += 1
				}
			}
		}

		license := getRepoLicense(repo.Owner(), repo.Name())
		if license == "" {
			license = "NOASSERTION"
		}

		doc := spdx.MakeDoc(version, license, repo.Host(), repo.Owner(), repo.Name(), packages)

		jsonBinary, err := json.Marshal(&doc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonBinary))
	}
}
