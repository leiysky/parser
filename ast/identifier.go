package ast

type SchemaNameType int

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

// SymbolicNameType is enum of SymbolicNameNode types
type SymbolicNameType int

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

type VariableNode struct {
	baseNode

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

type NodeLabelNode struct {
	baseNode

	LabelName *SchemaNameNode
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

type DecimalInteger = int
type ParameterType int

const (
	ParameterSymbolicname ParameterType = iota
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
	case ParameterSymbolicname:
		n.SymbolicName.Accept(v)
	case ParameterDecimalInteger:
		break
	}
	return v.Leave(n)
}

type LiteralType int

const (
	LiteralNumber LiteralType = iota
	LiteralString
	LiteralBoolean
	LiteralNull
	LiteralMap
	LiteralList
)

type LiteralNode struct {
	baseNode

	Type    LiteralType
	Number  *NumberLiteral
	String  string
	Boolean bool
	Map     *MapLiteral
	List    *ListLiteral
}

func (n *LiteralNode) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*LiteralNode)
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

type NumberLiteralType int

const (
	NumberLiteralInteger NumberLiteralType = iota
	NumberLiteralDouble
)

type NumberLiteral struct {
	baseNode

	Type    NumberLiteralType
	Integer int
	Double  float64
}

func (n *NumberLiteral) Accept(v Visitor) (Node, bool) {
	newNode, skip := v.Enter(n)
	if skip {
		return v.Leave(n)
	}
	n = newNode.(*NumberLiteral)
	return v.Leave(n)
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
