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

type CypherStmtType byte

const (
	CypherStmtQuery CypherStmtType = iota
	CypherStmtStandaloneCall
)

type CypherStmt struct {
	baseStmt

	Type           CypherStmtType
	Query          *QueryStmt
	StandaloneCall *StandaloneCall
}

func (n *CypherStmt) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*CypherStmt)
	switch n.Type {
	case CypherStmtQuery:
		n.Query.Accept(v)
	case CypherStmtStandaloneCall:
		n.StandaloneCall.Accept(v)
	}
	return v.Leave(n)
}

func (n *CypherStmt) Restore(ctx *RestoreContext) {
	switch n.Type {
	case CypherStmtQuery:
		n.Query.Restore(ctx)
	}
}

type QueryStmt struct {
	baseStmt

	Clauses []Stmt
}

func (n *QueryStmt) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*QueryStmt)
	for _, c := range n.Clauses {
		c.Accept(v)
	}
	return v.Leave(n)
}

func (n *QueryStmt) Restore(ctx *RestoreContext) {
	for i, c := range n.Clauses {
		if i > 0 {
			ctx.Write(" ")
		}
		c.Restore(ctx)
	}
}

type UnionClause struct {
	baseStmt

	All     bool
	Clauses []Stmt
}

func (n *UnionClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*UnionClause)
	for _, c := range n.Clauses {
		c.Accept(v)
	}
	return v.Leave(n)
}

func (n *UnionClause) Restore(ctx *RestoreContext) {
	ctx.WriteKeyword("UNION ")
	if n.All {
		ctx.WriteKeyword("ALL ")
	}
	for _, c := range n.Clauses {
		c.Restore(ctx)
	}
}

type StandaloneCall struct {
	baseNode
}

func (n *StandaloneCall) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*StandaloneCall)
	// TODO
	return v.Leave(n)
}
