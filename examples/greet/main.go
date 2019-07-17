package main

import (
	"fmt"
	"strings"

	"github.com/qdamm/flaq"
)

var name = "world"
var yell = false

func main() {
	flaq.String(&name, "name", "n", "name of the person to greet")
	flaq.Bool(&yell, "yell", "", "greet the person loudly", false)
	flaq.Parse()

	if yell {
		name = strings.ToUpper(name)
	}
	fmt.Printf("Hello %s!\n", name)
}
