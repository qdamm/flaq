# flaq

[![Build Status](https://travis-ci.org/qdamm/flaq.svg?branch=master)](https://travis-ci.org/qdamm/flaq)
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

There are mainly 2 ways to use the library: using one variable per option or storing them in a struct.

The first approach relies on a similar API than the [flag](https://godoc.org/flag) package in stdlib:

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

The above code specifies a `--name` option with a `-n` shorthand, accepting a string argument.
It also defines a `--yell` bool option.

The second approach is to define options through a struct:

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

Struct fields tags are expected to follow the `-s, --long type  description` pattern, where:

- `-s` is the option shortand.
- `--long` is the option long form.
- `type` is the type of argument the option accepts, if any.
    It can be one of the following: `string`, `bool`, `int`, `count`, or `duration`.
- `description` is the option description.

In order to ease alignment for readability, there can be an arbitrary number of whitespaces
between the different sections.

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
