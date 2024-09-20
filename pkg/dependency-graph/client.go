package dg

import (
	"log"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
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

type DependencyMap map[string]map[string]map[string]bool

func GetDependencies(repoOwner, repoName string) DependencyMap {
	dependencies := make(DependencyMap)

	opts := api.ClientOptions{
		Headers: map[string]string{"Accept": "application/vnd.github.hawkgirl-preview+json"},
	}

	client, err := api.NewGraphQLClient(opts)
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

func makeQuery(client *api.GraphQLClient, repoOwner, repoName string, manifestCursor, dependencyCursor *graphql.String, query *Query, dependencies *DependencyMap) {
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
			packageVersion := ""
			if len(eachDependencyNode.Requirements) > 2 {
				packageVersion = eachDependencyNode.Requirements[2:]
			}
			if strings.HasPrefix(packageVersion, "v") {
				packageVersion = packageVersion[1:]
			}

			if _, ok := (*dependencies)[packageManager]; !ok {
				(*dependencies)[packageManager] = make(map[string]map[string]bool)
			}
			if _, ok := (*dependencies)[packageManager][packageName]; !ok {
				(*dependencies)[packageManager][packageName] = make(map[string]bool)
			}
			if _, ok := (*dependencies)[packageManager][packageName][packageVersion]; !ok {
				(*dependencies)[packageManager][packageName][packageVersion] = true
			}
		}

		dependencyCursor = (*graphql.String)(&eachManifestNode.Dependencies.PageInfo.EndCursor)

		if eachManifestNode.Dependencies.PageInfo.HasNextPage {
			var newQuery Query
			makeQuery(client, repoOwner, repoName, manifestCursor, dependencyCursor, &newQuery, dependencies)
		}
	}
}
