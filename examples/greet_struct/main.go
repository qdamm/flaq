package main

import (
	"fmt"
	"strings"

	"github.com/qdamm/flaq"
)

type options struct {
	Name string `flaq:"-n, --name string  name of the person to greet"`
	Yell bool   `flaq:"    --yell         greet the person loudly"`
}

var opts = options{Name: "world"}

func main() {
	flaq.Struct(&opts)
	flaq.Parse()

	if opts.Yell {
		opts.Name = strings.ToUpper(opts.Name)
	}
	fmt.Printf("Hello %s!\n", opts.Name)
}
