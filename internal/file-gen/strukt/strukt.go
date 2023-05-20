package strukt

import (
	"bytes"
	"embed"
	"github.com/rmarken5/lava/internal/file-gen/property"
	"text/template"
)

//go:embed strukt.gotmpl
var fs embed.FS

type (
	Strukt struct {
		Name             string
		PropertyBuilders []property.BuildProperty
		PropertyPrinter  propertyPrinter
	}
	StruktBuilder struct {
		struktModifiers []struktModifier
	}
	struktModifier func(strukt *Strukt)
	BuildStrukt    func(b *StruktBuilder) *Strukt

	propertyPrinter func(b property.BuildProperty) ([]byte, error)
)

func (b *StruktBuilder) Named(name string) *StruktBuilder {
	b.struktModifiers = append(b.struktModifiers, func(strukt *Strukt) {
		strukt.Name = name
	})
	return b
}

func (b *StruktBuilder) AddPropertyBuilder(fn property.BuildProperty) *StruktBuilder {
	b.struktModifiers = append(b.struktModifiers, func(strukt *Strukt) {
		strukt.PropertyBuilders = append(strukt.PropertyBuilders, fn)
	})
	return b
}
func (b *StruktBuilder) AddPropertyPrinter(fn propertyPrinter) *StruktBuilder {
	b.struktModifiers = append(b.struktModifiers, func(strukt *Strukt) {
		strukt.PropertyPrinter = fn
	})
	return b
}
func (b *StruktBuilder) Build() *Strukt {
	strukt := &Strukt{}
	for _, mod := range b.struktModifiers {
		mod(strukt)
	}
	return strukt
}

func PrintStrukt(b BuildStrukt) ([]byte, error) {
	builder := &StruktBuilder{}
	strukt := b(builder)
	return printStrukt(strukt)
}

func printStrukt(strukt *Strukt) ([]byte, error) {
	file, err := template.ParseFS(fs, "strukt.gotmpl")
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(nil)

	err = file.Execute(buff, strukt)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
