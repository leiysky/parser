package parser

import (
	"fmt"
	"testing"

	"github.com/leiysky/parser/ast"
)

type testVisitor struct {
	ast.Visitor
	depth int
}

func (v *testVisitor) Enter(node ast.Node) (ast.Node, bool) {
	for i := 0; i < v.depth; i++ {
		fmt.Print("\033[32m| \033[0m")
	}
	if expr, ok := node.(ast.Expr); ok {
		v.VisitExpr(expr)
	} else {
		fmt.Printf("\033[33m%T\n\033[0m", node)
	}
	v.depth++

	return node, false
}

func (v *testVisitor) Leave(node ast.Node) (ast.Node, bool) {
	v.depth--
	return node, true
}

func (v *testVisitor) VisitExpr(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		fmt.Printf("\033[33m%T \033[0m \033[31m%s\n\033[0m", e, e.Op)
	case *ast.UnaryExpr:
		fmt.Printf("\033[33m%T \033[0m \033[31m%s\n\033[0m", e, e.Op)
	default:
		fmt.Printf("\033[33m%T\n\033[0m", e)
	}
}

func TestParser(t *testing.T) {
	parser := New()
	parser.Parse(`
MATCH (n)-[:Label *0123]-()
WHERE n.id = n +1
SET n = {}
RETURN n`)
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

func TestExpression(t *testing.T) {
	parser := New()
	stmt := parser.Parse(`
	MATCH (n)
	WHERE n OR n AND n XOR NOT n < n > n = n <= n >= n <> n + n - n * n / n % -n
	MATCH (n)
	WHERE 1 = 1 < 1
	RETURN n
	`)
	stmt.Accept(&testVisitor{})
}
