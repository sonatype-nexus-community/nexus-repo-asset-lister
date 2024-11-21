# Sonatype Nexus Repository Asset Lister

<!-- Badges Section -->
[![shield_gh-workflow-test]][link_gh-workflow-test]
[![shield_license]][license_file]
<!-- Add other badges or shields as appropriate -->

---

Utility to list identity of all assets in a repository within Sonatype Nexus Repository.

- [What does this tool do?](#what-does-this-tool-do)
- [Installation](#installation)
- [Usage](#usage)
- [Development](#development)
- [The Fine Print](#the-fine-print)

## What does this tool do?

This tool produces a JSON file per Proxy Repository that contains basic identity information for every Asset in the Repository.

A JSON files contents might look similar to this:

```json
[
      {
            "p": "/data-uri-to-buffer/-/data-uri-to-buffer-4.0.1.tgz",
            "h": {
                  "md5": "a4d7239289232437622e6342e771a153",
                  "sha1": "d8feb2b2881e6a4f58c2e08acfd0e2834e26222e",
                  "sha256": null
            }
      },
      ...
]
```

This tool purposefully skips Proxy Repositories of format `docker`.

## Installation

Obtain the binary for your Operating System and Architecture from the [GitHub Releases page](https://github.com/sonatype-nexus-community/nexus-repo-asset-lister/releases).

## Usage

You can see all options by running:

```
./nexus-repo-asset-lister -help

usage: nexus-repo-asset-lister [OPTIONS]
  -X    Enable debug logging
  -o string
        Directory to write asset lists to (default "/Users/phorton/Documents/GitHub/sonatype-nexus-community/nexus-repo-asset-lister")
  -password string
        Password used to authenticate to Sonatype Nexus Repository (can also be set using the environment variable NXRM_PASSWORD)
  -skipped
        Whether to ouptut skipped assets to a separate '-skipped.json' file
  -url string
        URL including protocol to your Sonatype Nexus Repository Manager (default "http://localhost:8081")
  -username string
        Username used to authenticate to Sonatype Nexus Repository (can also be set using the environment variable NXRM_USERNAME)
```

In order to produce the expected JSON files, you should run `nexus-repo-asset-lister` with Administrator credentials against your Sonatype Nexus Repository.

The URL of the Sonatype Nexus Repository sever is specified with the `-url` argument and should contain the protcol (e.g. `https://`) and any context path you may have set for the installation.

Credentials can be supplied as command line arguments or via Environment Variables.

```
NXRM_USERNAME=username NXRM_PASSWORD=password ./nexus-repo-asset-lister -url https://my-nexus-repository.tld
```

You can use User Token instead of actual username and password.

## Development

See [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## The Fine Print

Remember:

This project is part of the [Sonatype Nexus Community](https://github.com/sonatype-nexus-community) organization, which is not officially supported by Sonatype. Please review the latest pull requests, issues, and commits to understand this project's readiness for contribution and use.

* File suggestions and requests on this repo through GitHub Issues, so that the community can pitch in
* Use or contribute to this project according to your organization's policies and your own risk tolerance
* Don't file Sonatype support tickets related to this project— it won't reach the right people that way

Last but not least of all - have fun!

<!-- Links Section -->
[shield_gh-workflow-test]: https://img.shields.io/github/actions/workflow/status/sonatype-nexus-community/nexus-repo-asset-lister/build.yml?branch=main&logo=GitHub&logoColor=white "build"
[shield_license]: https://img.shields.io/github/license/sonatype-nexus-community/nexus-repo-asset-lister?logo=open%20source%20initiative&logoColor=white "license"

[link_gh-workflow-test]: https://github.com/sonatype-nexus-community/nexus-repo-asset-lister/actions/workflows/build.yml?query=branch%3Amain
[license_file]: https://github.com/sonatype-nexus-community/nexus-repo-asset-lister/blob/main/LICENSE