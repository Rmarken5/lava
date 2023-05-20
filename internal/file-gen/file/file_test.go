package file

import (
	"github.com/rmarken5/lava/internal/file-gen/strukt"
	"testing"
)

func TestPrintFile(t *testing.T) {
	struktBuilder1 := func(b *strukt.StruktBuilder) *strukt.Strukt { return nil }
	struktBuilder2 := func(b *strukt.StruktBuilder) *strukt.Strukt { return nil }

	tests := []struct {
		name           string
		fileBuilder    func(fb *FileBuilder) *File
		expectedOutput string
		expectedError  bool
	}{
		{
			name: "basic file",
			fileBuilder: func(fb *FileBuilder) *File {
				return fb.
					WithPackage("main").
					WithImports([]string{"fmt"}).
					WithStruktBuilder(struktBuilder1).
					WithStruktPrinter(func(b strukt.BuildStrukt) ([]byte, error) {
						return []byte{}, nil
					}).
					Build()
			},
			expectedOutput: `package main
        
        
        imports (
            "fmt"
        )`,
			expectedError: false,
		},
		{
			name: "file with struct builders",
			fileBuilder: func(fb *FileBuilder) *File {
				return fb.
					WithPackage("main").
					WithImports([]string{"fmt"}).
					WithStruktBuilder(struktBuilder1).
					WithStruktBuilder(struktBuilder2).
					WithStruktPrinter(func(b strukt.BuildStrukt) ([]byte, error) {
						return []byte{}, nil
					}).
					Build()
			},
			expectedOutput: `package main
        
        
        imports (
            "fmt"
        )`,
			expectedError: false,
		},
		{
			name: "file with invalid template",
			fileBuilder: func(fb *FileBuilder) *File {
				return fb.
					WithPackage("main").
					WithImports([]string{"fmt"}).
					WithStruktBuilder(struktBuilder1).
					WithStruktPrinter(func(b strukt.BuildStrukt) ([]byte, error) {
						return []byte{}, nil
					}).
					Build()
			},
			expectedOutput: "",
			expectedError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := PrintFile(test.fileBuilder)

			if test.expectedError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !test.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if string(output) != test.expectedOutput {
				t.Errorf("expected\n%s\nbut got\n%s", test.expectedOutput, string(output))
			}
		})
	}
}
