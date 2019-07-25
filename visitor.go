package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/leiysky/parser/ast"
)

type visitor struct {
	CypherVisitor
	parser *CypherParser
	// *BaseCypherVisitor
}

func (v *visitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

// func (v *visitor) VisitChildren(node antlr.RuleNode) interface{} {
// 	children := node.GetChildren()
// 	for _, c := range children {
// 		switch n := c.(type) {
// 		case antlr.RuleNode:
// 			n.Accept(v)

// 		case antlr.TerminalNode:
// 			v.VisitTerminal(n)
// 		case antlr.ErrorNode:
// 			v.VisitErrorNode(n)
// 		}
// 	}
// 	return nil
// }

func (v *visitor) VisitChildren(node antlr.RuleNode) interface{} {
	return nil
}

func (v *visitor) VisitTerminal(node antlr.TerminalNode) interface{} { return nil }

func (v *visitor) VisitErrorNode(node antlr.ErrorNode) interface{} { return nil }

func (v *visitor) VisitCypher(ctx *CypherContext) interface{} {
	node := &ast.CypherStmt{}
	stmt := ctx.Stmt().Accept(v).(*StmtContext)
	query := stmt.Query().Accept(v).(*QueryContext)

	var i antlr.RuleContext
	if query.RegularQuery() != nil {
		i = query.RegularQuery()
	} else if query.StandaloneCall() != nil {
		i = query.StandaloneCall()
	}

	switch v.parser.RuleNames[i.GetRuleIndex()] {
	case "regularQuery":
		node.Type = ast.CypherStmtQuery
		node.Query = query.RegularQuery().Accept(v).(*ast.RegularQueryStmt)

	case "standaloneCall":
		node.Type = ast.CypherStmtStandaloneCall
		node.StandaloneCall = query.StandaloneCall().Accept(v).(*ast.StandaloneCall)
	}
	return node
}

func (v *visitor) VisitStmt(ctx *StmtContext) interface{} {
	return ctx
}

func (v *visitor) VisitQuery(ctx *QueryContext) interface{} {
	return ctx
}

func (v *visitor) VisitRegularQuery(ctx *RegularQueryContext) interface{} {
	regularQuery := &ast.RegularQueryStmt{}
	regularQuery.SingleQuery = ctx.SingleQuery().Accept(v).(*ast.SingleQueryStmt)

	var unionClauses []*ast.UnionClause
	for _, union := range ctx.AllUnionClause() {
		unionClauses = append(unionClauses, union.Accept(v).(*ast.UnionClause))
	}
	regularQuery.Unions = unionClauses
	return regularQuery
}

func (v *visitor) VisitSingleQuery(ctx *SingleQueryContext) interface{} {
	singleQuery := &ast.SingleQueryStmt{}

	var i antlr.RuleContext
	if ctx.SinglePartQuery() != nil {
		i = ctx.SinglePartQuery()
	} else if ctx.MultiPartQuery() != nil {
		i = ctx.MultiPartQuery()
	}

	switch v.parser.RuleNames[i.GetRuleIndex()] {
	case "singlePartQuery":
		singleQuery.Type = ast.SingleQueryStmtSinglePart
		singleQuery.SinglePart = ctx.SinglePartQuery().Accept(v).(*ast.SinglePartQueryStmt)

	case "multiPartQuery":
		singleQuery.Type = ast.SingleQueryStmtMultiPart
		singleQuery.MultiPart = ctx.SinglePartQuery().Accept(v).(*ast.MultiPartQueryStmt)
	}
	return singleQuery
}

func (v *visitor) VisitSinglePartQuery(ctx *SinglePartQueryContext) interface{} {
	singlePartQuery := &ast.SinglePartQueryStmt{}
	var readingClauses []*ast.ReadingClause
	for _, c := range ctx.AllReadingClause() {
		readingClauses = append(readingClauses, c.Accept(v).(*ast.ReadingClause))
	}
	return singlePartQuery
}

func (v *visitor) VisitMultiPartQuery(ctx *MultiPartQueryContext) interface{} {
	multiPartQuery := &ast.MultiPartQueryStmt{}

	return multiPartQuery
}

func (v *visitor) VisitReadingClause(ctx *ReadingClauseContext) interface{} {
	readingClause := &ast.ReadingClause{}

	var i antlr.RuleContext
	if ctx.MatchClause() != nil {
		i = ctx.MatchClause()
	} else if ctx.UnwindClause() != nil {
		i = ctx.UnwindClause()
	}

	switch v.parser.RuleNames[i.GetRuleIndex()] {
	case "matchClause":
		readingClause.Type = ast.ReadingClauseMatch
		readingClause.Match = ctx.MatchClause().Accept(v).(*ast.MatchClause)

	case "unwindClause":
		readingClause.Type = ast.ReadingClauseUnwind
		readingClause.Unwind = ctx.UnwindClause().Accept(v).(*ast.UnwindClause)
	}
	return readingClause
}

func (v *visitor) VisitMatchClause(ctx *MatchClauseContext) interface{} {
	match := &ast.MatchClause{}
	if ctx.OPTIONAL() != nil {
		match.Optional = true
	}
	if ctx.WhereClause() != nil {
		match.WithWhere = true
		match.Where = ctx.WhereClause().Accept(v).(*ast.Expr)
	}
	match.Pattern = ctx.Pattern().Accept(v).(*ast.Pattern)
	return match
}

func (v *visitor) VisitUnwindClause(ctx *UnwindClauseContext) interface{} {
	unwind := &ast.UnwindClause{}

	unwind.Expr = ctx.Expr().Accept(v).(*ast.Expr)
	unwind.Variable = ctx.Variable().Accept(v).(*ast.SymbolicNameNode)

	return unwind
}

func (v *visitor) VisitWhereClause(ctx *WhereClauseContext) interface{} {
	where := ctx.Expr().Accept(v).(*ast.Expr)
	return where
}

func (v *visitor) VisitPattern(ctx *PatternContext) interface{} {
	// TODO
	pattern := &ast.Pattern{}
	return pattern
}

func (v *visitor) VisitExpr(ctx *ExprContext) interface{} {
	// TODO
	expr := &ast.Expr{}
	return expr
}
