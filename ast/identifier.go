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

import (
	"strings"
)

type SchemaNameType byte

const (
	SchemaNameSymbolicName SchemaNameType = iota
	SchemaNameReservedWord
)

type SchemaNameNode struct {
	baseNode

	Type         SchemaNameType
	SymbolicName *SymbolicNameNode
	ReservedWord *ReservedWordNode
}

func (n *SchemaNameNode) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*SchemaNameNode)
	switch n.Type {
	case SchemaNameSymbolicName:
		n.SymbolicName.Accept(v)
	case SchemaNameReservedWord:
		n.ReservedWord.Accept(v)
	}
	return v.Leave(n)
}

func (n *SchemaNameNode) Restore(ctx *RestoreContext) {
	switch n.Type {
	case SchemaNameSymbolicName:
		n.SymbolicName.Restore(ctx)
	case SchemaNameReservedWord:
		n.ReservedWord.Restore(ctx)
	}
}

// SymbolicNameType is enum of SymbolicNameNode types
type SymbolicNameType byte

// SymbolicNameNode types
const (
	SymbolicNameUnescaped SymbolicNameType = iota
	SymbolicNameEscaped
	SymbolicNameHexLetter
	SymbolicNameCount
	SymbolicNameFilter
	SymbolicNameExtract
	SymbolicNameAny
	SymbolicNameNone
	SymbolicNameSingle
)

type SymbolicNameNode struct {
	baseNode
	Type  SymbolicNameType
	Value string
}

func (n *SymbolicNameNode) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*SymbolicNameNode)
	return v.Leave(n)
}

func (n *SymbolicNameNode) Restore(ctx *RestoreContext) {
	switch n.Type {
	case SymbolicNameUnescaped, SymbolicNameEscaped, SymbolicNameHexLetter:
		ctx.Write(n.Value)
	case SymbolicNameAny:
		ctx.WriteKeyword("ANY")
	case SymbolicNameCount:
		ctx.WriteKeyword("COUNT")
	case SymbolicNameExtract:
		ctx.WriteKeyword("EXTRACT")
	case SymbolicNameFilter:
		ctx.WriteKeyword("FILTER")
	case SymbolicNameNone:
		ctx.WriteKeyword("NONE")
	case SymbolicNameSingle:
		ctx.WriteKeyword("SINGLE")
	}
}

type ReservedWordNode struct {
	baseNode

	Content string
}

func (n *ReservedWordNode) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ReservedWordNode)
	return v.Leave(n)
}

func (n *ReservedWordNode) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword(n.Content)
}

type VariableNode struct {
	baseExpr

	SymbolicName *SymbolicNameNode
}

func (n *VariableNode) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*VariableNode)
	n.SymbolicName.Accept(v)
	return v.Leave(n)
}

func (n *VariableNode) Restore(ctx *RestoreContext) {
	ctx.Write("`")
	n.SymbolicName.Restore(ctx)
	ctx.Write("`")
}

func (n *VariableNode) Name() string {
	return n.SymbolicName.Value
}

type NodeLabelNode struct {
	baseNode

	LabelName *SchemaNameNode
}

// NodeLabelNode is actually a string token
func (n *NodeLabelNode) String() string {
	var str strings.Builder
	ctx := NewRestoreContext(&str)
	n.LabelName.Restore(ctx)
	return str.String()
}

func (n *NodeLabelNode) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*NodeLabelNode)
	n.LabelName.Accept(v)
	return v.Leave(n)
}

func (n *NodeLabelNode) Restore(ctx *RestoreContext) {
	ctx.Write(":")
	n.LabelName.Restore(ctx)
}

type DecimalInteger = int
type ParameterType byte

const (
	ParameterSymbolicName ParameterType = iota
	ParameterDecimalInteger
)

type ParameterNode struct {
	baseNode

	Type           ParameterType
	SymbolicName   *SymbolicNameNode
	DecimalInteger DecimalInteger
}

func (n *ParameterNode) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ParameterNode)
	switch n.Type {
	case ParameterSymbolicName:
		n.SymbolicName.Accept(v)
	case ParameterDecimalInteger:
		break
	}
	return v.Leave(n)
}

func (n *ParameterNode) Restore(ctx *RestoreContext) {
	switch n.Type {
	case ParameterSymbolicName:
		n.SymbolicName.Restore(ctx)
	case ParameterDecimalInteger:
		ctx.Write(n.DecimalInteger)
	}
}

type LiteralType byte

const (
	LiteralNumber LiteralType = iota
	LiteralString
	LiteralBoolean
	LiteralNull
	LiteralMap
	LiteralList
)

type LiteralExpr struct {
	baseExpr

	Type    LiteralType
	Number  *NumberLiteral
	String  string
	Boolean bool
	Map     *MapLiteral
	List    *ListLiteral
}

func (n *LiteralExpr) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*LiteralExpr)
	switch n.Type {
	case LiteralNumber:
		n.Number.Accept(v)
	case LiteralMap:
		n.Map.Accept(v)
	case LiteralList:
		n.List.Accept(v)
	}
	return v.Leave(n)
}

func (n *LiteralExpr) Restore(ctx *RestoreContext) {
	switch n.Type {
	case LiteralNumber:
		n.Number.Restore(ctx)
	case LiteralString:
		ctx.WriteString(n.String)
	case LiteralBoolean:
		if n.Boolean {
			ctx.WriteKeyword("TRUE")
		} else {
			ctx.WriteKeyword("FALSE")
		}
	case LiteralList:
		n.List.Restore(ctx)
	case LiteralMap:
		n.Map.Restore(ctx)
	}
}

type NumberLiteralType byte

const (
	NumberLiteralInteger NumberLiteralType = iota
	NumberLiteralDouble
)

// NumberLiteral represents literal value of number.
type NumberLiteral struct {
	baseNode

	Type    NumberLiteralType
	Integer int
	Double  float64
}

// Accept imlements ast.Node interface
func (n *NumberLiteral) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*NumberLiteral)
	return v.Leave(n)
}

// Restore would restore NumberLiteral value.
// It will use default format for int and float64
func (n *NumberLiteral) Restore(ctx *RestoreContext) {
	switch n.Type {
	case NumberLiteralInteger:
		ctx.Writef("%d", n.Integer)
	case NumberLiteralDouble:
		ctx.Writef("%f", n.Double)
	}
}

// type StringLiteral struct {
// 	baseNode

// 	String string
// }

// type BooleanLiteral struct {
// 	baseNode

// 	Boolean bool
// }

type MapLiteral struct {
	baseNode

	PropertyKeys []*SchemaNameNode
	Exprs        []Expr
}

func (n *MapLiteral) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*MapLiteral)
	for _, key := range n.PropertyKeys {
		key.Accept(v)
	}
	for _, expr := range n.Exprs {
		expr.Accept(v)
	}
	return v.Leave(n)
}

func (n *MapLiteral) Restore(ctx *RestoreContext) {
	ctx.Write("{")
	for i := range n.PropertyKeys {
		if i > 0 {
			ctx.Write(", ")
		}
		n.PropertyKeys[i].Restore(ctx)
		ctx.Write(": ")
		n.Exprs[i].Restore(ctx)
	}
	ctx.Write("}")
}

type ListLiteral struct {
	baseNode

	Exprs []Expr
}

func (n *ListLiteral) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ListLiteral)
	for _, expr := range n.Exprs {
		expr.Accept(v)
	}
	return v.Leave(n)
}

func (n *ListLiteral) Restore(ctx *RestoreContext) {
	ctx.Write("[")
	for i, expr := range n.Exprs {
		if i > 0 {
			ctx.Write(", ")
		}
		expr.Restore(ctx)
	}
}
