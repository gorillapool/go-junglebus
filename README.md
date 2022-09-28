# JungleBus: Go Client
> Get started using [JungleBus](https://junglebus.gorillapool.io) in five minutes

[![Release](https://img.shields.io/github/release-pre/GorillaPool/go-junglebus.svg?logo=github&style=flat&v=2)](https://github.com/GorillaPool/go-junglebus/releases)
[![Build Status](https://img.shields.io/github/workflow/status/GorillaPool/go-junglebus/run-go-tests?logo=github&v=2)](https://github.com/GorillaPool/go-junglebus/actions)
[![Report](https://goreportcard.com/badge/github.com/GorillaPool/go-junglebus?style=flat&v=2)](https://goreportcard.com/report/github.com/GorillaPool/go-junglebus)
[![codecov](https://codecov.io/gh/GorillaPool/go-junglebus/branch/master/graph/badge.svg?v=2)](https://codecov.io/gh/GorillaPool/go-junglebus)
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://api.mergify.com/v1/badges/GorillaPool/go-junglebus&style=flat&v=2)](https://mergify.io)
[![Go](https://img.shields.io/github/go-mod/go-version/GorillaPool/go-junglebus?v=2)](https://golang.org/)
<br>
[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod&v=2)](https://gitpod.io/#https://github.com/GorillaPool/go-junglebus)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat&v=2)](https://github.com/RichardLitt/standard-readme)
[![Makefile Included](https://img.shields.io/badge/Makefile-Supported%20-brightgreen?=flat&logo=probot&v=2)](Makefile)
[![Sponsor](https://img.shields.io/badge/sponsor-mrz1836-181717.svg?logo=github&style=flat&v=2)](https://github.com/sponsors/GorillaPool)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat&v=2)](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=go-junglebusclient&utm_term=go-junglebusclient&utm_content=go-junglebusclient)

<br/>

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/GorillaPool/go-junglebus"
    "github.com/GorillaPool/go-junglebus/models"
)

func main() {
    junglebusClient, err := junglebus.New(
        junglebus.WithHTTP("https://junglebus.gorillapool.io"),
    )
    if err != nil {
        log.Fatalln(err.Error())
    }

    subscriptionID := "..." // fill in the ID for the subscription
    fromBlock := uint64(750000)

    eventHandler := junglebus.EventHandler{
        // do not set this function to leave out mined transactions
        OnTransaction: func(tx *models.TransactionResponse) {
            log.Printf("[TX]: %d: %v", tx.BlockHeight, tx.Id)
        },
        // do not set this function to leave out mempool transactions
        OnMempool: func(tx *models.TransactionResponse) {
            log.Printf("[MEMPOOL TX]: %v", tx.Id)
        },
        OnStatus: func(status *models.ControlResponse) {
            log.Printf("[STATUS]: %v", status)
        },
        OnError: func(err error) {
            log.Printf("[ERROR]: %v", err)
        },
    }

    var subscription *junglebus.Subscription
    if subscription, err = junglebusClient.Subscribe(context.Background(), subscriptionID, fromBlock, eventHandler); err != nil {
        log.Printf("ERROR: failed getting subscription %s", err.Error())
    } else {
        time.Sleep(10 * time.Second) // stop after 10 seconds
        if err = subscription.Unsubscribe(); err != nil {
            log.Printf("ERROR: failed unsubscribing %s", err.Error())
        }
    }
}
```

## Table of Contents
- [What is JungleBus?](#what-is-junglebus)
- [Installation](#installation)
- [Documentation](#documentation)
- [Examples & Tests](#examples--tests)
- [Benchmarks](#benchmarks)
- [Code Standards](#code-standards)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

<br/>

## What is JungleBus?
[Read more about JungleBus](https://getjunglebus.io)

<br/>

## Installation

**go-junglebusclient** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/GorillaPool/go-junglebus
```

<br/>

## Documentation
View the generated [documentation](https://pkg.go.dev/github.com/GorillaPool/go-junglebus)

[![GoDoc](https://godoc.org/github.com/GorillaPool/go-junglebus?status.svg&style=flat&v=2)](https://pkg.go.dev/github.com/GorillaPool/go-junglebus)

<br/>

<details>
<summary><strong><code>Repository Features</code></strong></summary>
<br/>

This repository was created using [MrZ's `go-template`](https://github.com/mrz1836/go-template#about)

#### Built-in Features
- Continuous integration via [GitHub Actions](https://github.com/features/actions)
- Build automation via [Make](https://www.gnu.org/software/make)
- Dependency management using [Go Modules](https://github.com/golang/go/wiki/Modules)
- Code formatting using [gofumpt](https://github.com/mvdan/gofumpt) and linting with [golangci-lint](https://github.com/golangci/golangci-lint) and [yamllint](https://yamllint.readthedocs.io/en/stable/index.html)
- Unit testing with [testify](https://github.com/stretchr/testify), [race detector](https://blog.golang.org/race-detector), code coverage [HTML report](https://blog.golang.org/cover) and [Codecov report](https://codecov.io/)
- Releasing using [GoReleaser](https://github.com/goreleaser/goreleaser) on [new Tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging)
- Dependency scanning and updating thanks to [Dependabot](https://dependabot.com) and [Nancy](https://github.com/sonatype-nexus-community/nancy)
- Security code analysis using [CodeQL Action](https://docs.github.com/en/github/finding-security-vulnerabilities-and-errors-in-your-code/about-code-scanning)
- Automatic syndication to [pkg.go.dev](https://pkg.go.dev/) on every release
- Generic templates for [Issues and Pull Requests](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/configuring-issue-templates-for-your-repository) in Github
- All standard Github files such as `LICENSE`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, and `SECURITY.md`
- Code [ownership configuration](.github/CODEOWNERS) for Github
- All your ignore files for [vs-code](.editorconfig), [docker](.dockerignore) and [git](.gitignore)
- Automatic sync for [labels](.github/labels.yml) into Github using a pre-defined [configuration](.github/labels.yml)
- Built-in powerful merging rules using [Mergify](https://mergify.io/)
- Welcome [new contributors](.github/mergify.yml) on their first Pull-Request
- Follows the [standard-readme](https://github.com/RichardLitt/standard-readme/blob/master/spec.md) specification
- [Visual Studio Code](https://code.visualstudio.com) configuration with [Go](https://code.visualstudio.com/docs/languages/go)
- (Optional) [Slack](https://slack.com), [Discord](https://discord.com) or [Twitter](https://twitter.com) announcements on new Github Releases
- (Optional) Easily add [contributors](https://allcontributors.org/docs/en/bot/installation) in any Issue or Pull-Request

</details>

<details>
<summary><strong><code>Package Dependencies</code></strong></summary>
<br/>

- [stretchr/testify](https://github.com/stretchr/testify)
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

Releases are automatically created when you create a new [git tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging)!

If you want to manually make releases, please install GoReleaser:

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to Github and can be installed:
- **using make:** `make install-releaser`
- **using brew:** `brew install goreleaser`

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

<br/>

### Automatic Releases on Tag Creation (recommended)
Automatic releases via [Github Actions](.github/workflows/release.yml) from creating a new tag:
```shell
make tag version=1.2.3
```

<br/>

### Manual Releases (optional)
Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production (manually).

<br/>

</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
all                           Runs multiple commands
clean                         Remove previous builds and any cached data
clean-mods                    Remove all the Go mod cache
coverage                      Shows the test coverage
diff                          Show the git diff
generate                      Runs the go generate command in the base of the repo
godocs                        Sync the latest tag with GoDocs
help                          Show this help message
install                       Install the application
install-all-contributors      Installs all contributors locally
install-go                    Install the application (Using Native Go)
install-releaser              Install the GoReleaser application
lint                          Run the golangci-lint application (install if not found)
release                       Full production release (creates release in Github)
release                       Runs common.release then runs godocs
release-snap                  Test the full release (build binaries)
release-test                  Full production test release (everything except deploy)
replace-version               Replaces the version in HTML/JS (pre-deploy)
tag                           Generate a new tag and push (tag version=0.0.0)
tag-remove                    Remove a tag if found (tag-remove version=0.0.0)
tag-update                    Update an existing tag to current commit (tag-update version=0.0.0)
test                          Runs lint and ALL tests
test-ci                       Runs all tests via CI (exports coverage)
test-ci-no-race               Runs all tests via CI (no race) (exports coverage)
test-ci-short                 Runs unit tests via CI (exports coverage)
test-no-lint                  Runs just tests
test-short                    Runs vet, lint and tests (excludes integration tests)
test-unit                     Runs tests and outputs coverage
uninstall                     Uninstall the application (and remove files)
update-contributors           Regenerates the contributors html/list
update-linter                 Update the golangci-lint package (macOS only)
vet                           Run the Go vet application
```
</details>

<br/>

## Examples & Tests
All unit tests and [examples](examples) run via [Github Actions](https://github.com/GorillaPool/go-junglebus/actions) and
uses [Go version 1.18.x](https://golang.org/doc/go1.18). View the [configuration file](.github/workflows/run-tests.yml).

<br/>

Run all tests (including integration tests)
```shell script
make test
```

<br/>

Run tests (excluding integration tests)
```shell script
make test-short
```

<br/>

## Benchmarks
Run the Go benchmarks:
```shell script
make bench
```

<br/>

## Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## Usage
Checkout all the [examples](examples)!

<br/>

## Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

<br/>

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/GorillaPool) :clap:
or by making a [**bitcoin donation**](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=go-junglebusclient&utm_term=go-junglebusclient&utm_content=go-junglebusclient) to ensure this journey continues indefinitely! :rocket:

[![Stars](https://img.shields.io/github/stars/GorillaPool/go-junglebus?label=Please%20like%20us&style=social&v=2)](https://github.com/GorillaPool/go-junglebus/stargazers)

<br/>

### Contributors ✨
Thank you to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/icellan"><img src="https://avatars.githubusercontent.com/u/4411176?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Siggi</b></sub></a><br /><a href="#infra-icellan" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="https://github.com/GorillaPool/go-junglebus/commits?author=icellan" title="Code">💻</a> <a href="#security-icellan" title="Security">🛡️</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

> This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification.

<br/>

## License

[![License](https://img.shields.io/github/license/GorillaPool/go-junglebus.svg?style=flat&v=2)](LICENSE)
