package ast

// Node is basic type of ast.
type Node interface {
	Accept(visitor Visitor) (node Node, ok bool)
	Text() string
	SetText(text string)
	Restore(ctx *RestoreContext)
}

// Stmt represents statement node.
type Stmt interface {
	Node
	statementNode()
}

// Expr represents expression node.
type Expr interface {
	Node
	exprNode()
}

// Visitor can visit ast.
type Visitor interface {
	// Enter would be called at the begin of Accept.
	// The visited node will be replaced by returned Node.
	// If skipChildren is true, the child nodes' Accept wouldn't be called.
	Enter(n Node) (node Node, skipChildren bool)

	// Accept will directly return Leave.
	Leave(n Node) (node Node, ok bool)
}
