package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	gh "github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/repository"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/spf13/pflag"
)

type Query struct {
	Repository struct {
		DependencyGraphManifests struct {
			Nodes []struct {
				Filename     string
				Dependencies struct {
					Nodes []struct {
						PackageManager string
						PackageName    string
						Requirements   string
					}
					PageInfo struct {
						HasNextPage bool
						EndCursor   string
					}
				} `graphql:"dependencies(first: $first, after: $dependencyCursor)"`
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		} `graphql:"dependencyGraphManifests(first: $first, after: $manifestCursor)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

type DependencyMap map[string]map[string]map[string]string

func getDependencies(repoOwner, repoName string) DependencyMap {
	dependencies := make(DependencyMap)

	opts := api.ClientOptions{
		Headers: map[string]string{"Accept": "application/vnd.github.hawkgirl-preview+json"},
	}

	client, err := gh.GQLClient(&opts)
	if err != nil {
		log.Fatal(err)
	}

	var manifestCursor, dependencyCursor *string

	for {
		var query Query

		makeQuery(client, repoOwner, repoName, (*graphql.String)(manifestCursor), (*graphql.String)(dependencyCursor), &query, &dependencies)

		manifestCursor = &query.Repository.DependencyGraphManifests.PageInfo.EndCursor

		if !query.Repository.DependencyGraphManifests.PageInfo.HasNextPage {
			break
		}
	}

	return dependencies
}

func makeQuery(client api.GQLClient, repoOwner, repoName string, manifestCursor, dependencyCursor *graphql.String, query *Query, dependencies *DependencyMap) {
	variables := map[string]interface{}{
		"name":             graphql.String(repoName),
		"owner":            graphql.String(repoOwner),
		"first":            graphql.Int(100),
		"manifestCursor":   manifestCursor,
		"dependencyCursor": dependencyCursor,
	}

	err := client.Query("RepositoryDependencies", &query, variables)
	if err != nil {
		log.Fatal(err)
	}

	for _, eachManifestNode := range query.Repository.DependencyGraphManifests.Nodes {
		for _, eachDependencyNode := range eachManifestNode.Dependencies.Nodes {
			packageManager := strings.ToLower(eachDependencyNode.PackageManager)
			packageName := strings.ToLower(eachDependencyNode.PackageName)

			if _, ok := (*dependencies)[eachManifestNode.Filename]; !ok {
				(*dependencies)[eachManifestNode.Filename] = make(map[string]map[string]string)
			}
			if _, ok := (*dependencies)[eachManifestNode.Filename][packageManager]; !ok {
				(*dependencies)[eachManifestNode.Filename][packageManager] = make(map[string]string)
			}
			(*dependencies)[eachManifestNode.Filename][packageManager][packageName] = eachDependencyNode.Requirements
		}

		dependencyCursor = (*graphql.String)(&eachManifestNode.Dependencies.PageInfo.EndCursor)

		if eachManifestNode.Dependencies.PageInfo.HasNextPage {
			var newQuery Query
			makeQuery(client, repoOwner, repoName, manifestCursor, dependencyCursor, &newQuery, dependencies)
		}
	}
}

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

func getLicense(packageManager, packageName, version string) (string, string, error) {
	client := &http.Client{}

	cdType := packageManager
	cdProvider := packageManager
	cdNamespace := "-"
	cdName := packageName
	cdRevision := version

	if packageManager == "pip" {
		cdType = "pypi"
		cdProvider = "pypi"
	} else if packageManager == "npm" {
		cdType = "npm"
		cdProvider = "npmjs"
		if strings.HasPrefix(packageName, "@") {
			packageParts := strings.SplitN(packageName, "/", 2)
			cdNamespace = packageParts[0]
			cdName = packageParts[1]
		}
	} else if packageManager == "go" {
		cdType = "go"
		cdProvider = "golang"
		packageParts := strings.SplitN(packageName, "/", 3)
		if len(packageParts) != 3 {
			return "", "", errors.New("Unable to parse go package " + packageName)
		}
		cdNamespace = packageParts[0] + "/" + packageParts[1]
		cdName = packageParts[2]
		cdRevision = "v" + version
	}

	// Useful for debugging ecosystems!
	//log.Printf("%s %s %s %s %s", cdType, cdProvider, cdNamespace, cdName, cdRevision)

	req, err := http.NewRequest("GET", "https://api.clearlydefined.io/definitions/"+url.PathEscape(cdType)+"/"+url.PathEscape(cdProvider)+"/"+url.PathEscape(cdNamespace)+"/"+url.PathEscape(cdName)+"/"+url.PathEscape(cdRevision), nil)
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

type Package struct {
	PackageName             string
	SPDXID                  string
	PackageVersion          string
	PackageDownloadLocation string
	FilesAnalyzed           bool
	PackageLicenseConcluded string
	PackageLicenseDeclared  string
}

type SPDXDoc struct {
	SPDXVersion  string
	DataLicense  string
	SPDXID       string
	DocumentName string
	Creator      string
	Created      string
	Packages     []Package
}

func main() {
	repoOverride := pflag.StringP("repository", "r", "", "Repository to query. Current directory used by default.")
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

	dependencies := getDependencies(repo.Owner(), repo.Name())
	packages := []Package{}

	i := 0

	for _, manifestMap := range dependencies {
		for packageManager, packageManagerMap := range manifestMap {
			for packageName, requirements := range packageManagerMap {
				pkg := Package{
					PackageName:             packageName,
					SPDXID:                  fmt.Sprintf("SPDXRef-%d", i),
					PackageVersion:          requirements[2:],
					PackageDownloadLocation: "NOASSERTION",
					FilesAnalyzed:           false,
					PackageLicenseDeclared:  "NOASSERTION",
					PackageLicenseConcluded: "NOASSERTION",
				}

				declared, discovered, err := getLicense(packageManager, packageName, requirements[2:])
				if err == nil && len(declared) > 0 {
					pkg.PackageLicenseDeclared = declared
				} else if err == nil && len(discovered) > 0 {
					pkg.PackageLicenseConcluded = discovered
				}
				packages = append(packages, pkg)
				i += 1
			}
		}
	}

	doc := SPDXDoc{
		SPDXVersion:  "SPDX-2.3",
		DataLicense:  "CC0-1.0",
		SPDXID:       "SPDXRef-DOCUMENT",
		DocumentName: fmt.Sprintf("%s/%s/%s", repo.Host(), repo.Owner(), repo.Name()),
		Creator:      "Tool https://github.com/steiza/gh-sbom",
		Created:      time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		Packages:     packages,
	}

	jsonBinary, err := json.Marshal(&doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonBinary))
}
