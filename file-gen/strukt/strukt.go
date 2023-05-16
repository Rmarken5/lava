package strukt

import (
	"bytes"
	"embed"
	"github.com/rmarken5/lava/file-gen/property"
	"text/template"
)

//go:embed strukt.gotmpl
var fs embed.FS

type (
	Strukt struct {
		name       string
		properties []*property.Property
	}
	StruktBuilder struct {
		struktModifiers []struktModifier
	}
	struktModifier func(strukt *Strukt)
	struktBuild    func(b *StruktBuilder) *Strukt
)

func (b *StruktBuilder) AddPropBuilder(fn func(pb *property.PropertyBuilder) *property.Property) *StruktBuilder {
	builder := &property.PropertyBuilder{}
	b.struktModifiers = append(b.struktModifiers, func(strukt *Strukt) {
		strukt.properties = append(strukt.properties, fn(builder))
	})
	return b
}

func (b *StruktBuilder) Named(name string) *StruktBuilder {
	b.struktModifiers = append(b.struktModifiers, func(strukt *Strukt) {
		strukt.name = name
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

func PrintStrukt(b struktBuild) ([]byte, error) {
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
