# SHDoc [![Release](https://img.shields.io/github/release/essentialkaos/shdoc.svg)](https://github.com/essentialkaos/shdoc/releases/latest)

Tool for viewing and exporting docs for shell scripts.

### Installation

````
go get github.com/essentialkaos/shdoc
````

#### Usage

    Usage: shdoc <options> file
    
    Options:
    
      --output, -o file      Path to output file
      --template, -t file    Path to template file
      --name, -n name        Overwrite default name
      --no-color, -nc        Disable colors in output
      --help, -h             Show this help message
      --version, -v          Show version

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.io/shdoc/).

### Test & Coverage Status

| Branch | TravisCI | CodeCov |
|--------|----------|---------|
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/shdoc.svg?branch=master)](https://travis-ci.org/essentialkaos/shdoc) | [![codecov.io](https://codecov.io/github/essentialkaos/shdoc/coverage.svg?branch=master)](https://codecov.io/github/essentialkaos/shdoc?branch=master) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/shdoc.svg?branch=develop)](https://travis-ci.org/essentialkaos/shdoc) | [![codecov.io](https://codecov.io/github/essentialkaos/shdoc/coverage.svg?branch=develop)](https://codecov.io/github/essentialkaos/shdoc?branch=develop) |

### License

[EKOL](https://essentialkaos.com/ekol)
