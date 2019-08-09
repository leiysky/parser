package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/leiysky/parser/ast"
	ps "github.com/leiysky/parser/parser"
)

func New() *Parser {
	return &Parser{}
}

type Parser struct {
	l *ps.CypherLexer
}

func (p *Parser) Parse(cypher string) ast.Stmt {
	p.l = ps.NewCypherLexer(antlr.NewInputStream(cypher))
	tokenStream := antlr.NewCommonTokenStream(p.l, antlr.LexerDefaultTokenChannel)
	parser := ps.NewCypherParser(tokenStream)
	tree := parser.Cypher()
	v := ps.NewConvertVisitor(parser)
	cypherStmt := v.Visit(tree).(ast.Stmt)
	return cypherStmt
}
