package ast

type CypherStmtType int

const (
	CypherStmtQuery CypherStmtType = iota
	CypherStmtStandaloneCall
)

type CypherStmt struct {
	baseStmt

	Query          *RegularQueryStmt
	StandaloneCall *StandaloneCall
	Type           CypherStmtType
}

type RegularQueryStmt struct {
	baseStmt

	SingleQuery *SingleQueryStmt
	Unions      []*UnionClause
}

type UnionClause struct {
	baseStmt

	All         bool
	SingleQuery *SingleQueryStmt
}

type StandaloneCall struct {
	baseNode
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

type SinglePartQueryStmt struct {
	baseStmt

	ReadingClauses []*ReadingClause

	// If length of UpdatingClauses is greater than 0 then there could be no ReturnClause
	UpdatingClauses []*UpdatingClause
	Return          *ReturnClause
}

type MultiPartQueryStmt struct {
	baseStmt
}

type ReturnClause struct {
	baseStmt

	Distinct   bool
	ReturnBody *ReturnBody
}
