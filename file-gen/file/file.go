package file

import "github.com/rmarken5/lava/file-gen/strukt"

type File struct {
	Pkg     string
	Imports []string
	Structs []strukt.Strukt
}
