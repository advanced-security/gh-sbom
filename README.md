# gh-sbom

This is a `gh` CLI extension that outputs JSON SBOMs (in SPDX or CycloneDX format) for your GitHub repository using information from [Dependency graph](https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/about-the-dependency-graph).

It can optionally include license information with `-l`. License information comes from [ClearlyDefined](https://clearlydefined.io/)'s API.

Here's an example of generating a SPDX SBOM:
```
$ gh sbom -l | jq
{
  "spdxVersion": "SPDX-2.3",
  "dataLicense": "CC0-1.0",
  "SPDXID": "SPDXRef-DOCUMENT",
  "name": "github.com/advanced-security/gh-sbom",
  "documentNamespace": "https://spdx.org/spdxdocs/github.com/advanced-security/gh-sbom-81f6ee97-cae4-42a4-9be0-840bd3dde2a7",
  "creationInfo": {
    "creators": [
      "Organization: GitHub, Inc",
      "Tool: gh-sbom-0.0.8"
    ],
    "created": "2023-03-10T21:12:26Z"
  },
  "packages": [
    {
      "name": "gh-sbom",
      "SPDXID": "SPDXRef-mainPackage",
      "versionInfo": "",
      "downloadLocation": "git+https://github.com/advanced-security/gh-sbom.git",
      "filesAnalyzed": false,
      "externalRefs": [
        {
          "referenceCategory": "PACKAGE-MANAGER",
          "referenceType": "purl",
          "referenceLocator": "pkg:github/advanced-security/gh-sbom"
        }
      ],
      "licenseConcluded": "NOASSERTION",
      "licenseDeclared": "MIT",
      "supplier": "NOASSERTION"
    },
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
        "version": "0.0.8"
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
