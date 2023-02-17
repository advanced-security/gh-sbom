# gh-dependabot

You can install this extension with `$ gh ext install steiza/gh-sbom`.

This [gh CLI extension](https://docs.github.com/en/github-cli/github-cli/using-github-cli-extensions) outputs a SPDX SBOM from your [Dependency graph](https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/about-the-dependency-graph) from the command line:

```
$ gh sbom -r steiza/dependabot-example
{"SPDXVersion":"SPDX-2.3","DataLicense":"CC0-1.0","SPDXID":"SPDXRef-DOCUMENT","DocumentName":"github.com/steiza/dependabot-example","Creator":"Tool https://github.com/steiza/gh-sbom","Created":"2023-02-17T21:44:07Z","Packages":[{"PackageName":"pillow","SPDXID":"SPDXRef-0","PackageVersion":"8.1.0","PackageDownloadLocation":"NOASSERTION","FilesAnalyzed":false}...
```
