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

func (n *baseStmt) statement() {}

type Expr struct {
	baseNode
}

func (n *Expr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*Expr)
	// TODO
	return v.Leave(n)
}
