package ast

type Pattern struct {
	baseNode
	Parts []*PatternPart
}

type PatternPart struct {
	baseNode
	WithVariable bool
	Variable     *SymbolicNameNode
	Element      *PatternElement
}

type PatternElement struct {
	baseNode
}
