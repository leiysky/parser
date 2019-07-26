package ast

type ReturnBody struct {
	baseStmt
}

func (n *ReturnBody) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ReturnBody)
	// TODO
	return v.Leave(n)
}
