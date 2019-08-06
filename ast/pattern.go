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

type PatternPart struct {
	baseNode
	WithVariable bool
	Variable     *VariableNode
	Element      *PatternElement
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

type PatternElement struct {
	baseNode

	// Amount of Relationships is the same as that of Nodes
	// The pattern sequence should be (n is `len(Relationships)` also `len(Nodes)`):
	// StartNode, Relationships[0], Nodes[0], ... , Relationships[n-1], Nodes[n-1]
	StartNode     *NodePattern
	Relationships []*RelationshipPattern
	// Nodes represents NodePatterns exclude StartNode
	Nodes []*NodePattern
}

func (n *PatternElement) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*PatternElement)
	n.StartNode.Accept(v)
	for _, rel := range n.Relationships {
		rel.Accept(v)
	}
	for _, node := range n.Nodes {
		node.Accept(v)
	}
	return v.Leave(n)
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

// RelationshipType represents types of Relationships
type RelationshipType int

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

type PropertiesType int

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

type PatternComprehension struct {
	baseNode

	Variable       *VariableNode
	PatternElement *PatternElement
	Where          *Expr
	Expr           *Expr
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
