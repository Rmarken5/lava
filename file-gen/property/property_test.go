package property

import (
	"bytes"
	"testing"
)

func TestPrintProperties(t *testing.T) {
	tests := []struct {
		name     string
		builder  BuildProperty
		expected string
	}{
		{
			name: "no modifiers",
			builder: func(b *PropertyBuilder) *Property {
				return b.Build()
			},
			expected: "",
		},
		{
			name: "with name and type",
			builder: func(b *PropertyBuilder) *Property {
				return b.Named("Foo").OfType("string").Build()
			},
			expected: "Foo string",
		},
		{
			name: "with tag",
			builder: func(b *PropertyBuilder) *Property {
				return b.Tagged("json").Build()
			},
			expected: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PrintProperties(tt.builder)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(result, []byte(tt.expected)) {
				t.Errorf("expected: %q, got: %q", tt.expected, string(result))
			}
		})
	}
}
