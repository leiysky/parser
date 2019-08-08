package ast

// ReadingClauseType represents types of ReadingClause
type ReadingClauseType int

// Values of ReadingClauseType
const (
	ReadingClauseMatch ReadingClauseType = iota
	ReadingClauseUnwind
)

// ReadingClause represents Reading clause in cypher
type ReadingClause struct {
	baseStmt

	Type   ReadingClauseType
	Match  *MatchClause
	Unwind *UnwindClause
}

func (n *ReadingClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ReadingClause)
	switch n.Type {
	case ReadingClauseMatch:
		n.Match.Accept(v)
	case ReadingClauseUnwind:
		n.Unwind.Accept(v)
	}
	return v.Leave(n)
}

// MatchClause represents MATCH clause
type MatchClause struct {
	baseStmt

	Pattern  *Pattern
	Optional bool
	Where    Expr
}

func (n *MatchClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*MatchClause)
	n.Pattern.Accept(v)
	if n.Where != nil {
		n.Where.Accept(v)
	}
	return v.Leave(n)
}

type UnwindClause struct {
	baseStmt

	Expr     Expr
	Variable *SymbolicNameNode
}

func (n *UnwindClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*UnwindClause)
	n.Expr.Accept(v)
	if n.Variable != nil {
		n.Variable.Accept(v)
	}
	return v.Leave(n)
}
