package logic

import (
	"bytes"
	"embed"
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
	findAliasRegEx            = `(?im)\sAS|\s(\w+)`
	findAliasedTableNameRegEx = `(?im)(\w+)\sAS\s`
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

	tableDefinitions, err := l.getAllTableDefinitionsFromQuery(query)
	if err != nil {
		return "", err
	}
	l.buildColumnDefMap(tableDefinitions)
	structBytes, err := l.buildStructStringFromTemplate(tableDefinitions)
	if err != nil {
		return "", err
	}

	return string(structBytes), err

}

func (l *LogicImpl) getAllTableDefinitionsFromQuery(query string) ([]TableDef, error) {
	tableNames := l.getTablesFromQuery(query)
	tableDefinitions := make([]TableDef, len(tableNames))
	for i, table := range tableNames {
		var tableName = table
		alias := findAlias(table)
		if alias != "" {
			tableName = findAliasedTableName(table)
		}
		definition, err := l.dataAccess.InspectTable(tableName)
		if err != nil {
			l.logger.Printf("error in inspecting table: %s ", tableName)
			return tableDefinitions, err
		}
		tableDefinitions[i] = TableDef{
			Table:      definition,
			Alias:      alias,
			ColumnDefs: nil,
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

func (l *LogicImpl) buildColumnDefMap(tables []TableDef) {

	for i, table := range tables {
		columnDefs := make([]ColumnDef, len(table.Table.Columns))
		for j, col := range table.Table.Columns {
			def := ColumnDef{
				Column: table.Table.Columns[j],
				Kind:   pgTypeToPrimitive[col.DataType],
			}
			columnDefs[j] = def
		}
		tables[i].ColumnDefs = columnDefs
	}
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

func findAlias(tableName string) string {
	var alias string
	aliasMatch := regexp.MustCompile(findAliasRegEx).FindAllStringSubmatch(tableName, -1)
	// if there is a match, getting the last entry of the inner array gives you the alias
	if len(aliasMatch) != 0 {
		alias = aliasMatch[len(aliasMatch)-1][len(aliasMatch[0])-1]

	}
	return alias
}

func findAliasedTableName(tableName string) string {
	tableNameMatch := regexp.MustCompile(findAliasedTableNameRegEx).FindAllStringSubmatch(tableName, -1)
	if len(tableNameMatch) != 0 {
		return strings.TrimSuffix(tableNameMatch[0][len(tableNameMatch[0])-1], " ")
	}

	return ""
}

func (l *LogicImpl) WriteStructsToDir(dir string, structBytes []byte) error {
	//TODO implement me
	panic("implement me")
}
