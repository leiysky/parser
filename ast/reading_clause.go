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

// ReadingClauseType represents types of ReadingClause
type ReadingClauseType byte

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

func (n *ReadingClause) Restore(ctx *RestoreContext) {
	switch n.Type {
	case ReadingClauseMatch:
		n.Match.Restore(ctx)
	case ReadingClauseUnwind:
		n.Unwind.Restore(ctx)
	}
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

func (n *MatchClause) Restore(ctx *RestoreContext) {
	if n.Optional {
		ctx.Write("OPTIONAL MATCH ")
	} else {
		ctx.Write("MATCH ")
	}
	n.Pattern.Restore(ctx)
	if n.Where != nil {
		ctx.WriteKeyword(" WHERE ")
		n.Where.Restore(ctx)
	}
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

func (n *UnwindClause) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("UNWIND ")
	n.Expr.Restore(ctx)
	ctx.WriteKeyword(" AS ")
	n.Variable.Restore(ctx)
}
