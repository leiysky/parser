package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/leiysky/parser/ast"
)

func New() *Parser {
	return &Parser{}
}

type Parser struct {
	l *CypherLexer
}

func (p *Parser) Parse(cypher string) ast.Stmt {
	p.l = NewCypherLexer(antlr.NewInputStream(cypher))
	tokenStream := antlr.NewCommonTokenStream(p.l, antlr.LexerDefaultTokenChannel)
	parser := NewCypherParser(tokenStream)
	tree := parser.Cypher()
	v := &convertVisitor{parser}
	cypherStmt := v.Visit(tree).(ast.Stmt)
	return cypherStmt
}
