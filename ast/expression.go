package ast

type BinaryExpr struct {
	baseExpr

	Op OpType
	L  Expr
	R  Expr
}

type UnaryExpr struct {
	baseExpr

	Op OpType
	E  Expr
}

type ParenExpr struct {
	baseExpr

	InnerExpr Expr
}

type PropertyExpr struct {
	baseExpr

	Atom    *Atom
	Lookups []*PropertyLookup
}

func (n *PropertyExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*PropertyExpr)
	n.Atom.Accept(v)
	for _, lookup := range n.Lookups {
		lookup.Accept(v)
	}
	return v.Leave(n)
}

type OrExpr struct {
	baseExpr

	XorExprs []*XorExpr
}

func (n *OrExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*OrExpr)
	for _, expr := range n.XorExprs {
		expr.Accept(v)
	}
	return v.Leave(n)
}

type XorExpr struct {
	baseExpr

	AndExprs []*AndExpr
}

func (n *XorExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*XorExpr)
	for _, expr := range n.AndExprs {
		expr.Accept(v)
	}
	return v.Leave(n)
}

type AndExpr struct {
	baseExpr

	NotExprs []*NotExpr
}

func (n *AndExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*AndExpr)
	for _, expr := range n.NotExprs {
		expr.Accept(v)
	}
	return v.Leave(n)
}

type NotExpr struct {
	baseExpr

	ComparisonExpr *ComparisonExpr
}

func (n *NotExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*NotExpr)
	n.ComparisonExpr.Accept(v)
	return v.Leave(n)
}

type ComparisonExpr struct {
	baseExpr

	AddSubExpr             *AddSubExpr
	PartialComparisonExprs []*PartialComparisonExpr
}

func (n *ComparisonExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ComparisonExpr)
	n.AddSubExpr.Accept(v)
	for _, partial := range n.PartialComparisonExprs {
		partial.Accept(v)
	}
	return v.Leave(n)
}

type PartialComparisonType int

const (
	PartialComparisonEQ PartialComparisonType = iota
	PartialComparisonNE
	PartialComparisonLT
	PartialComparisonGT
	PartialComparisonLTE
	PartialComparisonGTE
)

type PartialComparisonExpr struct {
	baseExpr

	Type       PartialComparisonType
	AddSubExpr *AddSubExpr
}

func (n *PartialComparisonExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*PartialComparisonExpr)
	n.AddSubExpr.Accept(v)
	return v.Leave(n)
}

type OpType int

const (
	OpAdd OpType = iota
	OpSub
	OpMul
	OpDiv
	OpMod
)

type AddSubExpr struct {
	baseExpr

	// An AddSubExpr starts with a LExpr, and has arbitrary number of RExprs
	LExpr *MulDivModExpr
	// An Op is combined with a RExpr, in AddSubExpr there are only OpAdd or OpSub Ops
	// ASSERT: len(Ops) == len(RExprs) && Op[n] -> RExprs[n]
	Ops    []OpType
	RExprs []*MulDivModExpr
}

func (n *AddSubExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*AddSubExpr)
	n.LExpr.Accept(v)
	for _, expr := range n.RExprs {
		expr.Accept(v)
	}
	return v.Leave(n)
}

type MulDivModExpr struct {
	baseExpr

	// An MulDivModExpr starts with a LExpr, and has arbitrary number of RExprs
	LExpr *PowerOfExpr
	// An Op is combined with a RExpr, in MulDivModExpr there are OpMul ,OpDiv or OpMod Ops
	// ASSERT: len(Ops) == len(RExprs) && Op[n] -> RExprs[n]
	Ops    []OpType
	RExprs []*PowerOfExpr
}

func (n *MulDivModExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*MulDivModExpr)
	for _, expr := range n.RExprs {
		expr.Accept(v)
	}
	return v.Leave(n)
}

type PowerOfExpr struct {
	baseExpr

	UnaryAddSubExprs []*UnaryAddSubExpr
}

func (n *PowerOfExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*PowerOfExpr)
	for _, expr := range n.UnaryAddSubExprs {
		expr.Accept(v)
	}
	return v.Leave(n)
}

type UnaryAddSubExpr struct {
	baseExpr

	IsNegative bool
	SLN        *StringListNullExpr
}

func (n *UnaryAddSubExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*UnaryAddSubExpr)
	n.SLN.Accept(v)
	return v.Leave(n)
}

type StringListNullExpr struct {
	baseExpr

	PropertyOrLabelsExpr *PropertyOrLabelsExpr
	// elements would be `*StringOperationExpr`, `*ListOperationType` or `NullOperationType`
	SLNs []interface{}
}

type StringOperationType int

const (
	StringOperationStartsWith StringOperationType = iota
	StringOperationEndsWith
	StringOperationContains
)

type StringOperationExpr struct {
	baseExpr

	Type                 StringOperationType
	PropertyOrLabelsExpr *PropertyOrLabelsExpr
}

func (n *StringOperationExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*StringOperationExpr)
	n.PropertyOrLabelsExpr.Accept(v)
	return v.Leave(n)
}

type ListOperationType int

const (
	ListOperationIn ListOperationType = iota
	ListOperationSingle
	ListOperationRange
)

type ListOperationExpr struct {
	baseExpr

	InExpr     *PropertyOrLabelsExpr
	SingleExpr Expr
	// RangeExprs[0] is lower bound, RangeExprs[1] is upper bound
	RangeExprs [2]Expr
}

func (n *ListOperationExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ListOperationExpr)
	n.InExpr.Accept(v)
	n.SingleExpr.Accept(v)
	return v.Leave(n)
}

type NullOperationExpr struct {
	baseExpr

	// IsIsNull would be true with IS NULL, false with IS NOT NULL
	// TODO: looks so odd, change the name
	IsIsNull bool
}

func (n *NullOperationExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*NullOperationExpr)
	return v.Leave(n)
}

type PropertyOrLabelsExpr struct {
	baseExpr

	Atom            *Atom
	PropertyLookups []*PropertyLookup
	NodeLabels      []*NodeLabelNode
}

func (n *PropertyOrLabelsExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*PropertyOrLabelsExpr)
	n.Atom.Accept(v)
	for _, lookup := range n.PropertyLookups {
		lookup.Accept(v)
	}
	for _, label := range n.NodeLabels {
		label.Accept(v)
	}
	return v.Leave(n)
}

type AtomType int

const (
	AtomLiteral AtomType = iota
	AtomParameter
	AtomCase
	AtomCount
	AtomList
	AtomPatternComprehension
	AtomAllFilter
	AtomAnyFilter
	AtomNoneFilter
	AtomSingleFilter
	AtomPattern
	AtomParenthesizedExpr
	// AtomFuncInvocation
	AtomVariable
)

type Atom struct {
	Node

	Type                 AtomType
	Literal              *LiteralNode
	Parameter            *ParameterNode
	CaseExpr             *CaseExpr
	ListComprehension    *ListComprehension
	FilterExpr           *FilterExpr
	ParenthesizedExpr    Expr
	PatternComprehension *PatternComprehension
	PatternElement       *PatternElement
	Variable             *VariableNode
	// Function             *FunctionInvocation
}

func (n *Atom) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*Atom)
	switch n.Type {
	case AtomLiteral:
		n.Literal.Accept(v)
	case AtomParameter:
		n.Parameter.Accept(v)
	case AtomCase:
		n.CaseExpr.Accept(v)
	case AtomList:
		n.ListComprehension.Accept(v)
	case AtomAllFilter, AtomAnyFilter, AtomNoneFilter, AtomSingleFilter:
		n.FilterExpr.Accept(v)
	case AtomParenthesizedExpr:
		n.ParenthesizedExpr.Accept(v)
	case AtomPatternComprehension:
		n.PatternComprehension.Accept(v)
	case AtomPattern:
		n.PatternElement.Accept(v)
	case AtomVariable:
		n.Variable.Accept(v)
	}
	return v.Leave(n)
}

type PropertyLookup struct {
	Node

	PropertyKey *SchemaNameNode
}

func (n *PropertyLookup) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*PropertyLookup)
	n.PropertyKey.Accept(v)
	return v.Leave(n)
}

type CaseExpr struct {
	baseExpr

	Expr Expr
	Alts []*CaseAlt
	Else Expr
}

func (n *CaseExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*CaseExpr)
	n.Expr.Accept(v)
	for _, alt := range n.Alts {
		alt.Accept(v)
	}
	n.Else.Accept(v)
	return v.Leave(n)
}

type CaseAlt struct {
	baseExpr

	When Expr
	Then Expr
}

func (n *CaseAlt) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*CaseAlt)
	n.When.Accept(v)
	n.Then.Accept(v)
	return v.Leave(n)
}

type FilterExpr struct {
	baseExpr

	Variable *VariableNode
	In       Expr
	Where    Expr
}

func (n *FilterExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*FilterExpr)
	n.Variable.Accept(v)
	n.In.Accept(v)
	if n.Where != nil {
		n.Where.Accept(v)
	}
	return v.Leave(n)
}

type ListComprehension struct {
	baseExpr

	FilterExpr *FilterExpr
	Expr       Expr
}

func (n *ListComprehension) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ListComprehension)
	n.FilterExpr.Accept(v)
	if n.Expr != nil {
		n.Expr.Accept(v)
	}
	return v.Leave(n)
}

type FunctionInvocation struct {
	// TODO
	// don't support now
}
