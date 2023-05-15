package main

import (
	"fmt"
	file_gen "github.com/rmarken5/lava/file-gen/property"
)

func main() {
	properties, err := file_gen.PrintProperties(func(b *file_gen.PropertyBuilder) *file_gen.Property {
		return b.Named("ID").OfType("int").Tagged("id").Build()
	})
	if err != nil {
		return
	}
	fmt.Printf("%s", properties)
}
