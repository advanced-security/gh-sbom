# gh-sbom

This is a `gh` CLI extension that outputs JSON SBOMs (in SPDX or CycloneDX format) for your GitHub repository using information from [Dependency graph](https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/about-the-dependency-graph).

It can optionally include license information with `-l`. License information comes from [ClearlyDefined](https://clearlydefined.io/)'s API.

## Installation
This is an extension to the `gh` CLI. See [gh CLI installation](https://github.com/cli/cli#installation) if you haven't installed `gh` yet.

Once `gh` is installed, you can install this extension with:
```
$ gh ext install steiza/gh-sbom
```

## Examples

This [gh CLI extension](https://docs.github.com/en/github-cli/github-cli/using-github-cli-extensions) outputs a SPDX SBOM from your [Dependency graph](https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/about-the-dependency-graph) from the command line:

```
$ gh sbom -l -r steiza/dependabot-example | jq
{
  "SPDXVersion": "SPDX-2.3",
  "DataLicense": "CC0-1.0",
  "SPDXID": "SPDXRef-DOCUMENT",
  "DocumentName": "github.com/steiza/dependabot-example",
  "Creator": "Tool https://github.com/steiza/gh-sbom",
  "Created": "2023-03-06T00:17:59Z",
  "Packages": [
      {
      "PackageName": "flake8",
      "SPDXID": "SPDXRef-40",
      "PackageVersion": "3.8.3",
      "PackageDownloadLocation": "NOASSERTION",
      "FilesAnalyzed": false,
      "ExternalRef": "PACKAGE-MANAGER purl pkg:pypi/flake8@3.8.3",
      "PackageLicenseConcluded": "NOASSERTION",
      "PackageLicenseDeclared": "MIT"
    },
    ...
```

```
$ gh sbom -c -l -r steiza/dependabot-example | jq
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.4",
  "serialNumber": "",
  "version": 1,
  "components": [
    {
      "type": "library",
      "name": "flake8",
      "version": "3.8.3",
      "purl": "pkg:pypi/flake8@3.8.3",
      "licenses": [
        {
          "Expression": "MIT"
        }
      ]
    },
    ...
```
