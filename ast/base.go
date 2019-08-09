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

type baseNode struct {
	Node
	text string
}

func (n *baseNode) Text() string {
	return n.text
}

func (n *baseNode) SetText(text string) {
	n.text = text
}

type baseStmt struct {
	baseNode
}

func (n *baseStmt) statementNode() {}

type baseExpr struct {
	baseNode
}

func (n *baseExpr) exprNode() {}
