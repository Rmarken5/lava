package query

import (
	"github.com/rmarken5/lava/internal/database"
	"github.com/rmarken5/lava/internal/file-gen/strukt"
	"log"
)

type (
	Parser interface {
		ParseQuery(query string) ([]*strukt.Strukt, error)
	}
	parser struct {
		query       string
		logger      *log.Logger
		dbInspector database.Inspector
	}

	ParserBuilder struct {
		modifiers []parserModifier
	}
	parserModifier func(parser *parser)
)

func (b *ParserBuilder) WithQuery(q string) *ParserBuilder {
	b.modifiers = append(b.modifiers, func(parser *parser) {
		parser.query = q
	})
	return b
}

func (b *ParserBuilder) WithLogger(l *log.Logger) *ParserBuilder {
	b.modifiers = append(b.modifiers, func(parser *parser) {
		parser.logger = l
	})
	return b
}
func (b *ParserBuilder) WithInspector(i database.Inspector) *ParserBuilder {
	b.modifiers = append(b.modifiers, func(parser *parser) {
		parser.dbInspector = i
	})
	return b
}
