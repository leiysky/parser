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

package parser

import (
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/leiysky/parser/ast"
)

var (
	_ CypherVisitor = &ConvertVisitor{}
)

type ConvertVisitor struct {
	parser *CypherParser
}

func NewConvertVisitor(parser *CypherParser) CypherVisitor {
	return &ConvertVisitor{
		parser: parser,
	}
}

func (v *ConvertVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *ConvertVisitor) VisitChildren(node antlr.RuleNode) interface{} {
	return nil
}

func (v *ConvertVisitor) VisitTerminal(node antlr.TerminalNode) interface{} { return nil }

func (v *ConvertVisitor) VisitErrorNode(node antlr.ErrorNode) interface{} { return nil }

func (v *ConvertVisitor) VisitCypher(ctx *CypherContext) interface{} {
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
		node.Query = query.RegularQuery().Accept(v).(*ast.QueryStmt)

	case "standaloneCall":
		node.Type = ast.CypherStmtStandaloneCall
		node.StandaloneCall = query.StandaloneCall().Accept(v).(*ast.StandaloneCall)
	}
	return node
}

func (v *ConvertVisitor) VisitStmt(ctx *StmtContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitQuery(ctx *QueryContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitRegularQuery(ctx *RegularQueryContext) interface{} {
	query := &ast.QueryStmt{}
	clauses := ctx.SingleQuery().Accept(v).([]ast.Stmt)
	for _, union := range ctx.AllUnionClause() {
		clauses = append(clauses, union.Accept(v).(*ast.UnionClause))
	}
	query.Clauses = clauses
	return query
}

func (v *ConvertVisitor) VisitUnionClause(ctx *UnionClauseContext) interface{} {
	unionClause := &ast.UnionClause{}
	if ctx.ALL() != nil {
		unionClause.All = true
	}
	unionClause.Clauses = ctx.SingleQuery().Accept(v).([]ast.Stmt)
	return unionClause
}

func (v *ConvertVisitor) VisitSingleQuery(ctx *SingleQueryContext) interface{} {
	var clauses []ast.Stmt
	if ctx.SinglePartQuery() != nil {
		ss := ctx.SinglePartQuery().Accept(v).([]ast.Stmt)
		clauses = append(clauses, ss...)
	} else if ctx.MultiPartQuery() != nil {
		ss := ctx.MultiPartQuery().Accept(v).([]ast.Stmt)
		clauses = append(clauses, ss...)
	}
	return clauses
}

func (v *ConvertVisitor) VisitSinglePartQuery(ctx *SinglePartQueryContext) interface{} {
	var clauses []ast.Stmt
	for _, c := range ctx.AllReadingClause() {
		clauses = append(clauses, c.Accept(v).(ast.Stmt))
	}
	for _, c := range ctx.AllUpdatingClause() {
		clauses = append(clauses, c.Accept(v).(ast.Stmt))
	}
	if ctx.ReturnClause() != nil {
		clauses = append(clauses, ctx.ReturnClause().Accept(v).(ast.Stmt))
	}
	return clauses
}

func (v *ConvertVisitor) VisitMultiPartQuery(ctx *MultiPartQueryContext) interface{} {
	var clauses []ast.Stmt
	for _, p := range ctx.AllMultiPartQueryPartial() {
		clauses = append(clauses, p.Accept(v).([]ast.Stmt)...)
	}
	clauses = append(clauses, ctx.SinglePartQuery().Accept(v).([]ast.Stmt)...)
	return clauses
}

func (v *ConvertVisitor) VisitMultiPartQueryPartial(ctx *MultiPartQueryPartialContext) interface{} {
	var clauses []ast.Stmt
	for _, r := range ctx.AllReadingClause() {
		clauses = append(clauses, r.Accept(v).(ast.Stmt))
	}
	for _, u := range ctx.AllUpdatingClause() {
		clauses = append(clauses, u.Accept(v).(ast.Stmt))
	}
	clauses = append(clauses, ctx.WithClause().Accept(v).(ast.Stmt))
	return clauses
}

func (v *ConvertVisitor) VisitWithClause(ctx *WithClauseContext) interface{} {
	withClause := &ast.WithClause{}
	if ctx.DISTINCT() != nil {
		withClause.Distinct = true
	}
	withClause.ReturnBody = ctx.ReturnBody().Accept(v).(*ast.ReturnBody)
	if ctx.WhereClause() != nil {
		withClause.Where = ctx.WhereClause().Accept(v).(ast.Expr)
	}
	return withClause
}

func (v *ConvertVisitor) VisitReadingClause(ctx *ReadingClauseContext) interface{} {
	var n ast.Stmt
	if ctx.MatchClause() != nil {
		n = ctx.MatchClause().Accept(v).(*ast.MatchClause)
	} else if ctx.UnwindClause() != nil {
		n = ctx.UnwindClause().Accept(v).(*ast.UnwindClause)
	}
	return n
}

func (v *ConvertVisitor) VisitMatchClause(ctx *MatchClauseContext) interface{} {
	match := &ast.MatchClause{}
	if ctx.OPTIONAL() != nil {
		match.Optional = true
	}
	if ctx.WhereClause() != nil {
		match.Where = ctx.WhereClause().Accept(v).(ast.Expr)
	}
	match.Pattern = ctx.Pattern().Accept(v).(*ast.Pattern)
	return match
}

func (v *ConvertVisitor) VisitUnwindClause(ctx *UnwindClauseContext) interface{} {
	unwind := &ast.UnwindClause{}
	unwind.Expr = ctx.Expr().Accept(v).(ast.Expr)
	unwind.Variable = ctx.Variable().Accept(v).(*ast.SymbolicNameNode)
	return unwind
}

func (v *ConvertVisitor) VisitUpdatingClause(ctx *UpdatingClauseContext) interface{} {
	switch {
	case ctx.CreateClause() != nil:
		return ctx.CreateClause().Accept(v).(*ast.CreateClause)
	case ctx.MergeClause() != nil:
		return ctx.MergeClause().Accept(v).(*ast.MergeClause)
	case ctx.SetClause() != nil:
		return ctx.SetClause().Accept(v).(*ast.SetClause)
	case ctx.DeleteClause() != nil:
		return ctx.DeleteClause().Accept(v).(*ast.DeleteClause)
	case ctx.RemoveClause() != nil:
		return ctx.RemoveClause().Accept(v).(*ast.RemoveClause)
	default:
		return nil
	}
}

func (v *ConvertVisitor) VisitCreateClause(ctx *CreateClauseContext) interface{} {
	create := &ast.CreateClause{}
	create.Pattern = ctx.Pattern().Accept(v).(*ast.Pattern)
	return create
}

func (v *ConvertVisitor) VisitSetClause(ctx *SetClauseContext) interface{} {
	set := &ast.SetClause{}
	var items []*ast.SetItem
	for _, item := range ctx.AllSetItem() {
		items = append(items, item.Accept(v).(*ast.SetItem))
	}
	set.SetItems = items
	return set
}

func (v *ConvertVisitor) VisitSetItem(ctx *SetItemContext) interface{} {
	setItem := &ast.SetItem{}

	if ctx.PropertyExpr() != nil {
		setItem.Type = ast.SetItemProperty
		setItem.Property = ctx.PropertyExpr().Accept(v).(*ast.PropertyExpr)
	} else if len(ctx.GetTokens(3)) > 0 {
		// 3 presents '=' token, see Cypher.tokens
		setItem.Type = ast.SetItemVariableAssignment
		setItem.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
		setItem.Expr = ctx.Expr().Accept(v).(ast.Expr)
	} else if len(ctx.GetTokens(4)) > 0 {
		// 4 presents '+=' token, see Cypher.tokens
		setItem.Type = ast.SetItemVariableIncrement
		setItem.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
		setItem.Expr = ctx.Expr().Accept(v).(ast.Expr)
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

func (v *ConvertVisitor) VisitDeleteClause(ctx *DeleteClauseContext) interface{} {
	deleteClause := &ast.DeleteClause{}
	if ctx.DETACH() != nil {
		deleteClause.Detach = true
	}
	var exprs []ast.Expr
	for _, expr := range ctx.AllExpr() {
		exprs = append(exprs, expr.Accept(v).(ast.Expr))
	}
	deleteClause.Exprs = exprs
	return deleteClause
}

func (v *ConvertVisitor) VisitRemoveClause(ctx *RemoveClauseContext) interface{} {
	removeClause := &ast.RemoveClause{}
	var items []*ast.RemoveItem
	for _, item := range ctx.AllRemoveItem() {
		items = append(items, item.Accept(v).(*ast.RemoveItem))
	}
	removeClause.RemoveItems = items
	return removeClause
}

func (v *ConvertVisitor) VisitRemoveItem(ctx *RemoveItemContext) interface{} {
	removeItem := &ast.RemoveItem{}
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

func (v *ConvertVisitor) VisitMergeClause(ctx *MergeClauseContext) interface{} {
	mergeClause := &ast.MergeClause{}
	mergeClause.PatternPart = ctx.PatternPart().Accept(v).(*ast.PatternPart)
	var actions []*ast.MergeAction
	for _, action := range ctx.AllMergeAction() {
		actions = append(actions, action.Accept(v).(*ast.MergeAction))
	}
	mergeClause.MergeActions = actions
	return mergeClause
}

func (v *ConvertVisitor) VisitMergeAction(ctx *MergeActionContext) interface{} {
	mergeAction := &ast.MergeAction{}
	if ctx.CREATE() != nil {
		mergeAction.Type = ast.MergeActionCreate
	} else if ctx.MATCH() != nil {
		mergeAction.Type = ast.MergeActionMatch
	}
	mergeAction.Set = ctx.SetClause().Accept(v).(*ast.SetClause)
	return mergeAction
}

func (v *ConvertVisitor) VisitWhereClause(ctx *WhereClauseContext) interface{} {
	where := ctx.Expr().Accept(v).(ast.Expr)
	return where
}

func (v *ConvertVisitor) VisitVariable(ctx *VariableContext) interface{} {
	variable := &ast.VariableNode{}
	variable.SymbolicName = ctx.SymbolicName().Accept(v).(*ast.SymbolicNameNode)
	return variable
}

func (v *ConvertVisitor) VisitSymbolicName(ctx *SymbolicNameContext) interface{} {
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

func (v *ConvertVisitor) VisitNodeLabels(ctx *NodeLabelsContext) interface{} {
	var labels []*ast.NodeLabelNode
	for _, l := range ctx.AllNodeLabel() {
		labels = append(labels, l.Accept(v).(*ast.NodeLabelNode))
	}
	return labels
}

func (v *ConvertVisitor) VisitNodeLabel(ctx *NodeLabelContext) interface{} {
	nodeLabel := &ast.NodeLabelNode{}
	nodeLabel.LabelName = ctx.LabelName().Accept(v).(*ast.SchemaNameNode)
	return nodeLabel
}

func (v *ConvertVisitor) VisitLabelName(ctx *LabelNameContext) interface{} {
	return ctx.SchemaName().Accept(v).(*ast.SchemaNameNode)
}

func (v *ConvertVisitor) VisitPattern(ctx *PatternContext) interface{} {
	pattern := &ast.Pattern{}
	var parts []*ast.PatternPart
	for _, part := range ctx.AllPatternPart() {
		parts = append(parts, part.Accept(v).(*ast.PatternPart))
	}
	pattern.Parts = parts
	return pattern
}

func (v *ConvertVisitor) VisitPatternPart(ctx *PatternPartContext) interface{} {
	patternPart := &ast.PatternPart{}
	if ctx.Variable() != nil {
		patternPart.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
	}
	patternPart.Element = ctx.
		AnonymousPatternPart().
		Accept(v).(*AnonymousPatternPartContext).
		PatternElement().
		Accept(v).(*ast.PatternElement)
	return patternPart
}

func (v *ConvertVisitor) VisitAnonymousPatternPart(ctx *AnonymousPatternPartContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitPatternElement(ctx *PatternElementContext) interface{} {
	patternElement := &ast.PatternElement{}
	// strip useless parenthesises recursively
	if ctx.PatternElement() != nil {
		ctx = ctx.PatternElement().Accept(v).(*PatternElementContext)
		return ctx
	}

	var nodes []*ast.NodePattern
	nodes = append(nodes, ctx.NodePattern().Accept(v).(*ast.NodePattern))
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

func (v *ConvertVisitor) VisitNodePattern(ctx *NodePatternContext) interface{} {
	nodePattern := &ast.NodePattern{}
	if ctx.Variable() != nil {
		nodePattern.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
	}
	if ctx.NodeLabels() != nil {
		nodePattern.Labels = ctx.NodeLabels().Accept(v).([]*ast.NodeLabelNode)
	}
	if ctx.Properties() != nil {
		nodePattern.Properties = ctx.Properties().Accept(v).(*ast.Properties)
	}
	return nodePattern
}

func (v *ConvertVisitor) VisitPatternElementChain(ctx *PatternElementChainContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitRelationshipPattern(ctx *RelationshipPatternContext) interface{} {
	relationshipPattern := &ast.RelationshipPattern{}
	if ctx.LeftArrowHead() != nil && ctx.RightArrowHead() != nil {
		relationshipPattern.Type = ast.RelationshipBoth
	} else if ctx.LeftArrowHead() != nil {
		relationshipPattern.Type = ast.RelationshipIn
	} else if ctx.RightArrowHead() != nil {
		relationshipPattern.Type = ast.RelationshipOut
	} else {
		relationshipPattern.Type = ast.RelationshipAll
	}
	if ctx.RelationshipDetail() != nil {
		relationshipPattern.Detail = ctx.RelationshipDetail().Accept(v).(*ast.RelationshipDetail)
	}
	return relationshipPattern
}

func (v *ConvertVisitor) VisitLeftArrowHead(ctx *LeftArrowHeadContext) interface{} {
	panic("Shouln't be invoked")
}

func (v *ConvertVisitor) VisitRightArrowHead(ctx *RightArrowHeadContext) interface{} {
	panic("Shouln't be invoked")
}

func (v *ConvertVisitor) VisitDash(ctx *DashContext) interface{} {
	panic("Shouln't be invoked")
}

func (v *ConvertVisitor) VisitRelationshipDetail(ctx *RelationshipDetailContext) interface{} {
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

func (v *ConvertVisitor) VisitRangeLiteral(ctx *RangeLiteralContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitMinHops(ctx *MinHopsContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitMaxHops(ctx *MaxHopsContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitRelationshipTypes(ctx *RelationshipTypesContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitRelTypeName(ctx *RelTypeNameContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitSchemaName(ctx *SchemaNameContext) interface{} {
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

func (v *ConvertVisitor) VisitReturnClause(ctx *ReturnClauseContext) interface{} {
	returnClause := &ast.ReturnClause{}
	if ctx.DISTINCT() != nil {
		returnClause.Distinct = true
	}
	returnClause.ReturnBody = ctx.ReturnBody().Accept(v).(*ast.ReturnBody)
	return returnClause
}

func (v *ConvertVisitor) VisitReturnBody(ctx *ReturnBodyContext) interface{} {
	returnBody := &ast.ReturnBody{}
	returnItems := ctx.ReturnItems().Accept(v).([]*ast.ReturnItem)
	returnBody.ReturnItems = returnItems
	if ctx.OrderClause() != nil {
		returnBody.OrderBy = ctx.OrderClause().Accept(v).(*ast.OrderClause)
	}
	if ctx.SkipClause() != nil {
		returnBody.Skip = ctx.SkipClause().(*SkipClauseContext).Expr().Accept(v).(ast.Expr)
	}
	if ctx.LimitClause() != nil {
		returnBody.Limit = ctx.LimitClause().(*LimitClauseContext).Expr().Accept(v).(ast.Expr)
	}
	return returnBody
}

func (v *ConvertVisitor) VisitReturnItems(ctx *ReturnItemsContext) interface{} {
	var returnItems []*ast.ReturnItem
	if len(ctx.GetTokens(5)) > 0 {
		returnItems = []*ast.ReturnItem{
			&ast.ReturnItem{
				Wildcard: true,
			},
		}
		return returnItems
	}
	for _, item := range ctx.AllReturnItem() {
		returnItems = append(returnItems, item.Accept(v).(*ast.ReturnItem))
	}
	return returnItems
}

func (v *ConvertVisitor) VisitReturnItem(ctx *ReturnItemContext) interface{} {
	returnItem := &ast.ReturnItem{}
	returnItem.Expr = ctx.Expr().Accept(v).(ast.Expr)
	if ctx.AS() != nil {
		returnItem.As = true
		returnItem.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
	}
	return returnItem
}

func (v *ConvertVisitor) VisitOrderClause(ctx *OrderClauseContext) interface{} {
	orderClause := &ast.OrderClause{}
	var sortItems []*ast.SortItem
	for _, item := range ctx.AllSortItem() {
		sortItems = append(sortItems, item.Accept(v).(*ast.SortItem))
	}
	return orderClause
}

func (v *ConvertVisitor) VisitSkipClause(ctx *SkipClauseContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitLimitClause(ctx *LimitClauseContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitSortItem(ctx *SortItemContext) interface{} {
	sortItem := &ast.SortItem{}
	sortItem.Expr = ctx.Expr().Accept(v).(ast.Expr)
	if ctx.ASC() != nil || ctx.ASCENDING() != nil {
		sortItem.Type = ast.SortAscending
	} else if ctx.DESC() != nil || ctx.DESCENDING() != nil {
		sortItem.Type = ast.SortDescending
	}
	return sortItem
}

func (v *ConvertVisitor) VisitExpr(ctx *ExprContext) interface{} {
	expr := ctx.OrExpr().Accept(v).(ast.Expr)
	return expr
}

func (v *ConvertVisitor) VisitOrExpr(ctx *OrExprContext) interface{} {
	if len(ctx.AllXorExpr()) == 1 {
		return ctx.XorExpr(0).Accept(v)
	}
	var exprs []ast.Expr
	for _, expr := range ctx.AllXorExpr() {
		exprs = append(exprs, expr.Accept(v).(ast.Expr))
	}
	binaryExpr := &ast.BinaryExpr{}
	binaryExpr.L = exprs[0]
	exprs = exprs[1:]
	for i, expr := range exprs {
		binaryExpr.Op = ast.OpOr
		binaryExpr.R = expr
		if i < len(exprs)-1 {
			binaryExpr = &ast.BinaryExpr{
				L: binaryExpr,
			}
		}
	}
	return binaryExpr
}

func (v *ConvertVisitor) VisitXorExpr(ctx *XorExprContext) interface{} {
	if len(ctx.AllAndExpr()) == 1 {
		return ctx.AndExpr(0).Accept(v)
	}
	var exprs []ast.Expr
	for _, expr := range ctx.AllAndExpr() {
		exprs = append(exprs, expr.Accept(v).(ast.Expr))
	}
	binaryExpr := &ast.BinaryExpr{}
	binaryExpr.L = exprs[0]
	exprs = exprs[1:]
	for i, expr := range exprs {
		binaryExpr.Op = ast.OpXor
		binaryExpr.R = expr
		if i < len(exprs)-1 {
			binaryExpr = &ast.BinaryExpr{
				L: binaryExpr,
			}
		}
	}
	return binaryExpr
}

func (v *ConvertVisitor) VisitAndExpr(ctx *AndExprContext) interface{} {
	if len(ctx.AllNotExpr()) == 1 {
		return ctx.NotExpr(0).Accept(v)
	}
	var exprs []ast.Expr
	for _, expr := range ctx.AllNotExpr() {
		exprs = append(exprs, expr.Accept(v).(ast.Expr))
	}
	binaryExpr := &ast.BinaryExpr{}
	binaryExpr.L = exprs[0]
	exprs = exprs[1:]
	for i, expr := range exprs {
		binaryExpr.Op = ast.OpAnd
		binaryExpr.R = expr
		if i < len(exprs)-1 {
			binaryExpr = &ast.BinaryExpr{
				L: binaryExpr,
			}
		}
	}
	return binaryExpr
}

func (v *ConvertVisitor) VisitNotExpr(ctx *NotExprContext) interface{} {
	if len(ctx.AllNOT()) == 0 {
		return ctx.ComparisonExpr().Accept(v)
	}
	unaryExpr := &ast.UnaryExpr{}
	unaryExpr.V = ctx.ComparisonExpr().Accept(v).(ast.Expr)
	for i := range ctx.AllNOT() {
		unaryExpr.Op = ast.OpNot
		if i < len(ctx.AllNOT())-1 {
			unaryExpr = &ast.UnaryExpr{
				V: unaryExpr,
			}
		}
	}
	return unaryExpr
}

// used for simplify parsing
type partialComparisonExpr struct {
	Type ast.OpType
	Expr ast.Expr
}

func (v *ConvertVisitor) VisitComparisonExpr(ctx *ComparisonExprContext) interface{} {
	if len(ctx.AllPartialComparisonExpr()) == 0 {
		return ctx.AddOrSubtractExpr().Accept(v)
	}
	expr := ctx.AddOrSubtractExpr().Accept(v).(ast.Expr)
	var partialExprs []partialComparisonExpr
	for _, expr := range ctx.AllPartialComparisonExpr() {
		partialExprs = append(partialExprs, expr.Accept(v).(partialComparisonExpr))
	}

	binaryExpr := &ast.BinaryExpr{}
	binaryExpr.L = expr
	for i, expr := range partialExprs {
		binaryExpr.Op = expr.Type
		binaryExpr.R = expr.Expr
		if i < len(partialExprs)-1 {
			binaryExpr = &ast.BinaryExpr{
				L: binaryExpr,
			}
		}
	}
	return binaryExpr
}

func (v *ConvertVisitor) VisitPartialComparisonExpr(ctx *PartialComparisonExprContext) interface{} {
	partialComparisonExpr := partialComparisonExpr{}
	comp := ctx.GetChild(0).(antlr.TerminalNode)
	switch comp.GetText() {
	case "=":
		partialComparisonExpr.Type = ast.OpEQ
	case "<>":
		partialComparisonExpr.Type = ast.OpNE
	case "<":
		partialComparisonExpr.Type = ast.OpLT
	case ">":
		partialComparisonExpr.Type = ast.OpGT
	case "<=":
		partialComparisonExpr.Type = ast.OpLTE
	case ">=":
		partialComparisonExpr.Type = ast.OpGTE
	}
	partialComparisonExpr.Expr = ctx.AddOrSubtractExpr().Accept(v).(ast.Expr)
	return partialComparisonExpr
}

func (v *ConvertVisitor) VisitAddOrSubtractExpr(ctx *AddOrSubtractExprContext) interface{} {
	if len(ctx.AllMultiplyDivideModuloExpr()) == 1 {
		return ctx.MultiplyDivideModuloExpr(0).Accept(v)
	}
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
	var exprs []ast.Expr
	for _, expr := range ctx.AllMultiplyDivideModuloExpr() {
		exprs = append(exprs, expr.Accept(v).(ast.Expr))
	}

	binaryExpr := &ast.BinaryExpr{}
	binaryExpr.L = exprs[0]
	exprs = exprs[1:]
	for i, expr := range exprs {
		binaryExpr.Op = ops[i]
		binaryExpr.R = expr
		if i < len(exprs)-1 {
			binaryExpr = &ast.BinaryExpr{
				L: binaryExpr,
			}
		}
	}
	return binaryExpr
}

func (v *ConvertVisitor) VisitMultiplyDivideModuloExpr(ctx *MultiplyDivideModuloExprContext) interface{} {
	if len(ctx.AllPowerOfExpr()) == 1 {
		return ctx.PowerOfExpr(0).Accept(v)
	}
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
	var exprs []ast.Expr
	for _, expr := range ctx.AllPowerOfExpr() {
		exprs = append(exprs, expr.Accept(v).(ast.Expr))
	}

	binaryExpr := &ast.BinaryExpr{}
	binaryExpr.L = exprs[0]
	exprs = exprs[1:]
	for i, expr := range exprs {
		binaryExpr.Op = ops[i]
		binaryExpr.R = expr
		if i < len(exprs)-1 {
			binaryExpr = &ast.BinaryExpr{
				L: binaryExpr,
			}
		}
	}
	return binaryExpr
}

func (v *ConvertVisitor) VisitPowerOfExpr(ctx *PowerOfExprContext) interface{} {
	if len(ctx.AllUnaryAddOrSubtractExpr()) == 1 {
		return ctx.UnaryAddOrSubtractExpr(0).Accept(v)
	}
	var exprs []ast.Expr
	for _, expr := range ctx.AllUnaryAddOrSubtractExpr() {
		exprs = append(exprs, expr.Accept(v).(ast.Expr))
	}

	binaryExpr := &ast.BinaryExpr{}
	binaryExpr.L = exprs[0]
	exprs = exprs[1:]
	for i, expr := range exprs {
		binaryExpr.Op = ast.OpPow
		binaryExpr.R = expr
		if i < len(exprs)-1 {
			binaryExpr = &ast.BinaryExpr{
				L: binaryExpr,
			}
		}
	}
	return binaryExpr
}

func (v *ConvertVisitor) VisitUnaryAddOrSubtractExpr(ctx *UnaryAddOrSubtractExprContext) interface{} {
	if len(ctx.GetTokens(13))+len(ctx.GetTokens(14)) == 0 {
		return ctx.StringListNullOperatorExpr().Accept(v)
	}
	var ops []string
	for _, child := range ctx.GetChildren() {
		if n, ok := child.GetPayload().(*antlr.CommonToken); ok && n.GetTokenType() != CypherLexerSP {
			ops = append(ops, n.GetText())
		}
	}

	unaryExpr := &ast.UnaryExpr{}
	unaryExpr.V = ctx.StringListNullOperatorExpr().Accept(v).(ast.Expr)
	for i, op := range ops {
		if op == "+" {
			unaryExpr.Op = ast.OpPlus
		} else if op == "-" {
			unaryExpr.Op = ast.OpMinus
		}
		if i < len(ops)-1 {
			unaryExpr = &ast.UnaryExpr{
				V: unaryExpr,
			}
		}
	}
	return unaryExpr
}

func (v *ConvertVisitor) VisitStringListNullOperatorExpr(ctx *StringListNullOperatorExprContext) interface{} {
	if len(ctx.AllStringOperatorExpr())+len(ctx.AllListOperatorExpr())+len(ctx.AllNullOperatorExpr()) == 0 {
		return ctx.PropertyOrLabelsExpr().Accept(v)
	}

	predicationExpr := &ast.PredicationExpr{}
	for i, child := range ctx.GetChildren() {
		if i == 0 {
			continue
		}
		switch expr := child.GetPayload().(type) {
		case StringOperatorExprContext:
			predicationExpr.Type = ast.PredicationStringOp
			predicationExpr.Expr = expr.Accept(v).(ast.Expr)
		case ListOperatorExprContext:
			predicationExpr.Type = ast.PredicationListOp
			predicationExpr.Expr = expr.Accept(v).(ast.Expr)
		case NullOperatorExprContext:
			predicationExpr.Type = ast.PredicationNullOp
			predicationExpr.Expr = expr.Accept(v).(ast.Expr)
		}

		if i < len(ctx.GetChildren())-2 {
			predicationExpr = &ast.PredicationExpr{
				Expr: predicationExpr,
			}
		}
	}

	return predicationExpr
}

func (v *ConvertVisitor) VisitListOperatorExpr(ctx *ListOperatorExprContext) interface{} {
	listOperatorExpr := &ast.ListOperationExpr{}
	if ctx.PropertyOrLabelsExpr() != nil {
		listOperatorExpr.InExpr = ctx.PropertyOrLabelsExpr().Accept(v).(*ast.PropertyOrLabelsExpr)
	} else if len(ctx.GetTokens(12)) > 0 {
		listOperatorExpr.LowerBound = ctx.Expr(0).Accept(v).(ast.Expr)
		listOperatorExpr.UpperBound = ctx.Expr(1).Accept(v).(ast.Expr)
	} else if len(ctx.GetTokens(12)) == 0 {
		expr := ctx.Expr(0).Accept(v).(ast.Expr)
		listOperatorExpr.SingleExpr = expr
	}
	return listOperatorExpr
}

func (v *ConvertVisitor) VisitStringOperatorExpr(ctx *StringOperatorExprContext) interface{} {
	stringOperatorExpr := &ast.StringOperationExpr{}
	if ctx.STARTS() != nil {
		stringOperatorExpr.Type = ast.StringOperationStartsWith
	} else if ctx.ENDS() != nil {
		stringOperatorExpr.Type = ast.StringOperationEndsWith
	} else if ctx.CONTAINS() != nil {
		stringOperatorExpr.Type = ast.StringOperationContains
	}
	stringOperatorExpr.Expr = ctx.PropertyOrLabelsExpr().Accept(v).(*ast.PropertyOrLabelsExpr)
	return stringOperatorExpr
}

func (v *ConvertVisitor) VisitNullOperatorExpr(ctx *NullOperatorExprContext) interface{} {
	nullOperatorExpr := &ast.NullOperationExpr{}
	if ctx.NOT() == nil {
		nullOperatorExpr.IsIsNull = true
	}
	return nullOperatorExpr
}

func (v *ConvertVisitor) VisitPropertyOrLabelsExpr(ctx *PropertyOrLabelsExprContext) interface{} {
	if len(ctx.AllPropertyLookup()) == 0 && ctx.NodeLabels() == nil {
		return ctx.Atom().Accept(v)
	}
	propertyOrLabelsExpr := &ast.PropertyOrLabelsExpr{}
	propertyOrLabelsExpr.Atom = ctx.Atom().Accept(v).(ast.Expr)
	if ctx.NodeLabels() != nil {
		var labels []*ast.NodeLabelNode
		for _, label := range ctx.NodeLabels().Accept(v).([]*ast.NodeLabelNode) {
			labels = append(labels, label)
		}
		propertyOrLabelsExpr.NodeLabels = labels
	}
	var lookups []*ast.PropertyLookup
	for _, lookup := range ctx.AllPropertyLookup() {
		lookups = append(lookups, lookup.Accept(v).(*ast.PropertyLookup))
	}
	propertyOrLabelsExpr.PropertyLookups = lookups
	return propertyOrLabelsExpr
}

func (v *ConvertVisitor) VisitPropertyLookup(ctx *PropertyLookupContext) interface{} {
	propertyLookup := &ast.PropertyLookup{}
	propertyLookup.PropertyKey = ctx.PropertyKeyName().Accept(v).(*PropertyKeyNameContext).SchemaName().Accept(v).(*ast.SchemaNameNode)
	return propertyLookup
}

func (v *ConvertVisitor) VisitPropertyKeyName(ctx *PropertyKeyNameContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitCaseExpr(ctx *CaseExprContext) interface{} {
	caseExpr := &ast.CaseExpr{}
	var alts []*ast.CaseAlt
	for _, alt := range ctx.AllCaseAlternatives() {
		alts = append(alts, alt.Accept(v).(*ast.CaseAlt))
	}
	caseExpr.Alts = alts
	if ctx.ELSE() != nil {
		// means there are at least 1 Expr
		if len(ctx.AllExpr()) > 1 {
			caseExpr.Expr = ctx.Expr(0).Accept(v).(ast.Expr)
			caseExpr.Else = ctx.Expr(1).Accept(v).(ast.Expr)
		} else if len(ctx.AllExpr()) == 1 {
			caseExpr.Expr = ctx.Expr(0).Accept(v).(ast.Expr)
		}
	} else if len(ctx.AllExpr()) > 0 {
		caseExpr.Expr = ctx.Expr(0).Accept(v).(ast.Expr)
	}
	return caseExpr
}

func (v *ConvertVisitor) VisitCaseAlternatives(ctx *CaseAlternativesContext) interface{} {
	caseAlt := &ast.CaseAlt{}
	caseAlt.When = ctx.Expr(0).Accept(v).(ast.Expr)
	caseAlt.Then = ctx.Expr(1).Accept(v).(ast.Expr)
	return caseAlt
}

func (v *ConvertVisitor) VisitAtom(ctx *AtomContext) interface{} {
	switch {
	case ctx.Literal() != nil:
		return ctx.Literal().Accept(v)
	case ctx.Parameter() != nil:
		return ctx.Parameter().Accept(v)
	case ctx.CaseExpr() != nil:
		return ctx.CaseExpr().Accept(v)
	case ctx.COUNT() != nil:
		return &ast.CountAllExpr{}
	case ctx.ListComprehension() != nil:
		return ctx.ListComprehension().Accept(v)
	case ctx.PatternComprehension() != nil:
		return ctx.PatternComprehension().Accept(v)
	case ctx.ALL() != nil:
		filter := ctx.FilterExpr().Accept(v).(*ast.FilterExpr)
		filter.Type = ast.FilterAll
		return filter
	case ctx.ANY() != nil:
		filter := ctx.FilterExpr().Accept(v).(*ast.FilterExpr)
		filter.Type = ast.FilterAny
		return filter
	case ctx.NONE() != nil:
		filter := ctx.FilterExpr().Accept(v).(*ast.FilterExpr)
		filter.Type = ast.FilterNone
		return filter
	case ctx.SINGLE() != nil:
		filter := ctx.FilterExpr().Accept(v).(*ast.FilterExpr)
		filter.Type = ast.FilterSingle
		return filter
	case ctx.RelationshipsPattern() != nil:
		return ctx.RelationshipsPattern().Accept(v)
	case ctx.ParenthesizedExpr() != nil:
		return ctx.ParenthesizedExpr().Accept(v)
	case ctx.FunctionInvocation() != nil:
		// TODO
		panic("FunctionInvocation not support now")
	case ctx.Variable() != nil:
		return ctx.Variable().Accept(v)
	default:
		panic("")
	}
}

func (v *ConvertVisitor) VisitRelationshipsPattern(ctx *RelationshipsPatternContext) interface{} {
	patternElement := &ast.PatternElement{}
	// same as VisitPatternElement
	var nodes []*ast.NodePattern
	nodes = append(nodes, ctx.NodePattern().Accept(v).(*ast.NodePattern))
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

func (v *ConvertVisitor) VisitParenthesizedExpr(ctx *ParenthesizedExprContext) interface{} {
	parenExpr := &ast.ParenExpr{}
	parenExpr.Expr = ctx.Expr().Accept(v).(ast.Expr)
	return parenExpr
}

func (v *ConvertVisitor) VisitFilterExpr(ctx *FilterExprContext) interface{} {
	filterExpr := &ast.FilterExpr{}
	idInColl := ctx.IdInColl().Accept(v).(*IdInCollContext)
	filterExpr.Variable = idInColl.Variable().Accept(v).(*ast.VariableNode)
	filterExpr.In = idInColl.Expr().Accept(v).(ast.Expr)
	if ctx.WhereClause() != nil {
		filterExpr.Where = ctx.WhereClause().Accept(v).(ast.Expr)
	}
	return filterExpr
}

func (v *ConvertVisitor) VisitIdInColl(ctx *IdInCollContext) interface{} {
	return ctx
}

func (v *ConvertVisitor) VisitListComprehension(ctx *ListComprehensionContext) interface{} {
	listComprehension := &ast.ListComprehension{}
	listComprehension.FilterExpr = ctx.FilterExpr().Accept(v).(*ast.FilterExpr)
	if ctx.Expr() != nil {
		listComprehension.Expr = ctx.Expr().Accept(v).(ast.Expr)
	}
	return listComprehension
}

func (v *ConvertVisitor) VisitPatternComprehension(ctx *PatternComprehensionContext) interface{} {
	patternComprehension := &ast.PatternComprehension{}
	if ctx.Variable() != nil {
		patternComprehension.Variable = ctx.Variable().Accept(v).(*ast.VariableNode)
	}
	patternComprehension.PatternElement = ctx.RelationshipsPattern().Accept(v).(*ast.PatternElement)
	if ctx.WHERE() != nil {
		patternComprehension.Where = ctx.Expr(0).Accept(v).(ast.Expr)
	}
	patternComprehension.Expr = ctx.AllExpr()[len(ctx.AllExpr())-1].Accept(v).(ast.Expr)
	return patternComprehension
}

func unquote(s string) string {
	if s[0] == s[len(s)-1] &&
		(s[0] == '\'' || s[0] == '"' || s[0] == '`') {
		return s[1 : len(s)-1]
	} else {
		return s
	}
}

func (v *ConvertVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	literal := &ast.LiteralExpr{}
	if ctx.NumberLiteral() != nil {
		literal.Type = ast.LiteralNumber
		literal.Number = ctx.NumberLiteral().Accept(v).(*ast.NumberLiteral)
	} else if ctx.StringLiteral() != nil {
		literal.Type = ast.LiteralString
		literal.String = unquote(ctx.StringLiteral().GetText())
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

func (v *ConvertVisitor) VisitNumberLiteral(ctx *NumberLiteralContext) interface{} {
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

func (v *ConvertVisitor) VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{} {
	panic("Shouln't be invoked")
}

func (v *ConvertVisitor) VisitMapLiteral(ctx *MapLiteralContext) interface{} {
	mapLiteral := &ast.MapLiteral{}
	var keys []*ast.SchemaNameNode
	var exprs []ast.Expr
	for i := range ctx.AllPropertyKeyName() {
		keys = append(keys, ctx.PropertyKeyName(i).Accept(v).(*PropertyKeyNameContext).SchemaName().Accept(v).(*ast.SchemaNameNode))
		exprs = append(exprs, ctx.Expr(i).Accept(v).(ast.Expr))
	}
	mapLiteral.PropertyKeys = keys
	mapLiteral.Exprs = exprs
	return mapLiteral
}

func (v *ConvertVisitor) VisitListLiteral(ctx *ListLiteralContext) interface{} {
	listLiteral := &ast.ListLiteral{}
	var exprs []ast.Expr
	for _, expr := range ctx.AllExpr() {
		exprs = append(exprs, expr.Accept(v).(ast.Expr))
	}
	listLiteral.Exprs = exprs
	return listLiteral
}

func (v *ConvertVisitor) VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{} {
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

func (v *ConvertVisitor) VisitDoubleLiteral(ctx *DoubleLiteralContext) interface{} {
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

func (v *ConvertVisitor) VisitParameter(ctx *ParameterContext) interface{} {
	parameter := &ast.ParameterNode{}
	if ctx.SymbolicName() != nil {
		parameter.Type = ast.ParameterSymbolicName
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

func (v *ConvertVisitor) VisitProperties(ctx *PropertiesContext) interface{} {
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

func (v *ConvertVisitor) VisitPropertyExpr(ctx *PropertyExprContext) interface{} {
	propertyExpr := &ast.PropertyExpr{}
	propertyExpr.Atom = ctx.Atom().Accept(v).(ast.Expr)
	var lookups []*ast.PropertyLookup
	for _, lookup := range ctx.AllPropertyLookup() {
		lookups = append(lookups, lookup.Accept(v).(*ast.PropertyLookup))
	}
	propertyExpr.Lookups = lookups
	return propertyExpr
}

func (v *ConvertVisitor) VisitReservedWord(ctx *ReservedWordContext) interface{} {
	reservedWord := &ast.ReservedWordNode{}
	reservedWord.Content = ctx.GetText()
	return reservedWord
}

func (v *ConvertVisitor) VisitStandaloneCall(ctx *StandaloneCallContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *ConvertVisitor) VisitYieldItems(ctx *YieldItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ConvertVisitor) VisitYieldItem(ctx *YieldItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *ConvertVisitor) VisitFunctionInvocation(ctx *FunctionInvocationContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *ConvertVisitor) VisitFunctionName(ctx *FunctionNameContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *ConvertVisitor) VisitExplicitProcedureInvocation(ctx *ExplicitProcedureInvocationContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *ConvertVisitor) VisitImplicitProcedureInvocation(ctx *ImplicitProcedureInvocationContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *ConvertVisitor) VisitProcedureResultField(ctx *ProcedureResultFieldContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *ConvertVisitor) VisitProcedureName(ctx *ProcedureNameContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *ConvertVisitor) VisitNamespace(ctx *NamespaceContext) interface{} {
	// TODO
	panic("not implemented")
}

func (v *ConvertVisitor) VisitInQueryCall(ctx *InQueryCallContext) interface{} {
	// TODO
	panic("not implemented")
}
