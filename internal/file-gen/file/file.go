package file

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/rmarken5/lava/internal/file-gen/strukt"
	"text/template"
)

//go:embed file.gotmpl
var fs embed.FS

type (
	File struct {
		Pkg            string
		Imports        []string
		StruktBuilders []strukt.BuildStrukt
		StruktPrinter  struktPrinter
	}

	FileBuilder struct {
		modifiers []func(f *File)
	}

	struktPrinter func(b strukt.BuildStrukt) ([]byte, error)
)

func (fb *FileBuilder) WithPackage(pkg string) *FileBuilder {
	fb.modifiers = append(fb.modifiers, func(f *File) {
		f.Pkg = pkg
	})
	return fb
}

func (fb *FileBuilder) WithImports(imports []string) *FileBuilder {
	fb.modifiers = append(fb.modifiers, func(f *File) {
		f.Imports = imports
	})
	return fb
}

func (fb *FileBuilder) WithStruktBuilder(fn strukt.BuildStrukt) *FileBuilder {
	fb.modifiers = append(fb.modifiers, func(f *File) {
		f.StruktBuilders = append(f.StruktBuilders, fn)
	})
	return fb
}

func (fb *FileBuilder) WithStruktPrinter(fn struktPrinter) *FileBuilder {
	fb.modifiers = append(fb.modifiers, func(f *File) {
		f.StruktPrinter = fn
	})
	return fb
}

func (fb *FileBuilder) Build() *File {
	file := &File{}
	for _, mod := range fb.modifiers {
		mod(file)
	}
	return file
}

func PrintFile(fn func(fb *FileBuilder) *File) ([]byte, error) {
	builder := &FileBuilder{}
	file := fn(builder)
	return printFile(file)
}

func printFile(file *File) ([]byte, error) {
	f, err := template.ParseFS(fs, "file.gotmpl")
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(nil)

	err = f.Execute(buff, file)
	if err != nil {
		return nil, err
	}

	for _, sb := range file.StruktBuilders {
		b, err := file.StruktPrinter(sb)
		if err != nil {
			fmt.Println(err)
		}
		buff.Write(b)
	}

	return buff.Bytes(), nil
}
