package ast

type CypherStmtType byte

const (
	CypherStmtQuery CypherStmtType = iota
	CypherStmtStandaloneCall
)

type CypherStmt struct {
	baseStmt

	Type           CypherStmtType
	Query          *QueryStmt
	StandaloneCall *StandaloneCall
}

func (n *CypherStmt) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*CypherStmt)
	switch n.Type {
	case CypherStmtQuery:
		n.Query.Accept(v)
	case CypherStmtStandaloneCall:
		n.StandaloneCall.Accept(v)
	}
	return v.Leave(n)
}

func (n *CypherStmt) Restore(ctx *RestoreContext) {
	switch n.Type {
	case CypherStmtQuery:
		n.Query.Restore(ctx)
	}
}

type QueryStmt struct {
	baseStmt

	Clauses []Stmt
}

func (n *QueryStmt) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*QueryStmt)
	for _, c := range n.Clauses {
		c.Accept(v)
	}
	return v.Leave(n)
}

func (n *QueryStmt) Restore(ctx *RestoreContext) {
	for i, c := range n.Clauses {
		if i > 0 {
			ctx.Write(" ")
		}
		c.Restore(ctx)
	}
}

type UnionClause struct {
	baseStmt

	All     bool
	Clauses []Stmt
}

func (n *UnionClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*UnionClause)
	for _, c := range n.Clauses {
		c.Accept(v)
	}
	return v.Leave(n)
}

func (n *UnionClause) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("UNION ")
	if n.All {
		ctx.WriteKeyword("ALL ")
	}
	for _, c := range n.Clauses {
		c.Restore(ctx)
	}
}

type StandaloneCall struct {
	baseNode
}

func (n *StandaloneCall) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*StandaloneCall)
	// TODO
	return v.Leave(n)
}

type WithClause struct {
	baseStmt

	Distinct   bool
	ReturnBody *ReturnBody
	Where      Expr
}

func (n *WithClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*WithClause)
	n.ReturnBody.Accept(v)
	n.Where.Accept(v)
	return v.Leave(n)
}

func (n *WithClause) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("WITH ")
	if n.Distinct {
		ctx.WriteKeyword("DISTINCT ")
	}
	n.ReturnBody.Restore(ctx)
	ctx.Write(" ")
	if n.Where != nil {
		n.Where.Restore(ctx)
	}
}

type ReturnClause struct {
	baseStmt

	Distinct   bool
	ReturnBody *ReturnBody
}

func (n *ReturnClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ReturnClause)
	n.ReturnBody.Accept(v)
	return v.Leave(n)
}

func (n *ReturnClause) Restore(ctx *RestoreContext) {
	ctx.Write("RETURN ")
	if n.Distinct {
		ctx.WriteKeyword("DISTINCT ")
	}
	n.ReturnBody.Restore(ctx)
}
