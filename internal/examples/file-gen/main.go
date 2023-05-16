package main

import (
	"fmt"
	"github.com/rmarken5/lava/file-gen/property"
	"github.com/rmarken5/lava/file-gen/strukt"
)

func main() {
	properties, err := property.PrintProperties(func(b *property.PropertyBuilder) *property.Property {
		return b.Named("ID").OfType("int").Tagged("id").Build()
	})
	if err != nil {
		return
	}
	fmt.Printf("%s%s\n", properties, properties)

	printStrukt, err := strukt.PrintStrukt(func(b *strukt.StruktBuilder) *strukt.Strukt {
		return b.Named("MyStruct").AddPropBuilder(func(pb *property.PropertyBuilder) *property.Property {
			return pb.Named("ID").OfType("uuid.UUID").Tagged("id").Build()
		}).Build()
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(printStrukt))

}
