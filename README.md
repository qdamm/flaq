# flaq

[![Build Status](https://travis-ci.org/qdamm/flaq.svg?branch=master)](https://travis-ci.org/qdamm/flaq)
[![Coverage Status](https://coveralls.io/repos/github/qdamm/flaq/badge.svg?branch=master)](https://coveralls.io/github/qdamm/flaq?branch=master)
[![GoDoc](https://godoc.org/github.com/qdamm/flaq?status.svg)](https://godoc.org/github.com/qdamm/flaq)

flaq is a command-line options parsing library for Go.

The package follows [POSIX][1] / [GNU][2] conventions, it notably supports:

- short and long options
- options with required, optional, or no argument
- parsing options intermixed with args
- defining options in a struct
- customizable help usage
- option abbreviations

[1]: http://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html
[2]: https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html

## Usage

As an example, we will specify a `--name` option with a `-n` shorthand, accepting a string argument.
We will also define a `--yell` bool option.

There are mainly 2 approaches when using flaq.

### Using a struct

The first approach is to define options through a struct:

```go
package main

import (
	"fmt"
	"strings"

	"github.com/qdamm/flaq"
)

type Options struct {
	Name string `flaq:"-n, --name string  name of the person to greet"`
	Yell bool   `flaq:"    --yell         greet the person loudly"`
}

var opts = Options{Name: "world"}

func main() {
	flaq.Struct(&opts)
	flaq.Parse()

	if opts.Yell {
		opts.Name = strings.ToUpper(opts.Name)
	}
	fmt.Printf("Hello %s!\n", opts.Name)
}
```

Using a struct creates a self-documented piece of code. It also plays well with other projects
such as [caarlos0/env](https://github.com/caarlos0/env), if you want for example to read values
from environment variables as well.

See the [struct fields](#struct-fields) section for more information about this.

### Using a variable per option

The second approach relies on a similar API than the [flag](https://godoc.org/flag) package:

```go
package main

import (
	"fmt"

	"github.com/qdamm/flaq"
)

var name = "world"
var yell = false

func main() {
	flaq.String(&name, "name", "n", "name of the person to greet")
	flaq.Bool(&yell, "yell", "", "whether to yell or not", false)
	flaq.Parse()

	if yell {
		name = strings.ToUpper(name)
	}
	fmt.Printf("Hello %s!\n", name)
}
```

## Help usage

By default, help usage is printed whenever a `--help` or `-h` option is seen:

```Shell
$ greet --help
Usage: greet [options]

Options:
  -h, --help           show usage help
  -n, --name <string>  name of the person to greet
      --yell           greet loudly
```

This behavior can be removed or customized.

## Struct fields

Struct fields tags are expected to follow the `-s, --long type  description` pattern, where:

- `-s` is the option shortand.
- `--long` is the option long form.
- `type` is the type of argument the option accepts, if any.
    It can be one of the following: `string`, `bool`, `int`, `count`, or `duration`.
- `description` is the option description.

In order to ease alignment for readability, there can be an arbitrary number of whitespaces
between the different sections.

Here is a reference of the supported field types:

```go
type Options struct {
	Name     string        `flaq:"-n, --name string         name of the person to greet"`
	Yell     bool          `flaq:"    --yell                --yell will set the value to true"`
	Bool     bool          `flaq:"    --bool bool           --bool, --bool=true or --bool=false"`
	Int      int           `flaq:"    --int int             an int value"`
	Count    int           `flaq:"-c, --count count         -ccc will set this count value to 3"`
	Duration time.Duration `flaq:"    --duration duration   a duration eg. --duration=5min"`
}
```
