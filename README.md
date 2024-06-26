<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/shdoc"><img src="https://kaos.sh/r/shdoc.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/shdoc"><img src="https://kaos.sh/b/a4221ea2-3758-4fb6-adf0-08cd7199960a.svg" alt="codebeat badge" /></a>
  <a href="https://kaos.sh/y/shdoc"><img src="https://kaos.sh/y/d8aa5c8aa68f43f6aa91872929a1695f.svg" alt="Codacy badge" /></a>
  <a href="https://kaos.sh/c/shdoc"><img src="https://kaos.sh/c/shdoc.svg" alt="Coverage Status" /></a>
  <br/>
  <a href="https://kaos.sh/w/shdoc/ci"><img src="https://kaos.sh/w/shdoc/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/shdoc/codeql"><img src="https://kaos.sh/w/shdoc/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#usage-demo">Usage Demo</a> • <a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#test--coverage-status">Test & Coverage Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`shdoc` is a tool for viewing and exporting documentation for shell scripts.

### Usage Demo

[![demo](https://gh.kaos.st/shdoc-020.gif)](#usage-demo)

### Installation

#### From source

Make sure you have a working Go 1.21+ workspace ([instructions](https://go.dev/doc/install)), then:

```bash
go install github.com/essentialkaos/shdoc@latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/shdoc/latest).

To install the latest prebuilt version of bibop, do:

```bash
bash <(curl -fsSL https://apps.kaos.st/get) shdoc
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```
sudo shdoc --completion=bash 1> /etc/bash_completion.d/shdoc
```


ZSH:
```
sudo shdoc --completion=zsh 1> /usr/share/zsh/site-functions/shdoc
```


Fish:
```
sudo shdoc --completion=fish 1> /usr/share/fish/vendor_completions.d/shdoc.fish
```

### Usage

<img src=".github/images/usage.svg" />

### Test & Coverage Status

| Branch | CI       | Coveralls |
|--------|----------|-----------|
| `master` | [![CI](https://kaos.sh/w/shdoc/ci.svg?branch=master)](https://kaos.sh/w/shdoc/ci?query=branch:master) | [![Coverage Status](https://kaos.sh/c/shdoc.svg?branch=master)](https://kaos.sh/c/shdoc?branch=master) |
| `develop` | [![CI](https://kaos.sh/w/shdoc/ci.svg?branch=develop)](https://kaos.sh/w/shdoc/ci?query=branch:develop) | [![Coverage Status](https://kaos.sh/c/shdoc.svg?branch=develop)](https://kaos.sh/c/shdoc?branch=develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
