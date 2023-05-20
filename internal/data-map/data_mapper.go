package data_map

import "reflect"

var (
	// https://www.postgresql.org/docs/current/datatype.html
	PgTypeToPrimitive = map[string]string{
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
