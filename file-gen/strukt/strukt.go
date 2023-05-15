package strukt

import (
	"github.com/rmarken5/lava/file-gen/property"
)

type (
	Strukt struct {
		name       string
		Properties []property.Property
	}
	StruktBuilder struct {
		struktModifiers []struktModifier
	}
	struktModifier func(strukt *Strukt)
	struktBuild    func(b *StruktBuilder) *Strukt
)

func (b *StruktBuilder) AddProp(prop property.Property) *StruktBuilder {
	b.struktModifiers = append(b.struktModifiers, func(strukt *Strukt) {
		strukt.Properties = append(strukt.Properties, prop)
	})
	return b
}

func (b *StruktBuilder) Named(name string) *StruktBuilder {
	b.struktModifiers = append(b.struktModifiers, func(strukt *Strukt) {
		strukt.name = name
	})
	return b
}

func printStrukt(strukt Strukt) string {

}
