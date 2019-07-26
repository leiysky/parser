package parser

import (
	"fmt"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/leiysky/parser/ast"
)

var (
// _ CypherVisitor = &visitor{}
)

type visitor struct {
	CypherVisitor
	parser *CypherParser
}

func (v *visitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

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
	var updatingClauses []*ast.UpdatingClause
	for _, c := range ctx.AllUpdatingClause() {
		updatingClauses = append(updatingClauses, c.Accept(v).(*ast.UpdatingClause))
	}
	singlePartQuery.ReadingClauses = readingClauses
	singlePartQuery.UpdatingClauses = updatingClauses
	if ctx.ReturnClause() != nil {
		singlePartQuery.Return = ctx.ReturnClause().Accept(v).(*ast.ReturnClause)
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

func (v *visitor) VisitUpdatingClause(ctx *UpdatingClauseContext) interface{} {
	updatingClause := &ast.UpdatingClause{}

	var i antlr.RuleContext
	if i = ctx.CreateClause(); i != nil {
		updatingClause.Type = ast.UpdatingClauseCreate
		updatingClause.Create = ctx.CreateClause().Accept(v).(*ast.CreateClause)
	} else if i = ctx.MergeClause(); i != nil {
		updatingClause.Type = ast.UpdatingClauseMerge
		updatingClause.Merge = ctx.MergeClause().Accept(v).(*ast.MergeClause)
	} else if i = ctx.SetClause(); i != nil {
		updatingClause.Type = ast.UpdatingClauseSet
		updatingClause.Set = ctx.SetClause().Accept(v).(*ast.SetClause)
	} else if i = ctx.DeleteClause(); i != nil {
		updatingClause.Type = ast.UpdatingClauseDelete
		updatingClause.Delete = ctx.DeleteClause().Accept(v).(*ast.DeleteClause)
	} else if i = ctx.RemoveClause(); i != nil {
		updatingClause.Type = ast.UpdatingClauseRemove
		updatingClause.Remove = ctx.RemoveClause().Accept(v).(*ast.RemoveClause)
	}
	return updatingClause
}

func (v *visitor) VisitCreateClause(ctx *CreateClauseContext) interface{} {
	create := &ast.CreateClause{}
	create.Pattern = ctx.Pattern().Accept(v).(*ast.Pattern)
	return create
}

func (v *visitor) VisitSetClause(ctx *SetClauseContext) interface{} {
	set := &ast.SetClause{}
	var items []*ast.SetItemStmt
	for _, item := range ctx.AllSetItem() {
		items = append(items, item.Accept(v).(*ast.SetItemStmt))
	}
	set.SetItems = items
	return set
}

func (v *visitor) VisitSetItem(ctx *SetItemContext) interface{} {
	setItem := &ast.SetItemStmt{}

	if ctx.PropertyExpr() != nil {
		setItem.Type = ast.SetItemProperty
		setItem.Property = ctx.PropertyExpr().Accept(v).(*ast.PropertyExpr)
	} else if len(ctx.GetTokens(3)) > 0 {
		// 3 presents '=' token, see Cypher.tokens
		setItem.Type = ast.SetItemVariableAssignment
		setItem.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
		setItem.Expr = ctx.Expr().Accept(v).(*ast.Expr)
	} else if len(ctx.GetTokens(4)) > 0 {
		// 4 presents '+=' token, see Cypher.tokens
		setItem.Type = ast.SetItemVariableIncrement
		setItem.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
		setItem.Expr = ctx.Expr().Accept(v).(*ast.Expr)
	} else if allLabel := ctx.NodeLabels().Accept(v).(*NodeLabelsContext).AllNodeLabel(); len(allLabel) > 0 {
		var labels []*ast.NodeLabelNode
		for _, label := range allLabel {
			labels = append(labels, label.Accept(v).(*ast.NodeLabelNode))
		}
		setItem.Type = ast.SetItemVariableLabel
		setItem.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
		setItem.Labels = labels
	}
	return setItem
}

func (v *visitor) VisitDeleteClause(ctx *DeleteClauseContext) interface{} {
	deleteClause := &ast.DeleteClause{}
	if ctx.DETACH() != nil {
		deleteClause.Detach = true
	}
	var exprs []*ast.Expr
	for _, expr := range ctx.AllExpr() {
		exprs = append(exprs, expr.Accept(v).(*ast.Expr))
	}
	deleteClause.Exprs = exprs
	return deleteClause
}

func (v *visitor) VisitRemoveClause(ctx *RemoveClauseContext) interface{} {
	removeClause := &ast.RemoveClause{}
	var items []*ast.RemoveItemStmt
	for _, item := range ctx.AllRemoveItem() {
		items = append(items, item.Accept(v).(*ast.RemoveItemStmt))
	}
	removeClause.RemoveItems = items
	return removeClause
}

func (v *visitor) VisitRemoveItem(ctx *RemoveItemContext) interface{} {
	removeItem := &ast.RemoveItemStmt{}
	if ctx.Variable() != nil {
		removeItem.Type = ast.RemoveItemVariable
		removeItem.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
		var labels []*ast.NodeLabelNode
		for _, label := range ctx.NodeLabels().Accept(v).(*NodeLabelsContext).AllNodeLabel() {
			labels = append(labels, label.Accept(v).(*ast.NodeLabelNode))
		}
		removeItem.Labels = labels
	} else if ctx.PropertyExpr() != nil {
		removeItem.Type = ast.RemoveItemProperty
		removeItem.Property = ctx.Variable().Accept(v).(*ast.PropertyExpr)
	}
	return removeItem
}

func (v *visitor) VisitWhereClause(ctx *WhereClauseContext) interface{} {
	where := ctx.Expr().Accept(v).(*ast.Expr)
	return where
}

func (v *visitor) VisitExpr(ctx *ExprContext) interface{} {
	// TODO
	expr := &ast.Expr{}
	return expr
}

func (v *visitor) VisitVariable(ctx *VariableContext) interface{} {
	variable := &ast.VariableNode{}
	variable.SymbolicName = ctx.SymbolicName().Accept(v).(*ast.SymbolicNameNode)
	return variable
}

func (v *visitor) VisitSymbolicName(ctx *SymbolicNameContext) interface{} {
	symbolicName := &ast.SymbolicNameNode{}

	var i antlr.TerminalNode
	if i = ctx.UnescapedSymbolicName(); i != nil {
	} else if i = ctx.EscapedSymbolicName(); i != nil {
	} else if i = ctx.HexLetter(); i != nil {
	} else if i = ctx.COUNT(); i != nil {
	} else if i = ctx.FILTER(); i != nil {
	} else if i = ctx.EXTRACT(); i != nil {
	} else if i = ctx.ANY(); i != nil {
	} else if i = ctx.NONE(); i != nil {
	} else if i = ctx.SINGLE(); i != nil {
	}
	switch i.GetSymbol().GetTokenType() {
	case CypherLexerUnescapedSymbolicName:
		symbolicName.Type = ast.SymbolicNameUnescaped
		symbolicName.Value = i.GetText()
	case CypherLexerEscapedSymbolicName:
		symbolicName.Type = ast.SymbolicNameEscaped
		symbolicName.Value = i.GetText()
	case CypherLexerHexLetter:
		symbolicName.Type = ast.SymbolicNameHexLetter
		symbolicName.Value = i.GetText()
	case CypherLexerCOUNT:
		symbolicName.Type = ast.SymbolicNameCount
	case CypherLexerFILTER:
		symbolicName.Type = ast.SymbolicNameFilter
	case CypherLexerEXTRACT:
		symbolicName.Type = ast.SymbolicNameExtract
	case CypherLexerANY:
		symbolicName.Type = ast.SymbolicNameAny
	case CypherLexerNONE:
		symbolicName.Type = ast.SymbolicNameNone
	case CypherLexerSINGLE:
		symbolicName.Type = ast.SymbolicNameSingle
	}
	return symbolicName
}

func (v *visitor) VisitNodeLabels(ctx *NodeLabelsContext) interface{} {
	return ctx
}

func (v *visitor) VisitNodeLabel(ctx *NodeLabelContext) interface{} {
	nodeLabel := &ast.NodeLabelNode{}
	nodeLabel.LabelName = ctx.LabelName().Accept(v).(*ast.SchemaNameNode)
	return nodeLabel
}

func (v *visitor) VisitLabelName(ctx *LabelNameContext) interface{} {
	labelName := &ast.SchemaNameNode{}
	labelName.Type = ast.SchemaNameSymbolicName
	return labelName
}

func (v *visitor) VisitPattern(ctx *PatternContext) interface{} {
	pattern := &ast.Pattern{}
	var parts []*ast.PatternPart
	for _, part := range ctx.AllPatternPart() {
		parts = append(parts, part.Accept(v).(*ast.PatternPart))
	}
	pattern.Parts = parts
	return pattern
}

func (v *visitor) VisitPatternPart(ctx *PatternPartContext) interface{} {
	patternPart := &ast.PatternPart{}
	if ctx.Variable() != nil {
		patternPart.WithVariable = true
		patternPart.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
	}
	patternPart.Element = ctx.
		AnonymousPatternPart().
		Accept(v).(*AnonymousPatternPartContext).
		PatternElement().
		Accept(v).(*ast.PatternElement)
	return patternPart
}

func (v *visitor) VisitAnonymousPatternPart(ctx *AnonymousPatternPartContext) interface{} {
	return ctx
}

func (v *visitor) VisitPatternElement(ctx *PatternElementContext) interface{} {
	patternElement := &ast.PatternElement{}
	// strip useless parenthesises recursively
	if ctx.PatternElement() != nil {
		ctx = ctx.PatternElement().Accept(v).(*PatternElementContext)
		return ctx
	}

	patternElement.StartNode = ctx.NodePattern().Accept(v).(*ast.NodePattern)
	var nodes []*ast.NodePattern
	var relationships []*ast.RelationshipPattern
	for _, pair := range ctx.AllPatternElementChain() {
		p := pair.Accept(v).(*PatternElementChainContext)
		nodes = append(nodes, p.NodePattern().Accept(v).(*ast.NodePattern))
		relationships = append(relationships, p.RelationshipPattern().Accept(v).(*ast.RelationshipPattern))
	}
	patternElement.Nodes = nodes
	patternElement.Relationships = relationships
	return patternElement
}

func (v *visitor) VisitNodePattern(ctx *NodePatternContext) interface{} {
	nodePattern := &ast.NodePattern{}
	if ctx.Variable() != nil {
		nodePattern.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
	}
	var labels []*ast.NodeLabelNode
	if ctx.NodeLabels() != nil {
		labelsCtx := ctx.NodeLabels().Accept(v).(*NodeLabelsContext)
		fmt.Println(labelsCtx.AllNodeLabel())
		for _, label := range labelsCtx.AllNodeLabel() {
			labels = append(labels, label.Accept(v).(*ast.NodeLabelNode))
		}
	}
	nodePattern.Labels = labels
	if ctx.Properties() != nil {
		nodePattern.Properties = ctx.Properties().Accept(v).(*ast.Properties)
	}
	return nodePattern
}

func (v *visitor) VisitPatternElementChain(ctx *PatternElementChainContext) interface{} {
	return ctx
}

func (v *visitor) VisitRelationshipPattern(ctx *RelationshipPatternContext) interface{} {
	relationshipPattern := &ast.RelationshipPattern{}
	if len(ctx.GetTokens(19)) > 0 && len(ctx.GetTokens(20)) > 0 {
		// 19 represents '<', 20 represents '>', see Cypher.tokens
		relationshipPattern.Type = ast.RelationshipBoth
	} else if len(ctx.GetTokens(19)) > 0 {
		relationshipPattern.Type = ast.RelationshipLeft
	} else if len(ctx.GetTokens(20)) > 0 {
		relationshipPattern.Type = ast.RelationshipRight
	} else {
		relationshipPattern.Type = ast.RelationshipNone
	}
	relationshipPattern.Detail = ctx.RelationshipDetail().Accept(v).(*ast.RelationshipDetail)
	return relationshipPattern
}

func (v *visitor) VisitRelationshipDetail(ctx *RelationshipDetailContext) interface{} {
	relationshipDetail := &ast.RelationshipDetail{}
	if ctx.Variable() != nil {
		relationshipDetail.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
	}
	if ctx.RelationshipTypes() != nil {
		var relTypes []*ast.SchemaNameNode
		for _, relType := range ctx.RelationshipTypes().Accept(v).(*RelationshipTypesContext).AllRelTypeName() {
			relTypes = append(relTypes, relType.Accept(v).(*RelTypeNameContext).SchemaName().Accept(v).(*ast.SchemaNameNode))
		}
		relationshipDetail.RelationshipTypes = relTypes
	}

	if ctx.RangeLiteral() != nil {
		rangeLiteral := ctx.RangeLiteral().Accept(v).(*RangeLiteralContext)
		if len(ctx.GetTokens(12)) == 0 {
			// 12 represents '..', see Cypher.tokens
			// means there is no upper bound
			if len(rangeLiteral.AllIntegerLiteral()) == 0 {
				// condition: [*]
				// means there is no lower bound and no upper bound, so we make it [-1, -1]
				relationshipDetail.Range = [2]int{-1, -1}
			} else {
				// condition: [*n]
				// means there is lower bound but no upper bound, so we make it [n, n] (n represents lower bound)
				lower := rangeLiteral.AllIntegerLiteral()[0].Accept(v).(int)
				relationshipDetail.Range = [2]int{lower, lower}
			}
		} else {
			// means there is upper bound
			// there are 3 conditions: [*m..n] [*m..] [*..n]
			if len(rangeLiteral.AllIntegerLiteral()) == 2 {
				// condition: [*m..n]
				// means has both lower and upper bound, so we make it [m, n]
				lower := rangeLiteral.AllIntegerLiteral()[0].Accept(v).(int)
				upper := rangeLiteral.AllIntegerLiteral()[1].Accept(v).(int)
				relationshipDetail.Range = [2]int{lower, upper}
			} else if len(rangeLiteral.AllIntegerLiteral()) == 1 {
				// conditions: [*m..] [*..n]
				bound := rangeLiteral.AllIntegerLiteral()[0].Accept(v).(int)
				var isLower bool
				children := rangeLiteral.GetChildren()
				// if IntegerLiteral appears before '..', then we can know it's [*m..], other wise it's [*..n]
				for _, child := range children {
					switch child.GetPayload().(type) {
					case *IntegerLiteralContext:
						isLower = true
						goto Outloop
					case antlr.TerminalNode:
						goto Outloop
					}
				}
			Outloop:

				if isLower {
					relationshipDetail.Range = [2]int{bound, -1}
				} else {
					relationshipDetail.Range = [2]int{0, bound}
				}
			}
		}
	} // Parse RangeLiteral

	if ctx.Properties() != nil {
		relationshipDetail.Properties = ctx.Properties().Accept(v).(*ast.Properties)
	}
	return relationshipDetail
}

func (v *visitor) VisitRangeLiteral(ctx *RangeLiteralContext) interface{} {
	return ctx
}

func (v *visitor) VisitRelationshipTypes(ctx *RelationshipTypesContext) interface{} {
	return ctx
}

func (v *visitor) VisitRelTypeName(ctx *RelTypeNameContext) interface{} {
	return ctx
}

func (v *visitor) VisitSchemaName(ctx *SchemaNameContext) interface{} {
	schemaName := &ast.SchemaNameNode{}
	var i antlr.RuleContext
	if i = ctx.SymbolicName(); i != nil {
		schemaName.Type = ast.SchemaNameSymbolicName
		schemaName.SymbolicName = i.Accept(v).(*ast.SymbolicNameNode)
	} else if i = ctx.ReservedWord(); i != nil {
		schemaName.Type = ast.SchemaNameReservedWord
		schemaName.ReservedWord = i.Accept(v).(*ast.ReservedWordNode)
	}
	return schemaName
}

func (v *visitor) VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{} {
	var value int
	if ctx.HexInteger() != nil {
		hex := ctx.HexInteger().GetSymbol().GetText()
		v, err := strconv.ParseInt(hex, 0, 64)
		if err != nil {
			panic(err)
		}
		value = int(v)
	} else if ctx.OctalInteger() != nil {
		oct := ctx.OctalInteger().GetSymbol().GetText()
		v, err := strconv.ParseInt(oct, 0, 64)
		if err != nil {
			panic(err)
		}
		value = int(v)
	} else if ctx.DecimalInteger() != nil {
		dec := ctx.DecimalInteger().GetSymbol().GetText()
		v, err := strconv.ParseInt(dec, 0, 64)
		if err != nil {
			panic(err)
		}
		value = int(v)
	}

	return value
}

func (v *visitor) VisitReturnClause(ctx *ReturnClauseContext) interface{} {
	returnClause := &ast.ReturnClause{}
	if ctx.DISTINCT() != nil {
		returnClause.Distinct = true
	}
	returnClause.ReturnBody = ctx.ReturnBody().Accept(v).(*ast.ReturnBody)
	return returnClause
}

func (v *visitor) VisitReturnBody(ctx *ReturnBodyContext) interface{} {
	returnBody := &ast.ReturnBody{}
	var returnItems []*ast.ReturnItem
	for _, item := range ctx.ReturnItems().Accept(v).(*ReturnItemsContext).AllReturnItem() {
		returnItems = append(returnItems, item.Accept(v).(*ast.ReturnItem))
	}
	returnBody.ReturnItems = returnItems
	if ctx.OrderClause() != nil {
		returnBody.OrderBy = ctx.OrderClause().Accept(v).(*ast.OrderClause)
	}
	if ctx.SkipClause() != nil {
		returnBody.Skip = ctx.SkipClause().(*SkipClauseContext).Expr().Accept(v).(*ast.Expr)
	}
	if ctx.LimitClause() != nil {
		returnBody.Limit = ctx.LimitClause().(*LimitClauseContext).Expr().Accept(v).(*ast.Expr)
	}
	return returnBody
}

func (v *visitor) VisitReturnItems(ctx *ReturnItemsContext) interface{} {
	return ctx
}

func (v *visitor) VisitReturnItem(ctx *ReturnItemContext) interface{} {
	returnItem := &ast.ReturnItem{}
	returnItem.Expr = ctx.Expr().Accept(v).(*ast.Expr)
	if ctx.AS() != nil {
		returnItem.As = true
		returnItem.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
	}
	return returnItem
}

func (v *visitor) VisitOrderClause(ctx *OrderClauseContext) interface{} {
	orderClause := &ast.OrderClause{}
	var sortItems []*ast.SortItem
	for _, item := range ctx.AllSortItem() {
		sortItems = append(sortItems, item.Accept(v).(*ast.SortItem))
	}
	return orderClause
}

func (v *visitor) VisitSkipClause(ctx *SkipClauseContext) interface{} {
	return ctx
}

func (v *visitor) VisitLimitClause(ctx *LimitClauseContext) interface{} {
	return ctx
}

func (v *visitor) VisitSortItem(ctx *SortItemContext) interface{} {
	sortItem := &ast.SortItem{}
	sortItem.Expr = ctx.Expr().Accept(v).(*ast.Expr)
	if ctx.ASC() != nil || ctx.ASCENDING() != nil {
		sortItem.Type = ast.SortAscending
	} else if ctx.DESC() != nil || ctx.DESCENDING() != nil {
		sortItem.Type = ast.SortDescending
	}
	return sortItem
}
