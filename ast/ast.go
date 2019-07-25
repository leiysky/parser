package ast

type Node interface {
	Accept(visitor Visitor) (node Node, ok bool)
	Text() string
	SetText(text string)
}

type StmtNode interface {
	Node
	statement()
}

type ExprNode interface {
	Node
}

type Visitor interface {
	Enter(n Node) (node Node, skipChildren bool)

	Leave(n Node) (node Node, ok bool)
}
