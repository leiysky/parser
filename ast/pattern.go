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

type Pattern struct {
	baseNode
	Parts []*PatternPart
}

func (n *Pattern) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*Pattern)
	for _, part := range n.Parts {
		part.Accept(v)
	}
	return v.Leave(n)
}

func (n *Pattern) Restore(ctx *RestoreContext) {
	for i, part := range n.Parts {
		if i > 0 {
			ctx.Write(", ")
		}
		part.Restore(ctx)
	}
}

type PatternPart struct {
	baseNode

	Variable *VariableNode
	Element  *PatternElement
}

func (n *PatternPart) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*PatternPart)
	if n.Variable != nil {
		n.Variable.Accept(v)
	}
	n.Element.Accept(v)
	return v.Leave(n)
}

func (n *PatternPart) Restore(ctx *RestoreContext) {
	if n.Variable != nil {
		n.Variable.Restore(ctx)
		ctx.Write(" = ")
	}
	n.Element.Restore(ctx)
}

type PatternElement struct {
	baseNode

	// Amount of Relationships is the same as that of Nodes
	// Nodes[i], Relationships[i], Nodes[i+1]...
	Relationships []*RelationshipPattern
	Nodes         []*NodePattern
}

func (n *PatternElement) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*PatternElement)
	for i := range n.Relationships {
		n.Nodes[i].Accept(v)
		n.Relationships[i].Accept(v)
	}
	n.Nodes[len(n.Nodes)-1].Accept(v)
	return v.Leave(n)
}

func (n *PatternElement) Restore(ctx *RestoreContext) {
	// n.StartNode.Restore(ctx)
	for i := range n.Relationships {
		n.Nodes[i].Restore(ctx)
		n.Relationships[i].Restore(ctx)
	}
	n.Nodes[len(n.Nodes)-1].Restore(ctx)
}

type NodePattern struct {
	baseNode

	Variable   *VariableNode
	Labels     []*NodeLabelNode
	Properties *Properties
}

func (n *NodePattern) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*NodePattern)
	if n.Variable != nil {
		n.Variable.Accept(v)
	}
	for _, label := range n.Labels {
		label.Accept(v)
	}
	if n.Properties != nil {
		n.Properties.Accept(v)
	}
	return v.Leave(n)
}

func (n *NodePattern) Restore(ctx *RestoreContext) {
	ctx.Write("(")
	if n.Variable != nil {
		n.Variable.Restore(ctx)
	}
	for _, label := range n.Labels {
		label.Restore(ctx)
	}
	if n.Properties != nil {
		n.Properties.Restore(ctx)
	}
	ctx.Write(")")
}

// RelationshipType represents types of Relationships
type RelationshipType byte

const (
	// RelationshipLeft represents the Relationship only points to left
	RelationshipLeft RelationshipType = iota
	// RelationshipRight represents the Relationship only points to right
	RelationshipRight
	// RelationshipBoth represents the Relationship points to both left and right
	RelationshipBoth
	// RelationshipNone represents the Relationship has no direction
	RelationshipNone
)

type RelationshipPattern struct {
	baseNode

	Type   RelationshipType
	Detail *RelationshipDetail
}

func (n *RelationshipPattern) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*RelationshipPattern)
	n.Detail.Accept(v)
	return v.Leave(n)
}

func (n *RelationshipPattern) Restore(ctx *RestoreContext) {
	switch n.Type {
	case RelationshipLeft:
		ctx.Write("<-")
		n.Detail.Restore(ctx)
		ctx.Write("-")
	case RelationshipRight:
		ctx.Write("-")
		n.Detail.Restore(ctx)
		ctx.Write("->")
	case RelationshipBoth:
		ctx.Write("<-")
		n.Detail.Restore(ctx)
		ctx.Write("->")
	case RelationshipNone:
		ctx.Write("-")
		n.Detail.Restore(ctx)
		ctx.Write("-")
	}
}

type RelationshipDetail struct {
	baseNode

	Variable          *VariableNode
	RelationshipTypes []*SchemaNameNode
	// -1 represent wildcard
	// [1, 2] means the Relationship will be matched for 1 - 2 times
	// [-1, 2] means will be matched for 0 - 2 times
	// [1, -1] means will be matched for 1 time or more
	// [-1, -1] means will be matched for any times
	// There are 5 conditions:
	// - [*]
	// - [*1]
	// - [*1..]
	// - [*1..2]
	// - [*..2]
	MinHops    int
	MaxHops    int
	Properties *Properties
}

func (n *RelationshipDetail) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*RelationshipDetail)
	if n.Variable != nil {
		n.Variable.Accept(v)
	}
	for _, label := range n.RelationshipTypes {
		label.Accept(v)
	}
	if n.Properties != nil {
		n.Properties.Accept(v)
	}
	return v.Leave(n)
}

func (n *RelationshipDetail) Restore(ctx *RestoreContext) {
	ctx.Write("[")
	if n.Variable != nil {
		n.Variable.Restore(ctx)
	}
	for i, t := range n.RelationshipTypes {
		if i == 0 {
			ctx.Write(":")
		} else {
			ctx.Write("|")
		}
		t.Restore(ctx)
	}

	ctx.Write("*")
	if n.MinHops > -1 {
		ctx.Write(n.MinHops)
	}
	ctx.Write("..")
	if n.MaxHops > -1 {
		ctx.Write(n.MaxHops)
	}

	if n.Properties != nil {
		n.Properties.Restore(ctx)
	}
	ctx.Write("]")
}

type PropertiesType byte

const (
	PropertiesMapLiteral PropertiesType = iota
	PropertiesParameter
)

type Properties struct {
	baseNode

	Type       PropertiesType
	MapLiteral *MapLiteral
	Parameter  *ParameterNode
}

func (n *Properties) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*Properties)
	switch n.Type {
	case PropertiesMapLiteral:
		n.MapLiteral.Accept(v)
	case PropertiesParameter:
		n.Parameter.Accept(v)
	}
	return v.Leave(n)
}

func (n *Properties) Restore(ctx *RestoreContext) {
	switch n.Type {
	case PropertiesMapLiteral:
		n.MapLiteral.Restore(ctx)
	case PropertiesParameter:
		n.Parameter.Restore(ctx)
	}
}

type PatternComprehension struct {
	baseNode

	Variable       *VariableNode
	PatternElement *PatternElement
	Where          Expr
	Expr           Expr
}

func (n *PatternComprehension) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*PatternComprehension)
	if n.Variable != nil {
		n.Variable.Accept(v)
	}
	n.PatternElement.Accept(v)
	if n.Where != nil {
		n.Where.Accept(v)
	}
	n.Expr.Accept(v)
	return v.Leave(n)
}

func (n *PatternComprehension) Restore(ctx *RestoreContext) {
	ctx.Write("[")
	if n.Variable != nil {
		n.Variable.Restore(ctx)
		ctx.Write(" = ")
	}
	n.PatternElement.Restore(ctx)
	if n.Where != nil {
		ctx.WriteKeyword(" WHERE")
		n.Where.Restore(ctx)
	}
	ctx.Write(" | ")
	n.Expr.Restore(ctx)
	ctx.Write("]")
}
