package ast

type SchemaNameType int

const (
	SchemaNameSymbolicName SchemaNameType = iota
	SchemaNameReservedWord
)

type SchemaNameNode struct {
	baseNode

	Type SchemaNameType
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
