package ast

type ReturnBody struct {
	baseStmt

	ReturnItems []*ReturnItem
	OrderBy     *OrderClause
	Skip        Expr
	Limit       Expr
}

func (n *ReturnBody) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ReturnBody)
	for _, item := range n.ReturnItems {
		item.Accept(v)
	}
	if n.OrderBy != nil {
		n.OrderBy.Accept(v)
	}
	if n.Skip != nil {
		n.Skip.Accept(v)
	}
	if n.Limit != nil {
		n.Limit.Accept(v)
	}
	return v.Leave(n)
}

type ReturnItem struct {
	baseStmt

	Expr     Expr
	As       bool
	Variable *VariableNode
}

func (n *ReturnItem) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ReturnItem)
	n.Expr.Accept(v)
	if n.As {
		n.Variable.Accept(v)
	}
	return v.Leave(n)
}

type OrderClause struct {
	baseStmt

	SortItems []*SortItem
}

func (n *OrderClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*OrderClause)
	for _, item := range n.SortItems {
		item.Accept(v)
	}
	return v.Leave(n)
}

type SortType byte

const (
	SortAscending SortType = iota
	SortDescending
)

type SortItem struct {
	baseStmt

	Type SortType
	Expr Expr
}

func (n *SortItem) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*SortItem)
	n.Expr.Accept(v)
	return v.Leave(n)
}
