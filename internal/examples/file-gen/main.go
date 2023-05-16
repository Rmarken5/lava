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

	sBuilder := strukt.StruktBuilder{}
	s := sBuilder.Named("MyStruct").AddPropertyPrinter(property.PrintProperties).AddPropertyBuilder(func(b *property.PropertyBuilder) *property.Property {
		return b.Named("ID").OfType("string").Tagged("id").Build()
	}).Build()

	printer, _ := s.PropertyPrinter(s.PropertyBuilders[0])

	fmt.Println(string(printer))

	printStrukt, err := strukt.PrintStrukt(func(b *strukt.StruktBuilder) *strukt.Strukt {
		s := b.Named("MyStruct").AddPropertyPrinter(property.PrintProperties).AddPropertyBuilder(func(b *property.PropertyBuilder) *property.Property {
			return b.Named("ID").OfType("string").Tagged("id").Build()
		}).Build()
		return s
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(printStrukt))

}
