package main

import (
	"fmt"
	"github.com/rmarken5/lava/internal/file-gen/file"
	"github.com/rmarken5/lava/internal/file-gen/property"
	"github.com/rmarken5/lava/internal/file-gen/strukt"
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

	fileBytes, err := file.PrintFile(func(fb *file.FileBuilder) *file.File {
		return fb.WithPackage("github.com/rmarken/lava/db/models").
			WithImports([]string{"fmt", "github.com/rmarken/lava/file-gen/file", "time.Time"}).
			WithStruktPrinter(strukt.PrintStrukt).
			WithStruktBuilder(func(b *strukt.StruktBuilder) *strukt.Strukt {
				return b.Named("MyFirst").
					AddPropertyPrinter(property.PrintProperties).
					AddPropertyBuilder(func(b *property.PropertyBuilder) *property.Property {
						return b.Named("File").
							OfType("file.File").
							Tagged("file").
							Build()
					}).
					AddPropertyBuilder(func(b *property.PropertyBuilder) *property.Property {
						return b.Named("Created_At").
							OfType("time.Time").
							Tagged("createdAt").
							Build()
					}).Build()
			}).
			WithStruktBuilder(func(b *strukt.StruktBuilder) *strukt.Strukt {
				return b.Named("MySecond").
					AddPropertyPrinter(property.PrintProperties).
					AddPropertyBuilder(func(b *property.PropertyBuilder) *property.Property {
						return b.Named("Eh").
							OfType("file.File").
							Tagged("file").
							Build()
					}).
					AddPropertyBuilder(func(b *property.PropertyBuilder) *property.Property {
						return b.Named("Created_At").
							OfType("time.Time").
							Tagged("createdAt").
							Build()
					}).Build()
			}).Build()
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(fileBytes))
}
