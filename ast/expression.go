package ast

var (
	_ Expr = &BinaryExpr{}
	_ Expr = &UnaryExpr{}
	_ Expr = &PredicationExpr{}
	_ Expr = &Atom{}
	_ Expr = &PropertyExpr{}
	_ Expr = &PropertyOrLabelsExpr{}
	_ Expr = &LiteralExpr{}
)

// PropertyExpr represents a property lookup expression like `a.b.c`.
// Different from PropertyOrLabelsExpr, it has at least one PropertyLookup.
type PropertyExpr struct {
	baseExpr

	Atom    *Atom
	Lookups []*PropertyLookup
}

// Accept implements Node interface
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

// BinaryExpr represents a binary expression with left expression, right expression and an operator.
type BinaryExpr struct {
	baseExpr

	Op OpType
	L  Expr
	R  Expr
}

// Accept implements Node interface
func (n *BinaryExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*BinaryExpr)
	n.L.Accept(v)
	n.R.Accept(v)
	return v.Leave(n)
}

// UnaryExpr represents a unary expression with expression and an operator.
type UnaryExpr struct {
	baseExpr

	Op OpType
	V  Expr
}

// Accept implements Node interface
func (n *UnaryExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*UnaryExpr)
	n.V.Accept(v)
	return v.Leave(n)
}

// PredicationType represents types of PredicationExpr
type PredicationType byte

// There are 3 kinds of predication expression.
const (
	PredicationStringOp PredicationType = iota
	PredicationListOp
	PredicationNullOp
)

// PredicationExpr represents a expression with boolean value but not logical operation.
type PredicationExpr struct {
	baseExpr

	Type PredicationType
	Expr Expr
}

// Accept implements Node interface
func (n *PredicationExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*PredicationExpr)
	n.Expr.Accept(v)
	return v.Leave(n)
}

// OpType represents operator type of expression
type OpType byte

// OpTypes
const (
	OpAdd OpType = iota
	OpSub
	OpMul
	OpDiv
	OpMod
	OpPow
	OpPlus  // positive
	OpMinus // negative
	OpEQ
	OpNE
	OpLT
	OpGT
	OpLTE
	OpGTE
	OpOr
	OpAnd
	OpXor
	OpNot
)

// String implements fmt.Stringer interface
func (op OpType) String() string {
	switch op {
	case OpAdd, OpPlus:
		return "+"
	case OpSub, OpMinus:
		return "-"
	case OpMul:
		return "*"
	case OpDiv:
		return "/"
	case OpMod:
		return "%"
	case OpPow:
		return "^"
	case OpEQ:
		return "="
	case OpNE:
		return "<>"
	case OpLT:
		return "<"
	case OpGT:
		return ">"
	case OpLTE:
		return "<="
	case OpGTE:
		return ">="
	case OpOr:
		return "OR"
	case OpAnd:
		return "AND"
	case OpXor:
		return "XOR"
	case OpNot:
		return "NOT"
	default:
		return "<unknown>"
	}
}

type StringOperationType byte

const (
	StringOperationStartsWith StringOperationType = iota
	StringOperationEndsWith
	StringOperationContains
)

type StringOperationExpr struct {
	baseExpr

	Type StringOperationType
	Expr *PropertyOrLabelsExpr
}

func (n *StringOperationExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*StringOperationExpr)
	n.Expr.Accept(v)
	return v.Leave(n)
}

type ListOperationType byte

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
	LowerBound Expr
	UpperBound Expr
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

type AtomType byte

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
	baseExpr

	Type                 AtomType
	Literal              *LiteralExpr
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
