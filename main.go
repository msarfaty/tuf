package main

import (
	"fmt"
	"os"

	"mikesarfaty.com/tuf/pkg/parser"
)

func main() {
	start1 := `
module "foo" {
	bar = baz
}
`
	os.WriteFile("test.tf", []byte(start1[1:]), 0644)
	err := parser.MoveHclBlock(&parser.MoveOptions{
		Address:  "module.foo",
		FromFile: "test.tf",
		ToFile:   "out.tf",
	})
	fmt.Printf("%s", err)
}
