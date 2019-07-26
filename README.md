# Cypher Parser

A CypherQL parser with [openCypher](https://www.opencypher.org) EBNF grammar.

# Status
Working in progress.

# Usage

Install: 

```shell
$ go get -u github.com/leiysky/parser
```

Print the ast node types with visitor: 

```go
import (
  "fmt"
  "github.com/leiysky/parser"
)

type printVisitor struct {
	ast.Visitor
}

func (v *printVisitor) Enter(node ast.Node) (ast.Node, bool) {
	fmt.Printf("%T\n", node)
	return node, false
}

func (v *printVisitor) Leave(node ast.Node) (ast.Node, bool) {
	return node, true
}

func main() {
  cypher := `
MATCH (n)-[r:Label *1..10]->()
WHERE n.name = 'leiysky'
SET n.name = 'leiysky1'
RETURN n, r
`
  p := parser.New()
  stmt := p.Parse(cypher)
  stmt.Accept(&testVisitor{})
}
```