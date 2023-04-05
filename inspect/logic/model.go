package logic

import (
	data_access "github.com/rmarken5/lava/inspect/data-access"
)

type (
	ColumnDef struct {
		Column data_access.Column
		Kind   string
	}
)
