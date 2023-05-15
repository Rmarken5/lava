package property

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed property.tmpl
var tmpl embed.FS

type (
	Property struct {
		Name, DataType, Tag string
	}

	PropertyBuilder struct {
		modifiers []modifier `json:"modifiers"`
	}
	modifier func(p *Property)

	build func(b *PropertyBuilder) *Property
)

func (b *PropertyBuilder) Named(name string) *PropertyBuilder {
	b.modifiers = append(b.modifiers, func(p *Property) {
		p.Name = name
	})
	return b
}

func (b *PropertyBuilder) OfType(dataType string) *PropertyBuilder {
	b.modifiers = append(b.modifiers, func(p *Property) {
		p.DataType = dataType
	})
	return b
}

func (b *PropertyBuilder) Tagged(tag string) *PropertyBuilder {
	b.modifiers = append(b.modifiers, func(p *Property) {
		p.Tag = tag
	})
	return b
}

func (b *PropertyBuilder) Build() *Property {
	var property Property
	for _, mod := range b.modifiers {
		mod(&property)
	}
	return &property
}

func PrintProperties(b build) ([]byte, error) {
	builder := &PropertyBuilder{}
	prop := b(builder)
	return printProperties(prop)
}

func printProperties(prop *Property) ([]byte, error) {
	fs, err := template.ParseFS(tmpl, "property.tmpl")
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(nil)

	err = fs.Execute(buff, prop)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
