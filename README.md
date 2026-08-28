<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/y/shdoc"><img src="https://app.codacy.com/project/badge/Grade/d8aa5c8aa68f43f6aa91872929a1695f" alt="Codacy badge" /></a>
  <a href="https://kaos.sh/c/shdoc"><img src="https://coveralls.io/repos/github/essentialkaos/shdoc/badge.svg" alt="Coverage Status" /></a>
  <br/>
  <a href="https://kaos.sh/w/shdoc/ci"><img src="https://github.com/essentialkaos/shdoc/actions/workflows/ci.yml/badge.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/shdoc/codeql"><img src="https://github.com/essentialkaos/shdoc/actions/workflows/codeql.yml/badge.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#usage-demo">Usage Demo</a> • <a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#test--coverage-status">Test & Coverage Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`shdoc` is a tool for viewing and exporting documentation for shell scripts. The documentation format is described in the [specification](SPECIFICATION.md).

### Usage Demo

[![demo](https://github.com/user-attachments/assets/693ae7df-63ee-42ff-af46-e95fc652fd25)](#usage-demo)

### Installation

#### From source

Make sure you have a working Go [1.25+](https://github.com/essentialkaos/.github/blob/master/GO-VERSION-SUPPORT.md) workspace ([instructions](https://go.dev/doc/install)), then:

```bash
go install github.com/essentialkaos/shdoc@latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/shdoc/latest).

To install the latest prebuilt version of `shdoc`, do:

```bash
bash <(curl -fsSL https://apps.kaos.st/get) shdoc
```

### Upgrading

Since version `0.11.0` you can update `shdoc` to the latest release using [self-update feature](https://github.com/essentialkaos/.github/blob/master/APPS-UPDATE.md):

```bash
shdoc --update
```

This command will runs a self-update in interactive mode. If you want to run a quiet update (_no output_), use the following command:

```bash
shdoc --update=quiet
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
shdoc --completion=bash | sudo tee /etc/bash_completion.d/shdoc > /dev/null
```

ZSH:
```bash
shdoc --completion=zsh | sudo tee /usr/share/zsh/site-functions/shdoc > /dev/null
```

Fish:
```bash
shdoc --completion=fish | sudo tee /usr/share/fish/vendor_completions.d/shdoc.fish > /dev/null
```

### Usage

<img src=".github/images/usage.svg" />

### Test & Coverage Status

| Branch | CI       | Coveralls |
|--------|----------|-----------|
| `master` | [![CI](https://github.com/essentialkaos/shdoc/actions/workflows/ci.yml/badge.svg?branch=master)](https://kaos.sh/w/shdoc/ci?query=branch:master) | [![Coverage Status](https://coveralls.io/repos/github/essentialkaos/shdoc/badge.svg?branch=master)](https://kaos.sh/c/shdoc?branch=master) |
| `develop` | [![CI](https://github.com/essentialkaos/shdoc/actions/workflows/ci.yml/badge.svg?branch=develop)](https://kaos.sh/w/shdoc/ci?query=branch:develop) | [![Coverage Status](https://coveralls.io/repos/github/essentialkaos/shdoc/badge.svg?branch=develop)](https://kaos.sh/c/shdoc?branch=develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/.github/blob/master/CONTRIBUTING.md).

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://kaos.dev"><img src="https://raw.githubusercontent.com/essentialkaos/.github/refs/heads/master/images/ekgh.svg"/></a></p>
