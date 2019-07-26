package ast

type CypherStmtType int

const (
	CypherStmtQuery CypherStmtType = iota
	CypherStmtStandaloneCall
)

type CypherStmt struct {
	baseStmt

	Type           CypherStmtType
	Query          *RegularQueryStmt
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

type RegularQueryStmt struct {
	baseStmt

	SingleQuery *SingleQueryStmt
	Unions      []*UnionClause
}

func (n *RegularQueryStmt) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*RegularQueryStmt)
	n.SingleQuery.Accept(v)
	for _, union := range n.Unions {
		union.Accept(v)
	}
	return v.Leave(n)
}

type UnionClause struct {
	baseStmt

	All         bool
	SingleQuery *SingleQueryStmt
}

func (n *UnionClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*UnionClause)
	n.SingleQuery.Accept(v)
	return v.Leave(n)
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

type SingleQueryStmtType int

const (
	SingleQueryStmtSinglePart SingleQueryStmtType = iota
	SingleQueryStmtMultiPart
)

type SingleQueryStmt struct {
	baseStmt

	Type       SingleQueryStmtType
	SinglePart *SinglePartQueryStmt
	MultiPart  *MultiPartQueryStmt
}

func (n *SingleQueryStmt) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*SingleQueryStmt)
	switch n.Type {
	case SingleQueryStmtSinglePart:
		n.SinglePart.Accept(v)
	case SingleQueryStmtMultiPart:
		n.MultiPart.Accept(v)
	}
	return v.Leave(n)
}

type SinglePartQueryStmt struct {
	baseStmt

	ReadingClauses []*ReadingClause

	// If length of UpdatingClauses is greater than 0 then there could be no ReturnClause
	UpdatingClauses []*UpdatingClause
	Return          *ReturnClause
}

func (n *SinglePartQueryStmt) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*SinglePartQueryStmt)
	for _, reading := range n.ReadingClauses {
		reading.Accept(v)
	}
	for _, updating := range n.UpdatingClauses {
		updating.Accept(v)
	}
	n.Return.Accept(v)
	return v.Leave(n)
}

type MultiPartQueryStmt struct {
	baseStmt
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
