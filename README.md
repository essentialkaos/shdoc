<p align="center"><a href="#readme"><img src="https://gh.kaos.st/shdoc.svg"/></a></p>

<p align="center">
  <a href="https://github.com/essentialkaos/shdoc/actions"><img src="https://github.com/essentialkaos/shdoc/workflows/CI/badge.svg" alt="GitHub Actions Status" /></a>
  <a href="https://github.com/essentialkaos/shdoc/actions?query=workflow%3ACodeQL"><img src="https://github.com/essentialkaos/shdoc/workflows/CodeQL/badge.svg" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/shdoc"><img src="https://goreportcard.com/badge/github.com/essentialkaos/shdoc"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-shdoc-master"><img alt="codebeat badge" src="https://codebeat.co/badges/a4221ea2-3758-4fb6-adf0-08cd7199960a" /></a>
  <a href='https://coveralls.io/github/essentialkaos/shdoc?branch=master'><img src='https://coveralls.io/repos/github/essentialkaos/shdoc/badge.svg?branch=master' alt='Coverage Status' /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#usage-demo">Usage Demo</a> • <a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#test--coverage-status">Test & Coverage Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`shdoc` is a tool for viewing and exporting documentation for shell scripts.

### Usage Demo

[![demo](https://gh.kaos.st/shdoc-020.gif)](#usage-demo)

### Installation

#### From source

Make sure you have a working Go 1.14+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/essentialkaos/shdoc
```

If you want to update `shdoc` to latest stable release, do:

```
go get -u github.com/essentialkaos/shdoc
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.st/shdoc/latest).

To install the latest prebuilt version of bibop, do:

```bash
bash <(curl -fsSL https://apps.kaos.st/get) shdoc
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```
[sudo] shdoc --completion=bash 1> /etc/bash_completion.d/shdoc
```


ZSH:
```
[sudo] shdoc --completion=zsh 1> /usr/share/zsh/site-functions/shdoc
```


Fish:
```
[sudo] shdoc --completion=fish 1> /usr/share/fish/vendor_completions.d/shdoc.fish
```

### Usage

```
Usage: shdoc {options} file

Options

  --output, -o file      Path to output file
  --template, -t name    Name of template
  --name, -n name        Overwrite default name
  --no-color, -nc        Disable colors in output
  --help, -h             Show this help message
  --version, -v          Show version

Examples:

  shdoc script.sh
  Parse shell script and show docs in console

  shdoc script.sh -t markdown -o my_script.md
  Parse shell script and save docs using given export template

  shdoc script.sh someEntity
  Parse shell script and show docs for some constant, variable or method

```

### Test & Coverage Status

| Branch | TravisCI | Coveralls |
|--------|----------|---------|
| `master` | [![CI](https://github.com/essentialkaos/shdoc/workflows/CI/badge.svg?branch=master)](https://github.com/essentialkaos/shdoc/actions) | [![Coverage Status](https://coveralls.io/repos/github/essentialkaos/shdoc/badge.svg?branch=master)](https://coveralls.io/github/essentialkaos/shdoc?branch=master) |
| `develop` | [![CI](https://github.com/essentialkaos/shdoc/workflows/CI/badge.svg?branch=develop)](https://github.com/essentialkaos/shdoc/actions) | [![Coverage Status](https://coveralls.io/repos/github/essentialkaos/shdoc/badge.svg?branch=develop)](https://coveralls.io/github/essentialkaos/shdoc?branch=develop)

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
