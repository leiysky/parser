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

func (n *WithClause) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("WITH ")
	if n.Distinct {
		ctx.WriteKeyword("DISTINCT ")
	}
	n.ReturnBody.Restore(ctx)
	if n.Where != nil {
		ctx.Write(" ")
		n.Where.Restore(ctx)
	}
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

func (n *ReturnClause) Restore(ctx *RestoreContext) {
	ctx.Write("RETURN ")
	if n.Distinct {
		ctx.WriteKeyword("DISTINCT ")
	}
	n.ReturnBody.Restore(ctx)
}

type ReturnBody struct {
	baseStmt

	ReturnItems []*ReturnItem
	OrderBy     *OrderClause
	Skip        Expr
	Limit       Expr
}

func (n *ReturnBody) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ReturnBody)
	for _, item := range n.ReturnItems {
		item.Accept(v)
	}
	if n.OrderBy != nil {
		n.OrderBy.Accept(v)
	}
	if n.Skip != nil {
		n.Skip.Accept(v)
	}
	if n.Limit != nil {
		n.Limit.Accept(v)
	}
	return v.Leave(n)
}

func (n *ReturnBody) Restore(ctx *RestoreContext) {
	for i, item := range n.ReturnItems {
		if i > 0 {
			ctx.Write(", ")
		}
		item.Restore(ctx)
	}
	if n.OrderBy != nil {
		ctx.Write(" ")
		n.OrderBy.Restore(ctx)
	}
	if n.Skip != nil {
		ctx.Write(" ")
		n.Skip.Restore(ctx)
	}
	if n.Limit != nil {
		ctx.Write(" ")
		n.Limit.Restore(ctx)
	}
}

type ReturnItem struct {
	baseStmt

	Wildcard bool
	Expr     Expr
	As       bool
	Variable *VariableNode
}

func (n *ReturnItem) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*ReturnItem)
	n.Expr.Accept(v)
	if n.As {
		n.Variable.Accept(v)
	}
	return v.Leave(n)
}

func (n *ReturnItem) Restore(ctx *RestoreContext) {
	if n.Wildcard {
		ctx.Write("*")
		return
	}
	n.Expr.Restore(ctx)
	if n.As {
		ctx.WriteKeyword(" AS ")
		n.Variable.Restore(ctx)
	}
}

type OrderClause struct {
	baseStmt

	SortItems []*SortItem
}

func (n *OrderClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*OrderClause)
	for _, item := range n.SortItems {
		item.Accept(v)
	}
	return v.Leave(n)
}

func (n *OrderClause) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("ORDER BY ")
	for i, item := range n.SortItems {
		if i > 0 {
			ctx.Write(", ")
		}
		item.Restore(ctx)
	}
}

type SortType byte

const (
	SortAscending SortType = iota
	SortDescending
)

type SortItem struct {
	baseStmt

	Type SortType
	Expr Expr
}

func (n *SortItem) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*SortItem)
	n.Expr.Accept(v)
	return v.Leave(n)
}

func (n *SortItem) Restore(ctx *RestoreContext) {
	n.Expr.Restore(ctx)
	switch n.Type {
	case SortAscending:
		ctx.WriteKeyword(" ASC")
	case SortDescending:
		ctx.WriteKeyword(" DESC")
	}
}
