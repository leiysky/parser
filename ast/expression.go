// Copyright 2019 leiysky
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ast

var (
	_ Expr = &BinaryExpr{}
	_ Expr = &UnaryExpr{}
	_ Expr = &PredicationExpr{}
	_ Expr = &PropertyExpr{}
	_ Expr = &PropertyOrLabelsExpr{}
	_ Expr = &LiteralExpr{}
	_ Expr = &StringOperationExpr{}
	_ Expr = &NullOperationExpr{}
	_ Expr = &ListOperationExpr{}
	_ Expr = &Atom{}
	_ Expr = &CaseExpr{}
	_ Expr = &CaseAlt{}
	_ Expr = &ListComprehension{}
	_ Expr = &ParenExpr{}
	_ Expr = &FilterExpr{}
	_ Node = &PropertyLookup{}
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

func (n *PropertyExpr) Restore(ctx *RestoreContext) {
	n.Atom.Restore(ctx)
	for _, l := range n.Lookups {
		l.Restore(ctx)
	}
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

func (n *BinaryExpr) Restore(ctx *RestoreContext) {
	n.L.Restore(ctx)
	ctx.Writef(" %s ", n.Op)
	n.R.Restore(ctx)
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

func (n *UnaryExpr) Restore(ctx *RestoreContext) {
	ctx.Write(n.Op.String())
	n.V.Restore(ctx)
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

func (n *PredicationExpr) Restore(ctx *RestoreContext) {
	n.Expr.Restore(ctx)
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

func (n *StringOperationExpr) Restore(ctx *RestoreContext) {
	switch n.Type {
	case StringOperationStartsWith:
		ctx.WriteKeyword("STARTS WITH ")
		n.Expr.Restore(ctx)
	case StringOperationEndsWith:
		ctx.WriteKeyword("ENDS WITH ")
		n.Expr.Restore(ctx)
	case StringOperationContains:
		ctx.WriteKeyword("CONTAINS ")
		n.Expr.Restore(ctx)
	}
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

func (n *ListOperationExpr) Restore(ctx *RestoreContext) {
	if n.InExpr != nil {
		ctx.WriteKeyword("IN ")
		n.InExpr.Restore(ctx)
	} else if n.SingleExpr != nil {
		ctx.Write("[")
		n.SingleExpr.Restore(ctx)
		ctx.Write("]")
	} else {
		ctx.Write("[")
		n.LowerBound.Restore(ctx)
		ctx.Write("..")
		n.UpperBound.Restore(ctx)
		ctx.Write("]")
	}
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

func (n *NullOperationExpr) Restore(ctx *RestoreContext) {
	if n.IsIsNull {
		ctx.WriteKeyword("IS NULL")
	} else {
		ctx.WriteKeyword("IS NOT NULL")
	}
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

func (n *PropertyOrLabelsExpr) Restore(ctx *RestoreContext) {
	n.Atom.Restore(ctx)
	for _, l := range n.PropertyLookups {
		l.Restore(ctx)
	}
	for _, label := range n.NodeLabels {
		label.Restore(ctx)
	}
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

func (n *Atom) Restore(ctx *RestoreContext) {
	switch n.Type {
	case AtomLiteral:
		n.Literal.Restore(ctx)
	case AtomParameter:
		n.Parameter.Restore(ctx)
	case AtomCase:
		n.CaseExpr.Restore(ctx)
	case AtomCount:
		ctx.WriteKeyword("COUNT(*)")
	case AtomList:
		n.ListComprehension.Restore(ctx)
	case AtomPatternComprehension:
		n.PatternComprehension.Restore(ctx)
	case AtomAllFilter:
		ctx.WriteKeyword("ALL ")
		n.FilterExpr.Restore(ctx)
	case AtomAnyFilter:
		ctx.WriteKeyword("ANY ")
		n.FilterExpr.Restore(ctx)
	case AtomNoneFilter:
		ctx.WriteKeyword("NONE ")
		n.FilterExpr.Restore(ctx)
	case AtomSingleFilter:
		ctx.WriteKeyword("SINGLE ")
		n.FilterExpr.Restore(ctx)
	case AtomPattern:
		n.PatternElement.Restore(ctx)
	case AtomParenthesizedExpr:
		n.ParenthesizedExpr.Restore(ctx)
	case AtomVariable:
		n.Variable.Restore(ctx)
	}
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

func (n *PropertyLookup) Restore(ctx *RestoreContext) {
	ctx.Write(".")
	ctx.Write("`")
	n.PropertyKey.Restore(ctx)
	ctx.Write("`")
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

func (n *CaseExpr) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("CASE")
	if n.Expr != nil {
		ctx.Write(" ")
		n.Expr.Restore(ctx)
	}
	for _, alt := range n.Alts {
		ctx.Write(" ")
		alt.Restore(ctx)
	}
	if n.Else != nil {
		ctx.WriteKeyword(" ELSE ")
		n.Else.Restore(ctx)
	}
	ctx.WriteKeyword(" END")
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

func (n *CaseAlt) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("WHEN ")
	n.When.Restore(ctx)
	ctx.WriteKeyword("THEN ")
	n.Then.Restore(ctx)
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

func (n *FilterExpr) Restore(ctx *RestoreContext) {
	n.Variable.Restore(ctx)
	ctx.WriteKeyword(" IN ")
	n.In.Restore(ctx)
	if n.Where != nil {
		ctx.WriteKeyword(" WHERE ")
		n.Where.Restore(ctx)
	}
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

func (n *ListComprehension) Restore(ctx *RestoreContext) {
	ctx.Write("[")
	n.FilterExpr.Restore(ctx)
	if n.Expr != nil {
		ctx.Write(" | ")
		n.Expr.Restore(ctx)
	}
	ctx.Write("]")
}

type FunctionInvocation struct {
	// TODO
	// don't support now
}

type ParenExpr struct {
	baseExpr

	Expr Expr
}

func (n *ParenExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ParenExpr)
	n.Expr.Accept(v)
	return v.Leave(n)
}

func (n *ParenExpr) Restore(ctx *RestoreContext) {
	ctx.Write("(")
	n.Expr.Restore(ctx)
	ctx.Write(")")
}
