package logic

import (
	data_access "github.com/rmarken5/lava/inspect/data-access"
	"strings"
)

type (
	TableDef struct {
		Table                  data_access.Table
		Alias                  string
		EmbeddedStructTagValue string
		ColumnDefs             []ColumnDef
	}
	ColumnDef struct {
		Column data_access.Column
		Alias  string
		Kind   string
	}
)

func (t TableDef) TableName() string {
	parts := strings.Split(t.Table.Name, "_")
	variable := ""
	for _, part := range parts {
		variable += strings.ToUpper(part[:1]) + part[1:]
	}
	return variable
}

func (c ColumnDef) VariableName() string {
	parts := strings.Split(c.Column.Name, "_")
	variable := ""
	for _, part := range parts {
		variable += strings.ToUpper(part[:1]) + part[1:]
	}
	return variable
}

//func (c ColumnDef) TagName() string {
//	if c.Alias != "" {
//		return c.Alias
//	}
//	return c.Column.Name
//}
