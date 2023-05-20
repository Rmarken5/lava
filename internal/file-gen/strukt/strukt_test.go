package strukt

import (
	"github.com/rmarken5/lava/internal/file-gen/property"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStruktBuilder_Named(t *testing.T) {

	testCases := map[string]struct {
		wantName string
	}{
		"should print ryan": {
			wantName: "ryan",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sb := StruktBuilder{}
			s := sb.Named("ryan").Build()
			assert.Equal(t, s.Name, tc.wantName)
		})
	}

}

func TestStruktBuilder_AddPropertyBuilder(t *testing.T) {

	testCases := map[string]struct {
		wantPropertyBuilders []property.BuildProperty
	}{
		"should print ryan": {
			wantPropertyBuilders: []property.BuildProperty{myBuilderFunc},
		},
	}

	for name, _ := range testCases {
		t.Run(name, func(t *testing.T) {
			sb := StruktBuilder{}
			s := sb.AddPropertyBuilder(myBuilderFunc).Build()
			pb := s.PropertyBuilders[0]
			builder := &property.PropertyBuilder{}
			p := pb(builder)

			assert.Equal(t, p.Tag, "hello")
		})
	}
}

func myBuilderFunc(b *property.PropertyBuilder) *property.Property {
	return b.Tagged("hello").Build()
}
