package logic

import (
	"bytes"
	"embed"
	data_access "github.com/rmarken5/lava/inspect/data-access"
	"log"
	"reflect"
	"regexp"
	"text/template"
	"time"
)

//go:embed templates/struct.tmpl
var content embed.FS

var (
	pgTypeToPrimitive = map[string]reflect.Kind{
		"bigint":    reflect.Int64,
		"integer":   reflect.Int,
		"smallint":  reflect.Int16,
		"numeric":   reflect.Float64,
		"real":      reflect.Float32,
		"double":    reflect.Float64,
		"timestamp": reflect.TypeOf(time.Time{}).Kind(),
		"text":      reflect.String,
		"varchar":   reflect.String,
	}
)

const tablesInQueryRegEx = `(?i)\b(?:FROM|JOIN)\s+(\w+)(?:\s+AS\s+\w+)?(?:\s*,\s*(?:\w+)(?:\s+AS\s+\w+)?)*`

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

	columnDefMap := l.buildColumnDefMap(tableDefinitions)
	structBytes, err := l.buildStructStringFromTemplate(columnDefMap)
	if err != nil {
		return "", err
	}

	return string(structBytes), err

}

func (l *LogicImpl) getAllTableDefinitionsFromQuery(query string) ([]data_access.Table, error) {
	tableNames := l.getTablesFromQuery(query)
	tableDefinitions := make([]data_access.Table, len(tableNames))
	for i, name := range tableNames {
		definition, err := l.dataAccess.InspectTable(name)
		if err != nil {
			l.logger.Printf("error in inspecting table: %s ", name)
			return tableDefinitions, err
		}
		tableDefinitions[i] = definition
	}

	return tableDefinitions, nil
}

func (l *LogicImpl) getTablesFromQuery(query string) []string {
	l.logger.Printf("getting tables from query: %s: ", query)

	re := regexp.MustCompile(tablesInQueryRegEx)

	matches := re.FindAllStringSubmatch(query, -1)
	tables := make([]string, len(matches))
	for i, table := range matches {
		tables[i] = table[len(table)-1]
	}
	return tables
}

func (l *LogicImpl) buildColumnDefMap(tables []data_access.Table) map[string][]ColumnDef {
	mappedDefs := make(map[string][]ColumnDef, len(tables))

	for _, table := range tables {
		columnDefs := make([]ColumnDef, len(table.Columns))
		for j, col := range table.Columns {
			def := ColumnDef{
				Column: table.Columns[j],
				Kind:   pgTypeToPrimitive[col.DataType],
			}
			columnDefs[j] = def
		}
		mappedDefs[table.Name] = columnDefs
	}

	return mappedDefs
}

func (l *LogicImpl) buildStructStringFromTemplate(definitions map[string][]ColumnDef) ([]byte, error) {
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

func (l *LogicImpl) WriteStructsToDir(dir string, structBytes []byte) error {
	//TODO implement me
	panic("implement me")
}
