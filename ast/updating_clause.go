package ast

type UpdatingClauseType int

const (
	UpdatingClauseCreate UpdatingClauseType = iota
	UpdatingClauseMerge
	UpdatingClauseSet
)

type UpdatingClause struct {
	baseStmt

	Type   UpdatingClauseType
	Create *CreateClause
	Merge  *MergeClause
	Set    *SetClause
}

func (n *UpdatingClause) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*UpdatingClause)
	switch n.Type {
	case UpdatingClauseCreate:
		n.Create.Accept(v)
	case UpdatingClauseMerge:
		n.Merge.Accept(v)
	case UpdatingClauseSet:
		n.Set.Accept(v)
	}
	return v.Leave(n)
}

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

type MergeActionType int

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

// SetClause represents SET clause node
type SetClause struct {
	baseStmt

	SetItems []*SetItemStmt
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

type SetItemType int

const (
	SetItemProperty SetItemType = iota
	SetItemVariableAssignment
	SetItemVariableIncrement
	SetItemVariableLabel
)

type SetItemStmt struct {
	baseStmt

	Type     SetItemType
	Property *PropertyExpr
	Variable *SymbolicNameNode
	Expr     *Expr
	Labels   []*SchemaNameNode
}

func (n *SetItemStmt) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*SetItemStmt)
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

type DeleteClause struct {
	baseStmt

	Detach bool
	Exprs  []*Expr
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

type RemoveClause struct {
	baseStmt
}

type RemoveItemType int

const (
	RemoveItemVariable RemoveItemType = iota
	RemoveItemProperty
)

type RemoveItemStmt struct {
	baseStmt

	Type     RemoveItemType
	Variable *SymbolicNameNode
	Labels   []*SchemaNameNode
	Property *PropertyExpr
}

func (n *RemoveItemStmt) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*RemoveItemStmt)
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
