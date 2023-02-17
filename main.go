// https://github.com/jhutchings1/sbom-generator/blob/main/src/index.js
// https://github.com/spdx/tools-golang/blob/v0.4.0/examples/3-build/example_build.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
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

type Package struct {
	PackageName             string
	SPDXID                  string
	PackageVersion          string
	PackageDownloadLocation string
	FilesAnalyzed           bool
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
		for _, packageManagerMap := range manifestMap {
			for packageName, requirements := range packageManagerMap {
				pkg := Package{
					PackageName:             packageName,
					SPDXID:                  fmt.Sprintf("SPDXRef-%d", i),
					PackageVersion:          requirements[2:],
					PackageDownloadLocation: "NOASSERTION",
					FilesAnalyzed:           false,
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
