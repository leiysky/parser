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
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/leiysky/parser/ast"
)

type testCase struct {
	original string
	pass     bool
	target   string
}

var cases = []testCase{
	{"match (n) return n", true, "MATCH (`n`) RETURN `n`"},
	{"match (n:Label)-[r:Type *1..2 {hello:'world'}]->() return *", true, "MATCH (`n`:Label)-[`r`:Type*1..2{hello: 'world'}]->() RETURN *"},
	{"with n as n, a as a create (n)-[]-(a)", true, "WITH `n` AS `n`, `a` AS `a` CREATE (`n`)-[*1..1]-(`a`)"},
	{"match (n) where n.name > 1 AND 'abc' = 2 OR 1.2 <> -2 return n", true, "MATCH (`n`) WHERE `n`.`name` > 1 AND 'abc' = 2 OR 1.200000 <> -2 RETURN `n`"},
	{"match (n) return count(*)", true, "MATCH (`n`) RETURN COUNT(*)"},
	{"match (n) return [n in list | n+1]", true, "MATCH (`n`) RETURN [`n` IN `list` | `n` + 1]"},
	{"match (n) return any(n in list), all(n in list), single(n in list), none(n in list where TRUE)", true, "MATCH (`n`) RETURN ANY(`n` IN `list`), ALL(`n` IN `list`), SINGLE(`n` IN `list`), NONE(`n` IN `list` WHERE TRUE)"},
}

func runTestCase(t *testing.T, cases []testCase) {
	parser := New()
	var target strings.Builder
	for _, c := range cases {
		target.Reset()
		originalAst := parser.Parse(c.original)
		ctx := ast.NewRestoreContext(&target)
		originalAst.Restore(ctx)
		if target.String() != c.target {
			t.Fatalf("obtained: %s; expected: %s", target.String(), c.target)
		}
	}
}

func TestCases(t *testing.T) {
	runTestCase(t, cases)
}

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
	t.Skip()
	parser := New()
	parser.Parse(`
MATCH (n)-[:Label *0123]-()
WHERE n.id = n +1
SET n = {}
RETURN n`)
}

func TestMultiPartQuery(t *testing.T) {
	t.Skip()
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
	t.Skip()
	parser := New()
	stmt := parser.Parse(`
	MATCH (n:Label1:Label2)-[r:Type1|Type2*1..2{name:'hello'}]->(n1)
	WHERE n.a OR n AND n XOR NOT n < n > n = n <= n >= n <> n + n - n * n / n % -n
	MATCH (n)
	WHERE 1 = 1 < 1
	RETURN n
	`)
	stmt.Accept(&testVisitor{})
	rst := ast.NewRestoreContext(os.Stdout)
	stmt.Restore(rst)
}
