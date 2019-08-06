package parser

import (
	"fmt"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/leiysky/parser/ast"
)

var (
	_ CypherVisitor = &convertVisitor{}
)

type convertVisitor struct {
	parser *CypherParser
}

func (v *convertVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *convertVisitor) VisitChildren(node antlr.RuleNode) interface{} {
	return nil
}

func (v *convertVisitor) VisitTerminal(node antlr.TerminalNode) interface{} { return nil }

func (v *convertVisitor) VisitErrorNode(node antlr.ErrorNode) interface{} { return nil }

func (v *convertVisitor) VisitCypher(ctx *CypherContext) interface{} {
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

func (v *convertVisitor) VisitStmt(ctx *StmtContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitQuery(ctx *QueryContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitRegularQuery(ctx *RegularQueryContext) interface{} {
	regularQuery := &ast.RegularQueryStmt{}
	regularQuery.SingleQuery = ctx.SingleQuery().Accept(v).(*ast.SingleQueryStmt)

	var unionClauses []*ast.UnionClause
	for _, union := range ctx.AllUnionClause() {
		unionClauses = append(unionClauses, union.Accept(v).(*ast.UnionClause))
	}
	regularQuery.Unions = unionClauses
	return regularQuery
}

func (v *convertVisitor) VisitUnionClause(ctx *UnionClauseContext) interface{} {
	unionClause := &ast.UnionClause{}
	if ctx.ALL() != nil {
		unionClause.All = true
	}
	unionClause.SingleQuery = ctx.SingleQuery().Accept(v).(*ast.SingleQueryStmt)
	return unionClause
}

func (v *convertVisitor) VisitSingleQuery(ctx *SingleQueryContext) interface{} {
	singleQuery := &ast.SingleQueryStmt{}
	if ctx.SinglePartQuery() != nil {
		singleQuery.Type = ast.SingleQueryStmtSinglePart
		singleQuery.SinglePart = ctx.SinglePartQuery().Accept(v).(*ast.SinglePartQueryStmt)
	} else if ctx.MultiPartQuery() != nil {
		singleQuery.Type = ast.SingleQueryStmtMultiPart
		singleQuery.MultiPart = ctx.MultiPartQuery().Accept(v).(*ast.MultiPartQueryStmt)
	}
	return singleQuery
}

func (v *convertVisitor) VisitSinglePartQuery(ctx *SinglePartQueryContext) interface{} {
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

func (v *convertVisitor) VisitMultiPartQuery(ctx *MultiPartQueryContext) interface{} {
	multiPartQuery := &ast.MultiPartQueryStmt{}
	var parts []*ast.MultiPartQueryPartial
	var part *ast.MultiPartQueryPartial = &ast.MultiPartQueryPartial{}
	var readingCount int
	var updatingCount int
	var withCount int
	for _, child := range ctx.GetChildren() {
		if n, ok := child.GetPayload().(antlr.RuleContext); ok {
			switch v.parser.RuleNames[n.GetRuleIndex()] {
			case "readingClause":
				part.Readings = append(part.Readings, ctx.ReadingClause(readingCount).Accept(v).(*ast.ReadingClause))
				readingCount++
			case "updatingClause":
				part.Updatings = append(part.Updatings, ctx.UpdatingClause(updatingCount).Accept(v).(*ast.UpdatingClause))
				updatingCount++
			case "withClause":
				part.With = ctx.WithClause(withCount).Accept(v).(*ast.WithClause)
				withCount++
				parts = append(parts, part)
				part = &ast.MultiPartQueryPartial{}
			}
		}
	}
	multiPartQuery.MultiPart = parts
	multiPartQuery.SinglePart = ctx.SinglePartQuery().Accept(v).(*ast.SinglePartQueryStmt)
	return multiPartQuery
}

func (v *convertVisitor) VisitWithClause(ctx *WithClauseContext) interface{} {
	withClause := &ast.WithClause{}
	if ctx.DISTINCT() != nil {
		withClause.Distinct = true
	}
	withClause.ReturnBody = ctx.ReturnBody().Accept(v).(*ast.ReturnBody)
	if ctx.WhereClause() != nil {
		withClause.Where = ctx.WhereClause().Accept(v).(*ast.Expr)
	}
	return withClause
}

func (v *convertVisitor) VisitReadingClause(ctx *ReadingClauseContext) interface{} {
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

func (v *convertVisitor) VisitMatchClause(ctx *MatchClauseContext) interface{} {
	match := &ast.MatchClause{}
	if ctx.OPTIONAL() != nil {
		match.Optional = true
	}
	if ctx.WhereClause() != nil {
		match.Where = ctx.WhereClause().Accept(v).(*ast.Expr)
	}
	match.Pattern = ctx.Pattern().Accept(v).(*ast.Pattern)
	return match
}

func (v *convertVisitor) VisitUnwindClause(ctx *UnwindClauseContext) interface{} {
	unwind := &ast.UnwindClause{}
	unwind.Expr = ctx.Expr().Accept(v).(*ast.Expr)
	unwind.Variable = ctx.Variable().Accept(v).(*ast.SymbolicNameNode)
	return unwind
}

func (v *convertVisitor) VisitUpdatingClause(ctx *UpdatingClauseContext) interface{} {
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

func (v *convertVisitor) VisitCreateClause(ctx *CreateClauseContext) interface{} {
	create := &ast.CreateClause{}
	create.Pattern = ctx.Pattern().Accept(v).(*ast.Pattern)
	return create
}

func (v *convertVisitor) VisitSetClause(ctx *SetClauseContext) interface{} {
	set := &ast.SetClause{}
	var items []*ast.SetItemStmt
	for _, item := range ctx.AllSetItem() {
		items = append(items, item.Accept(v).(*ast.SetItemStmt))
	}
	set.SetItems = items
	return set
}

func (v *convertVisitor) VisitSetItem(ctx *SetItemContext) interface{} {
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

func (v *convertVisitor) VisitDeleteClause(ctx *DeleteClauseContext) interface{} {
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

func (v *convertVisitor) VisitRemoveClause(ctx *RemoveClauseContext) interface{} {
	removeClause := &ast.RemoveClause{}
	var items []*ast.RemoveItemStmt
	for _, item := range ctx.AllRemoveItem() {
		items = append(items, item.Accept(v).(*ast.RemoveItemStmt))
	}
	removeClause.RemoveItems = items
	return removeClause
}

func (v *convertVisitor) VisitRemoveItem(ctx *RemoveItemContext) interface{} {
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

func (v *convertVisitor) VisitMergeClause(ctx *MergeClauseContext) interface{} {
	mergeClause := &ast.MergeClause{}
	mergeClause.PatternPart = ctx.PatternPart().Accept(v).(*ast.PatternPart)
	var actions []*ast.MergeAction
	for _, action := range ctx.AllMergeAction() {
		actions = append(actions, action.Accept(v).(*ast.MergeAction))
	}
	mergeClause.MergeActions = actions
	return mergeClause
}

func (v *convertVisitor) VisitMergeAction(ctx *MergeActionContext) interface{} {
	mergeAction := &ast.MergeAction{}
	if ctx.CREATE() != nil {
		mergeAction.Type = ast.MergeActionCreate
	} else if ctx.MATCH() != nil {
		mergeAction.Type = ast.MergeActionMatch
	}
	mergeAction.Set = ctx.SetClause().Accept(v).(*ast.SetClause)
	return mergeAction
}

func (v *convertVisitor) VisitWhereClause(ctx *WhereClauseContext) interface{} {
	where := ctx.Expr().Accept(v).(*ast.Expr)
	return where
}

func (v *convertVisitor) VisitVariable(ctx *VariableContext) interface{} {
	variable := &ast.VariableNode{}
	variable.SymbolicName = ctx.SymbolicName().Accept(v).(*ast.SymbolicNameNode)
	return variable
}

func (v *convertVisitor) VisitSymbolicName(ctx *SymbolicNameContext) interface{} {
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

func (v *convertVisitor) VisitNodeLabels(ctx *NodeLabelsContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitNodeLabel(ctx *NodeLabelContext) interface{} {
	nodeLabel := &ast.NodeLabelNode{}
	nodeLabel.LabelName = ctx.LabelName().Accept(v).(*ast.SchemaNameNode)
	return nodeLabel
}

func (v *convertVisitor) VisitLabelName(ctx *LabelNameContext) interface{} {
	labelName := &ast.SchemaNameNode{}
	labelName.Type = ast.SchemaNameSymbolicName
	return labelName
}

func (v *convertVisitor) VisitPattern(ctx *PatternContext) interface{} {
	pattern := &ast.Pattern{}
	var parts []*ast.PatternPart
	for _, part := range ctx.AllPatternPart() {
		parts = append(parts, part.Accept(v).(*ast.PatternPart))
	}
	pattern.Parts = parts
	return pattern
}

func (v *convertVisitor) VisitPatternPart(ctx *PatternPartContext) interface{} {
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

func (v *convertVisitor) VisitAnonymousPatternPart(ctx *AnonymousPatternPartContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitPatternElement(ctx *PatternElementContext) interface{} {
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

func (v *convertVisitor) VisitNodePattern(ctx *NodePatternContext) interface{} {
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

func (v *convertVisitor) VisitPatternElementChain(ctx *PatternElementChainContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitRelationshipPattern(ctx *RelationshipPatternContext) interface{} {
	relationshipPattern := &ast.RelationshipPattern{}
	if ctx.LeftArrowHead() != nil && ctx.RightArrowHead() != nil {
		relationshipPattern.Type = ast.RelationshipBoth
	} else if ctx.LeftArrowHead() != nil {
		relationshipPattern.Type = ast.RelationshipLeft
	} else if ctx.RightArrowHead() != nil {
		relationshipPattern.Type = ast.RelationshipRight
	} else {
		relationshipPattern.Type = ast.RelationshipNone
	}
	if ctx.RelationshipDetail() != nil {
		relationshipPattern.Detail = ctx.RelationshipDetail().Accept(v).(*ast.RelationshipDetail)
	}
	return relationshipPattern
}

func (v *convertVisitor) VisitLeftArrowHead(ctx *LeftArrowHeadContext) interface{} {
	panic("Shouln't be invoked")
}

func (v *convertVisitor) VisitRightArrowHead(ctx *RightArrowHeadContext) interface{} {
	panic("Shouln't be invoked")
}

func (v *convertVisitor) VisitDash(ctx *DashContext) interface{} {
	panic("Shouln't be invoked")
}

func (v *convertVisitor) VisitRelationshipDetail(ctx *RelationshipDetailContext) interface{} {
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
		if rangeLiteral.MinHops() != nil {
			relationshipDetail.MinHops = rangeLiteral.MinHops().Accept(v).(*MinHopsContext).IntegerLiteral().Accept(v).(int)
		} else {
			relationshipDetail.MinHops = -1
		}
		if rangeLiteral.MaxHops() != nil {
			relationshipDetail.MaxHops = rangeLiteral.MaxHops().Accept(v).(*MaxHopsContext).IntegerLiteral().Accept(v).(int)
		} else {
			relationshipDetail.MaxHops = -1
		}
	} else {
		relationshipDetail.MinHops = 1
		relationshipDetail.MaxHops = 1
	}

	if ctx.Properties() != nil {
		relationshipDetail.Properties = ctx.Properties().Accept(v).(*ast.Properties)
	}
	return relationshipDetail
}

func (v *convertVisitor) VisitRangeLiteral(ctx *RangeLiteralContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitMinHops(ctx *MinHopsContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitMaxHops(ctx *MaxHopsContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitRelationshipTypes(ctx *RelationshipTypesContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitRelTypeName(ctx *RelTypeNameContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitSchemaName(ctx *SchemaNameContext) interface{} {
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

func (v *convertVisitor) VisitReturnClause(ctx *ReturnClauseContext) interface{} {
	returnClause := &ast.ReturnClause{}
	if ctx.DISTINCT() != nil {
		returnClause.Distinct = true
	}
	returnClause.ReturnBody = ctx.ReturnBody().Accept(v).(*ast.ReturnBody)
	return returnClause
}

func (v *convertVisitor) VisitReturnBody(ctx *ReturnBodyContext) interface{} {
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

func (v *convertVisitor) VisitReturnItems(ctx *ReturnItemsContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitReturnItem(ctx *ReturnItemContext) interface{} {
	returnItem := &ast.ReturnItem{}
	returnItem.Expr = ctx.Expr().Accept(v).(*ast.Expr)
	if ctx.AS() != nil {
		returnItem.As = true
		returnItem.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
	}
	return returnItem
}

func (v *convertVisitor) VisitOrderClause(ctx *OrderClauseContext) interface{} {
	orderClause := &ast.OrderClause{}
	var sortItems []*ast.SortItem
	for _, item := range ctx.AllSortItem() {
		sortItems = append(sortItems, item.Accept(v).(*ast.SortItem))
	}
	return orderClause
}

func (v *convertVisitor) VisitSkipClause(ctx *SkipClauseContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitLimitClause(ctx *LimitClauseContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitSortItem(ctx *SortItemContext) interface{} {
	sortItem := &ast.SortItem{}
	sortItem.Expr = ctx.Expr().Accept(v).(*ast.Expr)
	if ctx.ASC() != nil || ctx.ASCENDING() != nil {
		sortItem.Type = ast.SortAscending
	} else if ctx.DESC() != nil || ctx.DESCENDING() != nil {
		sortItem.Type = ast.SortDescending
	}
	return sortItem
}

func (v *convertVisitor) VisitExpr(ctx *ExprContext) interface{} {
	expr := &ast.Expr{}
	expr.OrExpr = ctx.OrExpr().Accept(v).(*ast.OrExpr)
	return expr
}

func (v *convertVisitor) VisitOrExpr(ctx *OrExprContext) interface{} {
	orExpr := &ast.OrExpr{}
	var exprs []*ast.XorExpr
	for _, expr := range ctx.AllXorExpr() {
		exprs = append(exprs, expr.Accept(v).(*ast.XorExpr))
	}
	orExpr.XorExprs = exprs
	return orExpr
}

func (v *convertVisitor) VisitXorExpr(ctx *XorExprContext) interface{} {
	xorExpr := &ast.XorExpr{}
	var exprs []*ast.AndExpr
	for _, expr := range ctx.AllAndExpr() {
		exprs = append(exprs, expr.Accept(v).(*ast.AndExpr))
	}
	xorExpr.AndExprs = exprs
	return xorExpr
}

func (v *convertVisitor) VisitAndExpr(ctx *AndExprContext) interface{} {
	andExpr := &ast.AndExpr{}
	var exprs []*ast.NotExpr
	for _, expr := range ctx.AllNotExpr() {
		exprs = append(exprs, expr.Accept(v).(*ast.NotExpr))
	}
	andExpr.NotExprs = exprs
	return andExpr
}

func (v *convertVisitor) VisitNotExpr(ctx *NotExprContext) interface{} {
	notExpr := &ast.NotExpr{}
	notExpr.ComparisonExpr = ctx.ComparisonExpr().Accept(v).(*ast.ComparisonExpr)
	return notExpr
}

func (v *convertVisitor) VisitComparisonExpr(ctx *ComparisonExprContext) interface{} {
	comparisonExpr := &ast.ComparisonExpr{}
	comparisonExpr.AddSubExpr = ctx.AddOrSubtractExpr().Accept(v).(*ast.AddSubExpr)
	var exprs []*ast.PartialComparisonExpr
	for _, expr := range ctx.AllPartialComparisonExpr() {
		exprs = append(exprs, expr.Accept(v).(*ast.PartialComparisonExpr))
	}
	comparisonExpr.PartialComparisonExprs = exprs
	return comparisonExpr
}

func (v *convertVisitor) VisitPartialComparisonExpr(ctx *PartialComparisonExprContext) interface{} {
	partialComparisonExpr := &ast.PartialComparisonExpr{}
	comp := ctx.GetChild(0).(antlr.TerminalNode)
	switch comp.GetText() {
	case "=":
		partialComparisonExpr.Type = ast.PartialComparisonEQ
	case "<>":
		partialComparisonExpr.Type = ast.PartialComparisonNE
	case "<":
		partialComparisonExpr.Type = ast.PartialComparisonLT
	case ">":
		partialComparisonExpr.Type = ast.PartialComparisonGT
	case "<=":
		partialComparisonExpr.Type = ast.PartialComparisonLTE
	case ">=":
		partialComparisonExpr.Type = ast.PartialComparisonGTE
	}
	partialComparisonExpr.AddSubExpr = ctx.AddOrSubtractExpr().Accept(v).(*ast.AddSubExpr)
	return partialComparisonExpr
}

func (v *convertVisitor) VisitAddOrSubtractExpr(ctx *AddOrSubtractExprContext) interface{} {
	addSubExpr := &ast.AddSubExpr{}
	exprs := ctx.AllMultiplyDivideModuloExpr()
	addSubExpr.LExpr = exprs[0].Accept(v).(*ast.MulDivModExpr)
	exprs = exprs[1:]
	var ops []ast.OpType
	for _, child := range ctx.GetChildren() {
		if n, ok := child.GetPayload().(*antlr.CommonToken); ok && n.GetTokenType() != CypherLexerSP {
			switch n.GetText() {
			case "+":
				ops = append(ops, ast.OpAdd)
			case "-":
				ops = append(ops, ast.OpSub)
			}
		}
	}
	var mulDivModExprs []*ast.MulDivModExpr
	for _, expr := range exprs {
		mulDivModExprs = append(mulDivModExprs, expr.Accept(v).(*ast.MulDivModExpr))
	}
	addSubExpr.RExprs = mulDivModExprs
	addSubExpr.Ops = ops
	return addSubExpr
}

func (v *convertVisitor) VisitMultiplyDivideModuloExpr(ctx *MultiplyDivideModuloExprContext) interface{} {
	mulDivModExpr := &ast.MulDivModExpr{}
	exprs := ctx.AllPowerOfExpr()
	mulDivModExpr.LExpr = exprs[0].Accept(v).(*ast.PowerOfExpr)
	exprs = exprs[1:]
	var ops []ast.OpType
	for _, child := range ctx.GetChildren() {
		if n, ok := child.GetPayload().(*antlr.CommonToken); ok && n.GetTokenType() != CypherLexerSP {
			switch n.GetText() {
			case "*":
				ops = append(ops, ast.OpMul)
			case "/":
				ops = append(ops, ast.OpDiv)
			case "%":
				ops = append(ops, ast.OpMod)
			}
		}
	}
	var powerOfExprs []*ast.PowerOfExpr
	for _, expr := range exprs {
		powerOfExprs = append(powerOfExprs, expr.Accept(v).(*ast.PowerOfExpr))
	}
	mulDivModExpr.RExprs = powerOfExprs
	mulDivModExpr.Ops = ops
	return mulDivModExpr
}

func (v *convertVisitor) VisitPowerOfExpr(ctx *PowerOfExprContext) interface{} {
	powerOfExpr := &ast.PowerOfExpr{}
	var exprs []*ast.UnaryAddSubExpr
	for _, expr := range ctx.AllUnaryAddOrSubtractExpr() {
		exprs = append(exprs, expr.Accept(v).(*ast.UnaryAddSubExpr))
	}
	powerOfExpr.UnaryAddSubExprs = exprs
	return powerOfExpr
}

func (v *convertVisitor) VisitUnaryAddOrSubtractExpr(ctx *UnaryAddOrSubtractExprContext) interface{} {
	unaryAddSub := &ast.UnaryAddSubExpr{}
	var ops []string
	for _, child := range ctx.GetChildren() {
		if n, ok := child.GetPayload().(*antlr.CommonToken); ok && n.GetTokenType() != CypherLexerSP {
			ops = append(ops, n.GetText())
		}
	}
	var isNegative bool
	for _, op := range ops {
		if op == "-" {
			isNegative = !isNegative
		}
	}
	unaryAddSub.IsNegative = isNegative
	return unaryAddSub
}

func (v *convertVisitor) VisitStringListNullOperatorExpr(ctx *StringListNullOperatorExprContext) interface{} {
	stringListNullExpr := &ast.StringListNullExpr{}
	// TODO
	return stringListNullExpr
}

func (v *convertVisitor) VisitListOperatorExpr(ctx *ListOperatorExprContext) interface{} {
	listOperatorExpr := &ast.ListOperationExpr{}
	if ctx.PropertyOrLabelsExpr() != nil {
		listOperatorExpr.InExpr = ctx.PropertyOrLabelsExpr().Accept(v).(*ast.PropertyOrLabelsExpr)
	} else if len(ctx.GetTokens(12)) > 0 {
		listOperatorExpr.RangeExprs[0] = ctx.Expr(0).Accept(v).(*ast.Expr)
		listOperatorExpr.RangeExprs[1] = ctx.Expr(1).Accept(v).(*ast.Expr)
	} else if len(ctx.GetTokens(12)) == 0 {
		expr := ctx.Expr(0).Accept(v).(*ast.Expr)
		listOperatorExpr.SingleExpr = expr
	}
	return listOperatorExpr
}

func (v *convertVisitor) VisitStringOperatorExpr(ctx *StringOperatorExprContext) interface{} {
	stringOperatorExpr := &ast.StringOperationExpr{}
	if ctx.STARTS() != nil {
		stringOperatorExpr.Type = ast.StringOperationStartsWith
	} else if ctx.ENDS() != nil {
		stringOperatorExpr.Type = ast.StringOperationEndsWith
	} else if ctx.CONTAINS() != nil {
		stringOperatorExpr.Type = ast.StringOperationContains
	}
	stringOperatorExpr.PropertyOrLabelsExpr = ctx.PropertyOrLabelsExpr().Accept(v).(*ast.PropertyOrLabelsExpr)
	return stringOperatorExpr
}

func (v *convertVisitor) VisitNullOperatorExpr(ctx *NullOperatorExprContext) interface{} {
	nullOperatorExpr := &ast.NullOperationExpr{}
	if ctx.NOT() == nil {
		nullOperatorExpr.IsIsNull = true
	}
	return nullOperatorExpr
}

func (v *convertVisitor) VisitPropertyOrLabelsExpr(ctx *PropertyOrLabelsExprContext) interface{} {
	propertyOrLabelsExpr := &ast.PropertyOrLabelsExpr{}
	propertyOrLabelsExpr.Atom = ctx.Atom().Accept(v).(*ast.Atom)
	if ctx.NodeLabels() != nil {
		var labels []*ast.NodeLabelNode
		for _, label := range ctx.NodeLabels().Accept(v).(*NodeLabelsContext).AllNodeLabel() {
			labels = append(labels, label.Accept(v).(*ast.NodeLabelNode))
		}
		propertyOrLabelsExpr.NodeLabels = labels
	}
	var lookups []*ast.PropertyLookup
	for _, lookup := range ctx.AllPropertyLookup() {
		lookups = append(lookups, lookup.Accept(v).(*ast.PropertyLookup))
	}
	return propertyOrLabelsExpr
}

func (v *convertVisitor) VisitPropertyLookup(ctx *PropertyLookupContext) interface{} {
	propertyLookup := &ast.PropertyLookup{}
	propertyLookup.PropertyKey = ctx.PropertyKeyName().Accept(v).(*PropertyKeyNameContext).SchemaName().Accept(v).(*ast.SchemaNameNode)
	return propertyLookup
}

func (v *convertVisitor) VisitPropertyKeyName(ctx *PropertyKeyNameContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitCaseExpr(ctx *CaseExprContext) interface{} {
	caseExpr := &ast.CaseExpr{}
	var alts []*ast.CaseAlt
	for _, alt := range ctx.AllCaseAlternatives() {
		alts = append(alts, alt.Accept(v).(*ast.CaseAlt))
	}
	caseExpr.Alts = alts
	if ctx.ELSE() != nil {
		// means there are at least 1 Expr
		if len(ctx.AllExpr()) > 1 {
			caseExpr.Expr = ctx.Expr(0).Accept(v).(*ast.Expr)
			caseExpr.Else = ctx.Expr(1).Accept(v).(*ast.Expr)
		} else if len(ctx.AllExpr()) == 1 {
			caseExpr.Expr = ctx.Expr(0).Accept(v).(*ast.Expr)
		}
	} else if len(ctx.AllExpr()) > 0 {
		caseExpr.Expr = ctx.Expr(0).Accept(v).(*ast.Expr)
	}
	return caseExpr
}

func (v *convertVisitor) VisitCaseAlternatives(ctx *CaseAlternativesContext) interface{} {
	caseAlt := &ast.CaseAlt{}
	caseAlt.When = ctx.Expr(0).Accept(v).(*ast.Expr)
	caseAlt.Then = ctx.Expr(1).Accept(v).(*ast.Expr)
	return caseAlt
}

func (v *convertVisitor) VisitAtom(ctx *AtomContext) interface{} {
	atom := &ast.Atom{}
	if ctx.Literal() != nil {
		atom.Type = ast.AtomLiteral
		atom.Literal = ctx.Literal().Accept(v).(*ast.LiteralNode)
	} else if ctx.Parameter() != nil {
		atom.Type = ast.AtomParameter
		atom.Parameter = ctx.Parameter().Accept(v).(*ast.ParameterNode)
	} else if ctx.CaseExpr() != nil {
		atom.Type = ast.AtomCase
		atom.CaseExpr = ctx.CaseExpr().Accept(v).(*ast.CaseExpr)
	} else if ctx.COUNT() != nil {
		atom.Type = ast.AtomCount
	} else if ctx.ListComprehension() != nil {
		atom.Type = ast.AtomList
		atom.ListComprehension = ctx.ListComprehension().Accept(v).(*ast.ListComprehension)
	} else if ctx.PatternComprehension() != nil {
		atom.Type = ast.AtomPatternComprehension
		atom.PatternComprehension = ctx.PatternComprehension().Accept(v).(*ast.PatternComprehension)
	} else if ctx.ALL() != nil {
		atom.Type = ast.AtomAllFilter
		atom.FilterExpr = ctx.FilterExpr().Accept(v).(*ast.FilterExpr)
	} else if ctx.ANY() != nil {
		atom.Type = ast.AtomAnyFilter
		atom.FilterExpr = ctx.FilterExpr().Accept(v).(*ast.FilterExpr)
	} else if ctx.NONE() != nil {
		atom.Type = ast.AtomNoneFilter
		atom.FilterExpr = ctx.FilterExpr().Accept(v).(*ast.FilterExpr)
	} else if ctx.SINGLE() != nil {
		atom.Type = ast.AtomSingleFilter
		atom.FilterExpr = ctx.FilterExpr().Accept(v).(*ast.FilterExpr)
	} else if ctx.RelationshipsPattern() != nil {
		atom.Type = ast.AtomPattern
		atom.PatternElement = ctx.RelationshipsPattern().Accept(v).(*ast.PatternElement)
	} else if ctx.ParenthesizedExpr() != nil {
		atom.Type = ast.AtomParenthesizedExpr
		atom.ParenthesizedExpr = ctx.ParenthesizedExpr().Accept(v).(*ParenthesizedExprContext).Expr().Accept(v).(*ast.Expr)
	} else if ctx.FunctionInvocation() != nil {
		// TODO
		panic("FunctionInvocation not support now")
	} else if ctx.Variable() != nil {
		atom.Type = ast.AtomVariable
		atom.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
	}
	return atom
}

func (v *convertVisitor) VisitRelationshipsPattern(ctx *RelationshipsPatternContext) interface{} {
	patternElement := &ast.PatternElement{}
	// same as VisitPatternElement
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

func (v *convertVisitor) VisitParenthesizedExpr(ctx *ParenthesizedExprContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitFilterExpr(ctx *FilterExprContext) interface{} {
	filterExpr := &ast.FilterExpr{}
	idInColl := ctx.IdInColl().Accept(v).(*IdInCollContext)
	filterExpr.Variable = idInColl.Variable().Accept(v).(*ast.VariableNode)
	filterExpr.In = idInColl.Expr().Accept(v).(*ast.Expr)
	if ctx.WhereClause() != nil {
		filterExpr.Where = ctx.WhereClause().Accept(v).(*ast.Expr)
	}
	return filterExpr
}

func (v *convertVisitor) VisitIdInColl(ctx *IdInCollContext) interface{} {
	return ctx
}

func (v *convertVisitor) VisitListComprehension(ctx *ListComprehensionContext) interface{} {
	listComprehension := &ast.ListComprehension{}
	listComprehension.FilterExpr = ctx.FilterExpr().Accept(v).(*ast.FilterExpr)
	if ctx.Expr() != nil {
		listComprehension.Expr = ctx.Expr().Accept(v).(*ast.Expr)
	}
	return listComprehension
}

func (v *convertVisitor) VisitPatternComprehension(ctx *PatternComprehensionContext) interface{} {
	patternComprehension := &ast.PatternComprehension{}
	if ctx.Variable() != nil {
		patternComprehension.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
	}
	patternComprehension.PatternElement = ctx.RelationshipsPattern().Accept(v).(*ast.PatternElement)
	if ctx.WHERE() != nil {
		patternComprehension.Where = ctx.Expr(0).Accept(v).(*ast.Expr)
	}
	patternComprehension.Expr = ctx.AllExpr()[len(ctx.AllExpr())-1].Accept(v).(*ast.Expr)
	return patternComprehension
}

func (v *convertVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	literal := &ast.LiteralNode{}
	if ctx.NumberLiteral() != nil {
		literal.Type = ast.LiteralNumber
		literal.Number = ctx.NumberLiteral().Accept(v).(*ast.NumberLiteral)
	} else if ctx.StringLiteral() != nil {
		var err error
		literal.Type = ast.LiteralString
		literal.String, err = strconv.Unquote(ctx.StringLiteral().GetText())
		if err != nil {
			panic(err)
		}
	} else if ctx.BooleanLiteral() != nil {
		literal.Type = ast.LiteralBoolean
		if ctx.BooleanLiteral().GetText() == "TRUE" {
			literal.Boolean = true
		}
	} else if ctx.NULL() != nil {
		literal.Type = ast.LiteralNull
	} else if ctx.MapLiteral() != nil {
		literal.Type = ast.LiteralMap
		literal.Map = ctx.MapLiteral().Accept(v).(*ast.MapLiteral)
	} else if ctx.ListLiteral() != nil {
		literal.Type = ast.LiteralList
		literal.List = ctx.ListLiteral().Accept(v).(*ast.ListLiteral)
	}
	return literal
}

func (v *convertVisitor) VisitNumberLiteral(ctx *NumberLiteralContext) interface{} {
	numberLiteral := &ast.NumberLiteral{}
	if ctx.IntegerLiteral() != nil {
		numberLiteral.Type = ast.NumberLiteralInteger
		numberLiteral.Integer = ctx.IntegerLiteral().Accept(v).(int)
	} else if ctx.DoubleLiteral() != nil {
		numberLiteral.Type = ast.NumberLiteralDouble
		numberLiteral.Double = ctx.DoubleLiteral().Accept(v).(float64)
	}
	return numberLiteral
}

func (v *convertVisitor) VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{} {
	panic("Shouln't be invoked")
}

func (v *convertVisitor) VisitMapLiteral(ctx *MapLiteralContext) interface{} {
	mapLiteral := &ast.MapLiteral{}
	var keys []*ast.SchemaNameNode
	var exprs []*ast.Expr
	for i := range ctx.AllPropertyKeyName() {
		keys = append(keys, ctx.PropertyKeyName(i).Accept(v).(*PropertyKeyNameContext).SchemaName().Accept(v).(*ast.SchemaNameNode))
		exprs = append(exprs, ctx.Expr(i).Accept(v).(*ast.Expr))
	}
	mapLiteral.PropertyKeys = keys
	mapLiteral.Exprs = exprs
	return mapLiteral
}

func (v *convertVisitor) VisitListLiteral(ctx *ListLiteralContext) interface{} {
	listLiteral := &ast.ListLiteral{}
	var exprs []*ast.Expr
	for _, expr := range ctx.AllExpr() {
		exprs = append(exprs, expr.Accept(v).(*ast.Expr))
	}
	listLiteral.Exprs = exprs
	return listLiteral
}

func (v *convertVisitor) VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{} {
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

func (v *convertVisitor) VisitDoubleLiteral(ctx *DoubleLiteralContext) interface{} {
	var value float64
	if ctx.RegularDecimalReal() != nil {
		regular := ctx.RegularDecimalReal().GetSymbol().GetText()
		v, err := strconv.ParseFloat(regular, 64)
		if err != nil {
			panic(err)
		}
		value = v
	} else if ctx.ExponentDecimalReal() != nil {
		exponent := ctx.ExponentDecimalReal().GetSymbol().GetText()
		v, err := strconv.ParseFloat(exponent, 64)
		if err != nil {
			panic(err)
		}
		value = v
	}
	return value
}

func (v *convertVisitor) VisitParameter(ctx *ParameterContext) interface{} {
	parameter := &ast.ParameterNode{}
	if ctx.SymbolicName() != nil {
		parameter.Type = ast.ParameterSymbolicname
		parameter.SymbolicName = ctx.SymbolicName().Accept(v).(*ast.SymbolicNameNode)
	} else if ctx.DecimalInteger() != nil {
		parameter.Type = ast.ParameterDecimalInteger
		v, err := strconv.Atoi(ctx.DecimalInteger().GetSymbol().GetText())
		if err != nil {
			panic(err)
		}
		parameter.DecimalInteger = v
	}
	return parameter
}

func (v *convertVisitor) VisitProperties(ctx *PropertiesContext) interface{} {
	properties := &ast.Properties{}
	if ctx.MapLiteral() != nil {
		properties.Type = ast.PropertiesMapLiteral
		properties.MapLiteral = ctx.MapLiteral().Accept(v).(*ast.MapLiteral)
	} else if ctx.Parameter() != nil {
		properties.Type = ast.PropertiesParameter
		properties.Parameter = ctx.Parameter().Accept(v).(*ast.ParameterNode)
	}
	return properties
}

func (v *convertVisitor) VisitPropertyExpr(ctx *PropertyExprContext) interface{} {
	propertyExpr := &ast.PropertyExpr{}
	propertyExpr.Atom = ctx.Atom().Accept(v).(*ast.Atom)
	var lookups []*ast.PropertyLookup
	for _, lookup := range ctx.AllPropertyLookup() {
		lookups = append(lookups, lookup.Accept(v).(*ast.PropertyLookup))
	}
	propertyExpr.Lookups = lookups
	return propertyExpr
}

func (v *convertVisitor) VisitReservedWord(ctx *ReservedWordContext) interface{} {
	reservedWord := &ast.ReservedWordNode{}
	reservedWord.Content = ctx.GetText()
	return reservedWord
}

func (v *convertVisitor) VisitStandaloneCall(ctx *StandaloneCallContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *convertVisitor) VisitYieldItems(ctx *YieldItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *convertVisitor) VisitYieldItem(ctx *YieldItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *convertVisitor) VisitFunctionInvocation(ctx *FunctionInvocationContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *convertVisitor) VisitFunctionName(ctx *FunctionNameContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *convertVisitor) VisitExplicitProcedureInvocation(ctx *ExplicitProcedureInvocationContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *convertVisitor) VisitImplicitProcedureInvocation(ctx *ImplicitProcedureInvocationContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *convertVisitor) VisitProcedureResultField(ctx *ProcedureResultFieldContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *convertVisitor) VisitProcedureName(ctx *ProcedureNameContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *convertVisitor) VisitNamespace(ctx *NamespaceContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *convertVisitor) VisitInQueryCall(ctx *InQueryCallContext) interface{} {
	// TODO
	panic("not implemented")
}
