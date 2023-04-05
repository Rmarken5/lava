package data_access

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type (
	Inspector interface {
		InspectTable(name string) (Table, error)
	}

	InspectorImpl struct {
		db     *sql.DB
		logger *log.Logger
	}
)

func NewInspector(db *sql.DB, logger *log.Logger) *InspectorImpl {
	return &InspectorImpl{
		db:     db,
		logger: logger,
	}
}

func (i *InspectorImpl) InspectTable(name string) (Table, error) {
	columns := make([]Column, 0)
	table := Table{
		Name: name,
	}

	i.logger.Printf("introspecting columns on table %s\n", name)
	rows, err := i.db.Query(fmt.Sprintf("SELECT column_name, data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '%s' order by ordinal_position", strings.ToLower(name)))
	if err != nil {
		i.logger.Println(err)
		return table, err
	}
	for rows.Next() {
		var col Column
		err := rows.Scan(&col.Name, &col.DataType)
		col.DataType = strings.Fields(col.DataType)[0]
		if err != nil {
			return table, err
		}
		columns = append(columns, col)
	}

	table.Columns = columns

	return table, nil
}
