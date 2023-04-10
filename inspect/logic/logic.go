package logic

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/rmarken5/lava/inspect/data-access"
	"log"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

//go:embed templates/struct.tmpl
var content embed.FS

var (
	// https://www.postgresql.org/docs/current/datatype.html
	pgTypeToPrimitive = map[string]string{
		"bigint":      reflect.Int64.String(),
		"int8":        reflect.Int64.String(),
		"bigserial":   reflect.Uint64.String(),
		"integer":     reflect.Int.String(),
		"int4":        reflect.Int.String(),
		"serial":      reflect.Uint32.String(),
		"serial4":     reflect.Uint32.String(),
		"int":         reflect.Int.String(),
		"smallint":    reflect.Int16.String(),
		"int2":        reflect.Int16.String(),
		"smallserial": reflect.Uint16.String(),
		"serial2":     reflect.Uint16.String(),
		"numeric":     reflect.Float64.String(),
		"real":        reflect.Float32.String(),
		"float4":      reflect.Float32.String(),
		"double":      reflect.Float64.String(),
		"float8":      reflect.Float64.String(),
		"timestamp":   "time.Time",
		"text":        reflect.String.String(),
		"varchar":     reflect.String.String(),
		"character":   reflect.String.String(),
		"char":        reflect.String.String(),
		"boolean":     reflect.Bool.String(),
		"bool":        reflect.Bool.String(),
	}
)

const (
	tablesInQueryRegEx        = `(?im)(?:FROM|JOIN)\s+([^\s,;()]+)(?:\s+AS\s+([^\s,;()]+))?`
	findTableAliasRegEx       = `(?im)\sAS|\s(\w+)`
	findColumnAliasRegEx      = `(?im)AS\s(\S+)`
	findAliasedTableNameRegEx = `(?im)(\w+)\sAS\s`
	allColumnsInQueryRegEx    = `(?im)(?:select\s|\s+)(.*?)(?:\s*(?:,|from))`
)

type (
	Logic interface {
		BuildStructsForQuery(query string) (string, error)
		WriteStructsToDir(dir string, structBytes []byte) error
	}

	LogicImpl struct {
		logger     *log.Logger
		dataAccess data_access.Inspector
	}
)

func New(logger *log.Logger, dataAccess data_access.Inspector) *LogicImpl {
	return &LogicImpl{
		logger:     logger,
		dataAccess: dataAccess,
	}
}

func (l *LogicImpl) BuildStructsForQuery(query string) (string, error) {

	rawColumnStrings := getColumnsInSelectClause(query)
	fmt.Println(rawColumnStrings)
	columnsMappedToTableAliases := tableColumnAliasMap(rawColumnStrings)
	fmt.Println(columnsMappedToTableAliases)
	tableDefinitions, err := l.getAllTableDefinitionsFromQuery(query)
	if err != nil {
		return "", err
	}
	l.buildColumnDefMap(tableDefinitions, columnsMappedToTableAliases)
	getEmbeddedStructTagsFromColumns(tableDefinitions)
	structBytes, err := l.buildStructStringFromTemplate(tableDefinitions)
	if err != nil {
		return "", err
	}

	return string(structBytes), err

}

func getEmbeddedStructTagsFromColumns(tables []TableDef) map[string]string {
	tagMapping := make(map[string]string)
	for i := 0; i < len(tables); i++ {
		table := &tables[i]
		table.EmbeddedStructTagValue = table.Table.Name
		if table.Alias != "" {
			table.EmbeddedStructTagValue = table.Alias
		}
		for _, columnDef := range table.ColumnDefs {
			if strings.Contains(columnDef.Alias, ".") {
				tableAlias := columnDef.Alias[:strings.Index(columnDef.Alias, ".")]
				table.EmbeddedStructTagValue = tableAlias
				break
			}
		}

	}

	return tagMapping
}

func (l *LogicImpl) getAllTableDefinitionsFromQuery(query string) ([]TableDef, error) {
	tableNames := l.getTablesFromQuery(query)
	tableDefinitions := make([]TableDef, len(tableNames))
	for i, table := range tableNames {
		var tableName = table
		alias := findTableAlias(table)
		if alias != "" {
			tableName = findAliasedTableName(table)
		}
		definition, err := l.dataAccess.InspectTable(tableName)
		if err != nil {
			l.logger.Printf("error in inspecting table: %s ", tableName)
			return tableDefinitions, err
		}
		tableDefinitions[i] = TableDef{
			Table: definition,
			Alias: alias,
		}
	}

	return tableDefinitions, nil
}

func (l *LogicImpl) getTablesFromQuery(query string) []string {
	tables := make([]string, 0)

	l.logger.Printf("getting tables from query: %s: ", query)

	re := regexp.MustCompile(tablesInQueryRegEx)

	matches := re.FindAllString(query, -1)
	for _, match := range matches {
		tableNames := regexp.MustCompile(`\w+`).FindAllString(match, -1)
		tables = append(tables, strings.Join(tableNames[1:], " ")) // skip the first match, which is "FROM"
	}
	return tables
}

//func (l *LogicImpl) buildColumnDefMap(tables []TableDef) {
//	for i, table := range tables {
//		columnDefs := make([]ColumnDef, len(table.Table.Columns))
//		for j, col := range table.Table.Columns {
//			alias := findTableAlias(col.Name)
//			def := ColumnDef{
//				Column: table.Table.Columns[j],
//				Kind:   pgTypeToPrimitive[col.DataType],
//				Alias:  alias,
//			}
//			columnDefs[j] = def
//		}
//		tables[i].ColumnDefs = columnDefs
//	}
//}

func (l *LogicImpl) buildColumnDefMap(tables []TableDef, columnsFromSelect map[string][]string) {
	for i := 0; i < len(tables); i++ {

		columnsFromAlias := columnsFromSelect[tables[i].Alias]
		if columnsFromAlias[0] == "*" {
			columnDefsForAllColumns(&tables[i])
		} else {
			columnDefsForAliases(&tables[i], columnsFromAlias)
		}
	}
}

func columnDefsForAliases(table *TableDef, columnsFromAlias []string) {

	columnDefs := make([]ColumnDef, len(columnsFromAlias))
	for j, cols := range columnsFromAlias {
		alias := findColumnAlias(cols)
		colName := findAliasedTableName(cols)
		matchedCol := table.Table.Columns.FindFirst(colName)
		def := ColumnDef{
			Column: *matchedCol,
			Kind:   pgTypeToPrimitive[matchedCol.DataType],
			Alias:  alias,
		}
		columnDefs[j] = def
	}
	table.ColumnDefs = columnDefs
}

func columnDefsForAllColumns(table *TableDef) {
	columnDefs := make([]ColumnDef, len(table.Table.Columns))
	for j, col := range table.Table.Columns {
		alias := findTableAlias(col.Name)
		def := ColumnDef{
			Column: table.Table.Columns[j],
			Kind:   pgTypeToPrimitive[col.DataType],
			Alias:  alias,
		}
		columnDefs[j] = def
	}
	table.ColumnDefs = columnDefs
}

func (l *LogicImpl) buildStructStringFromTemplate(definitions []TableDef) ([]byte, error) {
	tmpl, err := template.ParseFS(content, "templates/struct.tmpl")
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(nil)

	err = tmpl.Execute(buff, definitions)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// findTableAlias finds the table's alias (if it has one).
func findTableAlias(tableName string) string {
	var alias string
	aliasMatch := regexp.MustCompile(findTableAliasRegEx).FindAllStringSubmatch(tableName, -1)
	// if there is a match, getting the last entry of the inner array gives you the alias
	if len(aliasMatch) != 0 {
		alias = aliasMatch[len(aliasMatch)-1][len(aliasMatch[0])-1]
	}
	return alias
}

// findAliasedTableName finds the table name which contains an alias
func findAliasedTableName(tableName string) string {
	tableNameMatch := regexp.MustCompile(findAliasedTableNameRegEx).FindAllStringSubmatch(tableName, -1)
	if len(tableNameMatch) != 0 {
		return strings.TrimSuffix(tableNameMatch[0][len(tableNameMatch[0])-1], " ")
	}

	return ""
} // findColumnAlias finds the alias of a column
func findColumnAlias(columnName string) string {
	columnNameMatch := regexp.MustCompile(findColumnAliasRegEx).FindAllStringSubmatch(columnName, -1)
	if len(columnNameMatch) != 0 {
		return strings.TrimSuffix(columnNameMatch[0][len(columnNameMatch[0])-1], " ")
	}

	return ""
}

func getColumnsInSelectClause(query string) []string {
	var cols []string
	re := regexp.MustCompile(allColumnsInQueryRegEx)
	match := re.FindAllStringSubmatch(query, -1)
	if match != nil {
		for _, m := range match {
			cols = append(cols, m[1:]...)
		}
	}
	return cols
}

func tableColumnAliasMap(columnStrings []string) map[string][]string {
	colMap := make(map[string][]string)
	for _, col := range columnStrings {
		idx := strings.Index(col, ".")
		tableAlias := col[:idx]
		columnName := col[idx+1:]
		columnName = strings.ReplaceAll(columnName, "\"", "")
		if _, ok := colMap[tableAlias]; !ok {
			colMap[tableAlias] = []string{columnName}
		} else {
			colMap[tableAlias] = append(colMap[tableAlias], columnName)
		}
	}
	return colMap
}

func (l *LogicImpl) WriteStructsToDir(dir string, structBytes []byte) error {
	//TODO implement me
	panic("implement me")
}
