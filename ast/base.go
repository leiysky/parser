package ast

type baseNode struct {
	Node
	text string
}

func (n *baseNode) Text() string {
	return n.text
}

func (n *baseNode) SetText(text string) {
	n.text = text
}

type baseStmt struct {
	baseNode
}

func (n *baseStmt) statementNode() {}

type baseExpr struct {
	baseNode
}

func (n *baseExpr) exprNode() {}
