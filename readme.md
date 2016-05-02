## SHDoc [![Release](https://img.shields.io/github/release/essentialkaos/shdoc.svg)](https://github.com/essentialkaos/shdoc/releases/latest)

Tool for viewing and exporting docs for shell scripts.

* [Installation](#installation)
* [Usage](#usage)
* [Test & Coverage Status](#test--coverage-status)
* [Contributing](#contributing)
* [License](#license)

#### Installation

````
go get github.com/essentialkaos/shdoc
````

#### Usage

```
Usage: shdoc <options> file

Options:

  --output, -o file      Path to output file
  --template, -t file    Path to template file
  --name, -n name        Overwrite default name
  --no-color, -nc        Disable colors in output
  --help, -h             Show this help message
  --version, -v          Show version

Examples:

  shdoc script.sh
  Parse shell script and show docs in console

  shdoc script.sh -t path/to/template.tpl -o my_script.md
  Parse shell script and save docs using given export template

  shdoc script.sh someEntity
  Parse shell script and show docs for some constant, variable or method

```

#### Test & Coverage Status

| Branch | TravisCI | CodeCov |
|--------|----------|---------|
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/shdoc.svg?branch=master)](https://travis-ci.org/essentialkaos/shdoc) | [![codecov.io](https://codecov.io/github/essentialkaos/shdoc/coverage.svg?branch=master)](https://codecov.io/github/essentialkaos/shdoc?branch=master) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/shdoc.svg?branch=develop)](https://travis-ci.org/essentialkaos/shdoc) | [![codecov.io](https://codecov.io/github/essentialkaos/shdoc/coverage.svg?branch=develop)](https://codecov.io/github/essentialkaos/shdoc?branch=develop) |

#### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

#### License

[EKOL](https://essentialkaos.com/ekol)
