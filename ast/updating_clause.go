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

// CreateClause represents CREATE clause node
type CreateClause struct {
	baseStmt

	Pattern *Pattern
}

func (n *CreateClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*CreateClause)
	n.Pattern.Accept(v)
	return v.Leave(n)
}

func (n *CreateClause) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("CREATE ")
	n.Pattern.Restore(ctx)
}

// MergeClause represents MERGE clause node
type MergeClause struct {
	baseStmt

	PatternPart  *PatternPart
	MergeActions []*MergeAction
}

func (n *MergeClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*MergeClause)
	n.PatternPart.Accept(v)
	for _, action := range n.MergeActions {
		action.Accept(v)
	}
	return v.Leave(n)
}

func (n *MergeClause) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("MERGE ")
	n.PatternPart.Restore(ctx)
	for i, action := range n.MergeActions {
		if i > 0 {
			ctx.Write(" ")
		}
		action.Restore(ctx)
	}
}

type MergeActionType byte

const (
	MergeActionCreate MergeActionType = iota
	MergeActionMatch
)

// MergeAction represents optional action node (ON CREATE or ON MERGE) in MERGE clause node
type MergeAction struct {
	baseStmt

	// Action is with values "MERGE" and "CREATE"
	Type MergeActionType
	Set  *SetClause
}

func (n *MergeAction) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*MergeAction)
	n.Set.Accept(v)
	return v.Leave(n)
}

func (n *MergeAction) Restore(ctx *RestoreContext) {
	switch n.Type {
	case MergeActionCreate:
		ctx.WriteKeyword("ON CREATE ")
	case MergeActionMatch:
		ctx.WriteKeyword("ON MATCH ")
	}
	n.Set.Restore(ctx)
}

// SetClause represents SET clause node
type SetClause struct {
	baseStmt

	SetItems []*SetItem
}

func (n *SetClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*SetClause)
	for _, item := range n.SetItems {
		item.Accept(v)
	}
	return v.Leave(n)
}

func (n *SetClause) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("SET ")
	for i, item := range n.SetItems {
		if i > 0 {
			ctx.Write(", ")
		}
		item.Restore(ctx)
	}
}

type SetItemType byte

const (
	SetItemProperty SetItemType = iota
	SetItemVariableAssignment
	SetItemVariableIncrement
	SetItemVariableLabel
)

type SetItem struct {
	baseStmt

	Type     SetItemType
	Property *PropertyExpr
	Variable *VariableNode
	Expr     Expr
	Labels   []*NodeLabelNode
}

func (n *SetItem) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*SetItem)
	switch n.Type {
	case SetItemProperty:
		n.Property.Accept(v)
		n.Expr.Accept(v)

	case SetItemVariableAssignment, SetItemVariableIncrement:
		n.Variable.Accept(v)
		n.Expr.Accept(v)

	case SetItemVariableLabel:
		n.Variable.Accept(v)
		for _, label := range n.Labels {
			label.Accept(v)
		}
	}
	return v.Leave(n)
}

func (n *SetItem) Restore(ctx *RestoreContext) {
	switch n.Type {
	case SetItemProperty:
		n.Property.Restore(ctx)
		ctx.Write(" = ")
		n.Expr.Restore(ctx)
	case SetItemVariableAssignment:
		n.Variable.Restore(ctx)
		ctx.Write(" = ")
		n.Expr.Restore(ctx)
	case SetItemVariableIncrement:
		n.Variable.Restore(ctx)
		ctx.Write(" += ")
		n.Expr.Restore(ctx)
	case SetItemVariableLabel:
		n.Variable.Restore(ctx)
		for _, label := range n.Labels {
			label.Restore(ctx)
		}
	}
}

type DeleteClause struct {
	baseStmt

	Detach bool
	Exprs  []Expr
}

func (n *DeleteClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*DeleteClause)
	for _, expr := range n.Exprs {
		expr.Accept(v)
	}
	return v.Leave(n)
}

func (n *DeleteClause) Restore(ctx *RestoreContext) {
	if n.Detach {
		ctx.WriteKeyword("DETACH ")
	}
	ctx.WriteKeyword("DELETE ")
	for i, expr := range n.Exprs {
		if i > 0 {
			ctx.Write(", ")
		}
		expr.Restore(ctx)
	}
}

type RemoveClause struct {
	baseStmt

	RemoveItems []*RemoveItem
}

func (n *RemoveClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*RemoveClause)
	for _, item := range n.RemoveItems {
		item.Accept(v)
	}
	return v.Leave(n)
}

func (n *RemoveClause) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("REMOVE ")
	for i, item := range n.RemoveItems {
		if i > 0 {
			ctx.Write(", ")
		}
		item.Restore(ctx)
	}
}

type RemoveItemType byte

const (
	RemoveItemVariable RemoveItemType = iota
	RemoveItemProperty
)

type RemoveItem struct {
	baseStmt

	Type     RemoveItemType
	Variable *VariableNode
	Labels   []*NodeLabelNode
	Property *PropertyExpr
}

func (n *RemoveItem) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*RemoveItem)
	switch n.Type {
	case RemoveItemVariable:
		n.Variable.Accept(v)
		for _, label := range n.Labels {
			label.Accept(v)
		}

	case RemoveItemProperty:
		n.Property.Accept(v)
	}
	return v.Leave(n)
}

func (n *RemoveItem) Restore(ctx *RestoreContext) {
	switch n.Type {
	case RemoveItemVariable:
		n.Variable.Restore(ctx)
		for _, label := range n.Labels {
			label.Restore(ctx)
		}
	case RemoveItemProperty:
		n.Property.Restore(ctx)
	}
}
