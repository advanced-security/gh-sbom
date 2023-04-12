# gh-sbom

This is a `gh` CLI extension that outputs JSON SBOMs (in SPDX or CycloneDX format) for your GitHub repository using information from [Dependency graph](https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/about-the-dependency-graph).

SPDX output use the [Dependency Graph SBOM API](https://docs.github.com/en/rest/dependency-graph/sboms?apiVersion=2022-11-28), which quickly generates the SBOM server-side, and as such is faster, works for large repositories, and always includes license information.

CycloneDX output is generating by assembling the dependency information from the Dependency Graph GraphQL API, and license information (if you specify `-l`) from [ClearlyDefined](https://clearlydefined.io/)'s API. As such, CycloneDX output is slower, and may not work for large repositories.

Here's an example of generating a SPDX SBOM:
```
$ gh sbom | jq
{
  "SPDXID": "SPDXRef-DOCUMENT",
  "creationInfo": {
    "created": "2023-04-12T18:41:40Z",
    "creators": [
      "Tool: GitHub.com-Dependency-Graph"
    ]
  },
  "dataLicense": "CC0-1.0",
  "documentDescribes": [
    "com.github.advanced-security/gh-sbom"
  ],
  "documentNamespace": "https://github.com/advanced-security/gh-sbom/dependency_graph/sbom-fa3abb267af77b5d",
  "name": "com.github.advanced-security/gh-sbom",
  "packages": [
    {
      "SPDXID": "SPDXRef-go-github.com/cli/go-gh-1.1.0",
      "downloadLocation": "NOASSERTION",
      "externalRefs": [
        {
          "referenceCategory": "PACKAGE-MANAGER",
          "referenceLocator": "pkg:golang/github.com/cli/go-gh@1.1.0",
          "referenceType": "purl"
        }
      ],
...
```

Or for CycloneDX use `-c`:
```
$ gh sbom -c -l | jq
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.4",
  "version": 1,
  "metadata": {
    "timestamp": "2023-03-10T21:14:23Z",
    "tools": [
      {
        "vendor": "advanced-security",
        "name": "gh-sbom",
        "version": "0.0.9"
      }
    ],
    "licenses": [
      {
        "expression": "CC0-1.0"
      }
    ]
  },
  "components": [
    {
      "type": "library",
      "group": "github.com/cli",
      "name": "go-gh",
      "version": "1.2.1",
      "purl": "pkg:golang/github.com/cli/go-gh@1.2.1"
    },
    ...
```

## Background

There is not another planned release, but bug reports are welcome via issues and questions are welcome via discussion.

## Requirements

This is an extension to the `gh` CLI. See [gh CLI installation](https://github.com/cli/cli#installation) if you haven't installed `gh` yet.

Once `gh` is installed, you can install this extension with:
```
$ gh ext install advanced-security/gh-sbom
```

If you want to upgrade to the latest version you can remove and reinstall the extension:
```
$ gh ext remove advanced-security/gh-sbom
$ gh ext install advanced-security/gh-sbom
```

Finally, if you are planning to run this on a GHES instance, you will need to be on: `GHES 3.8` or higher. 

## License

This project is licensed under the terms of the MIT open source license. Please refer to [LICENSE.md](./LICENSE.md) for the full terms.

## Support

Bug reports are welcome via issues and questions are welcome via discussion. Please refer to [SUPPORT.md](./SUPPORT.md) for details.
This project is provided as-is. See
