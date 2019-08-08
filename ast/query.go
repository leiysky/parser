package ast

type CypherStmtType int

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

// type SingleQueryStmtType int

// const (
// 	SingleQueryStmtSinglePart SingleQueryStmtType = iota
// 	SingleQueryStmtMultiPart
// )

// type SingleQueryStmt struct {
// 	baseStmt

// 	Type       SingleQueryStmtType
// 	SinglePart *SinglePartQueryStmt
// 	MultiPart  *MultiPartQueryStmt
// }

// func (n *SingleQueryStmt) Accept(v Visitor) (Node, bool) {
// 	newNode, skip := v.Enter(n)
// 	if skip {
// 		return v.Leave(n)
// 	}
// 	n = newNode.(*SingleQueryStmt)
// 	switch n.Type {
// 	case SingleQueryStmtSinglePart:
// 		n.SinglePart.Accept(v)
// 	case SingleQueryStmtMultiPart:
// 		n.MultiPart.Accept(v)
// 	}
// 	return v.Leave(n)
// }

// type SinglePartQueryStmt struct {
// 	baseStmt

// 	ReadingClauses []StmtNode

// 	// If length of UpdatingClauses is greater than 0 then there could be no ReturnClause
// 	UpdatingClauses []StmtNode
// 	Return          *ReturnClause
// }

// func (n *SinglePartQueryStmt) Accept(v Visitor) (Node, bool) {
// 	newNode, skip := v.Enter(n)
// 	if skip {
// 		return v.Leave(n)
// 	}
// 	n = newNode.(*SinglePartQueryStmt)
// 	for _, reading := range n.ReadingClauses {
// 		reading.Accept(v)
// 	}
// 	for _, updating := range n.UpdatingClauses {
// 		updating.Accept(v)
// 	}
// 	n.Return.Accept(v)
// 	return v.Leave(n)
// }

// type MultiPartQueryStmt struct {
// 	baseStmt

// 	MultiPart  []*MultiPartQueryPartial
// 	SinglePart *SinglePartQueryStmt
// }

// func (n *MultiPartQueryStmt) Accept(v Visitor) (Node, bool) {
// 	newNode, skip := v.Enter(n)
// 	if skip {
// 		return v.Leave(n)
// 	}
// 	n = newNode.(*MultiPartQueryStmt)
// 	for _, part := range n.MultiPart {
// 		for _, reading := range part.Readings {
// 			reading.Accept(v)
// 		}
// 		for _, updating := range part.Updatings {
// 			updating.Accept(v)
// 		}
// 		part.With.Accept(v)
// 	}
// 	n.SinglePart.Accept(v)
// 	return v.Leave(n)
// }

// type MultiPartQueryPartial struct {
// 	Readings  []*ReadingClause
// 	Updatings []*UpdatingClause
// 	With      *WithClause
// }

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
