package logic

import (
	data_access "github.com/rmarken5/lava/inspect/data-access"
	"strings"
)

type (
	ColumnDef struct {
		Column data_access.Column
		Kind   string
	}
)

func (c ColumnDef) VariableName() string {
	parts := strings.Split(c.Column.Name, "_")
	variable := ""
	for _, part := range parts {
		variable += strings.ToUpper(part[:1]) + part[1:]
	}
	return variable
}
