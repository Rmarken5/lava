package query

import (
	"github.com/rmarken5/lava/internal/database"
	"github.com/rmarken5/lava/internal/file-gen/strukt"
	"strings"
)

var reservedWords = []string{"SELECT",
	"FROM",
	"WHERE",
	"GROUP BY",
	"HAVING",
	"ORDER BY",
	"INSERT INTO",
	"VALUES",
	"UPDATE",
	"SET",
	"DELETE",
	"JOIN",
	"INNER JOIN",
	"LEFT JOIN",
	"RIGHT JOIN",
	"FULL OUTER JOIN",
	"UNION",
	"DISTINCT",
	"AS",
	"ON",
	"BETWEEN",
	"LIKE",
	"IN",
	"IS NULL",
	"NOT",
	"AND",
	"OR",
	"ASC",
	"DESC",
	"LIMIT",
	"OFFSET"}

func (p *parser) ParseQuery(query string) ([]*strukt.Strukt, error) {
	strukts := make([]*strukt.Strukt, 0)
	// 1. Find tables in query.
	// 2. Find columns in select clause.
	// 3. Inspect table definitions
	// 4. Map select clause columns to tables
	// 5. Map go data-types to database data-type

	return strukts, nil
}

func (p *parser) findTablesInQuery(query string) []string {
	query = strings.ToUpper(query)
	trimmedQ := strings.TrimPrefix(query, "FROM ")
	for _, w := range reservedWords {
		trimmedQ = strings.TrimSuffix(query, w)
	}
	tables := strings.Split(trimmedQ, ",")

	return tables
}

func (p *parser) findColumnsInSelectClause(query string) ([]string, error) {
	columns := make([]string, 0)

	return columns, nil
}

func (p *parser) inspectTables(tables []string) ([]*database.Table, error) {
	dbDefs := make([]*database.Table, 0)

	return dbDefs, nil
}

func (p *parser) inspectTable(table string) (*database.Table, error) {
	dbDef := &database.Table{}

	return dbDef, nil
}

func (p *parser) mapColumnsToTables(tableDefs []*database.Table, columns []string) error {

	return nil
}
