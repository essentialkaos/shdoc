# SHDoc [![Build Status](https://travis-ci.org/essentialkaos/shdoc.svg?branch=master)](https://travis-ci.org/essentialkaos/shdoc) [![Coverage Status](https://coveralls.io/repos/github/essentialkaos/shdoc/badge.svg?branch=master)](https://coveralls.io/github/essentialkaos/shdoc?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/essentialkaos/shdoc)](https://goreportcard.com/report/github.com/essentialkaos/shdoc) [![codebeat badge](https://codebeat.co/badges/a4221ea2-3758-4fb6-adf0-08cd7199960a)](https://codebeat.co/projects/github-com-essentialkaos-shdoc-master) [![License](https://gh.kaos.st/ekol.svg)](https://essentialkaos.com/ekol)

Tool for viewing and exporting docs for shell scripts.

* [Usage Demo](#usage-demo)
* [Installation](#installation)
* [Usage](#usage)
* [Test & Coverage Status](#test--coverage-status)
* [Contributing](#contributing)
* [License](#license)

### Usage Demo

[![demo](https://gh.kaos.st/shdoc-020.gif)](#usage-demo)

### Installation

#### From source

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.6+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/essentialkaos/shdoc
```

If you want to update `shdoc` to latest stable release, do:

```
go get -u github.com/essentialkaos/shdoc
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.st/shdoc/latest).

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
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/shdoc.svg?branch=master)](https://travis-ci.org/essentialkaos/shdoc) | [![Coverage Status](https://coveralls.io/repos/github/essentialkaos/shdoc/badge.svg?branch=master)](https://coveralls.io/github/essentialkaos/shdoc?branch=master) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/shdoc.svg?branch=develop)](https://travis-ci.org/essentialkaos/shdoc) | [![Coverage Status](https://coveralls.io/repos/github/essentialkaos/shdoc/badge.svg?branch=develop)](https://coveralls.io/github/essentialkaos/shdoc?branch=develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>