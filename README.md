# gh-sbom

This is a `gh` CLI extension that outputs JSON SBOMs (in SPDX or CycloneDX format) for your GitHub repository using information from [Dependency graph](https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/about-the-dependency-graph).

It can optionally include license information with `-l`. License information comes from [ClearlyDefined](https://clearlydefined.io/)'s API.

Here's an example of generating a SPDX SBOM:
```
$ gh sbom -l -r steiza/dependabot-example | jq
{
  "spdxVersion": "SPDX-2.3",
  "dataLicense": "CC0-1.0",
  "SPDXID": "SPDXRef-DOCUMENT",
  "name": "github.com/steiza/dependabot-example",
  "documentNamespace": "https://spdx.org/spdxdocs/github.com/steiza/dependabot-example-316abfe8-962e-4d21-9887-a347027bb216",
  "creationInfo": {
    "creators": [
      "Organization: GitHub, Inc",
      "Tool: gh-sbom"
    ],
    "created": "2023-03-08T15:19:43Z"
  },
  "packages": [
    {
      "name": "urllib3",
      "SPDXID": "SPDXRef-35",
      "versionInfo": "1.25.10",
      "downloadLocation": "NOASSERTION",
      "filesAnalyzed": false,
      "externalRefs": [
        {
          "referenceCategory": "PACKAGE-MANAGER",
          "referenceType": "purl",
          "referenceLocator": "pkg:pypi/urllib3@1.25.10"
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
$ gh sbom -c -l -r steiza/dependabot-example | jq
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.4",
  "version": 1,
  "metadata": {
    "timestamp": "2023-03-08T15:21:32Z",
    "tools": [
      {
        "vendor": "advanced-security",
        "name": "gh-sbom"
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
      "name": "urllib3",
      "version": "1.25.10",
      "purl": "pkg:pypi/urllib3@1.25.10",
      "licenses": [
        {
          "expression": "MIT"
        }
      ]
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

## License

This project is licensed under the terms of the MIT open source license. Please refer to [LICENSE.md](./LICENSE.md) for the full terms.

## Support

Bug reports are welcome via issues and questions are welcome via discussion. Please refer to [SUPPORT.md](./SUPPORT.md) for details.
This project is provided as-is. See
