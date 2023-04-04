package data_access

type (
	Table struct {
		Name    string
		Columns []Column
	}

	Column struct {
		Name     string
		DataType string
	}
)
