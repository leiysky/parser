package parser

import (
	"fmt"
	"testing"

	"github.com/leiysky/parser/ast"
)

type testVisitor struct {
	ast.Visitor
}

func (v *testVisitor) Enter(node ast.Node) (ast.Node, bool) {
	fmt.Printf("%T\n", node)
	return node, false
}

func (v *testVisitor) Leave(node ast.Node) (ast.Node, bool) {
	return node, true
}

func TestParser(t *testing.T) {
	parser := New()
	parser.Parse(`
MATCH (n)-[:Label *0123]-()
WHERE n.id = n +1
SET n = {}
RETURN n`)
	// stmt.Accept(&testVisitor{})
}

func TestMultiPartQuery(t *testing.T) {
	parser := New()
	parser.Parse(`
MATCH (n)
SET n.id = 1
WITH n
MATCH (n)
RETURN n
`)
}
