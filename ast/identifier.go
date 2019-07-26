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

type MapLiteralNode struct {
	baseNode

	PropertyKeys []*SymbolicNameNode
	Exprs        []*Expr
}

type ParameterType int

type DecimalInteger int

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
