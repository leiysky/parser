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

// Node is basic type of ast.
type Node interface {
	Accept(visitor Visitor) (node Node, ok bool)
	Text() string
	SetText(text string)
	Restore(ctx *RestoreContext)
}

// Stmt represents statement node.
type Stmt interface {
	Node
	statementNode()
}

// Expr represents expression node.
type Expr interface {
	Node
	exprNode()
}

// Visitor can visit ast.
type Visitor interface {
	// Enter would be called at the begin of Accept.
	// The visited node will be replaced by returned Node.
	// If skipChildren is true, the child nodes' Accept wouldn't be called.
	Enter(n Node) (node Node, skipChildren bool)

	// Accept will directly return Leave.
	Leave(n Node) (node Node, ok bool)
}
