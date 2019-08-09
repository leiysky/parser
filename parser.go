// Copyright 2019 leiysky
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/leiysky/parser/ast"
	ps "github.com/leiysky/parser/parser"
)

// New will create a Parser
func New() *Parser {
	return &Parser{}
}

// Parser is used for parsing cypher
type Parser struct {
	l *ps.CypherLexer
}

// Parse will parse the cypher into ast.Stmt
func (p *Parser) Parse(cypher string) ast.Stmt {
	p.l = ps.NewCypherLexer(antlr.NewInputStream(cypher))
	tokenStream := antlr.NewCommonTokenStream(p.l, antlr.LexerDefaultTokenChannel)
	parser := ps.NewCypherParser(tokenStream)
	tree := parser.Cypher()
	v := ps.NewConvertVisitor(parser)
	cypherStmt := v.Visit(tree).(ast.Stmt)
	return cypherStmt
}
