package database

type (
	Table struct {
		Name    string
		Columns Columns
	}

	Column struct {
		Name     string
		DataType string
	}
	Columns []Column
)

// FindFirst finds first occurrence of column
func (cs Columns) FindFirst(name string) *Column {
	for _, col := range cs {
		if col.Name == name {
			return &col
		}
	}
	return nil
}
