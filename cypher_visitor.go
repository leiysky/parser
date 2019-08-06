// Generated from Cypher.g4 by ANTLR 4.7.

package parser // Cypher

import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by CypherParser.
type CypherVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by CypherParser#cypher.
	VisitCypher(ctx *CypherContext) interface{}

	// Visit a parse tree produced by CypherParser#stmt.
	VisitStmt(ctx *StmtContext) interface{}

	// Visit a parse tree produced by CypherParser#query.
	VisitQuery(ctx *QueryContext) interface{}

	// Visit a parse tree produced by CypherParser#regularQuery.
	VisitRegularQuery(ctx *RegularQueryContext) interface{}

	// Visit a parse tree produced by CypherParser#unionClause.
	VisitUnionClause(ctx *UnionClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#singleQuery.
	VisitSingleQuery(ctx *SingleQueryContext) interface{}

	// Visit a parse tree produced by CypherParser#singlePartQuery.
	VisitSinglePartQuery(ctx *SinglePartQueryContext) interface{}

	// Visit a parse tree produced by CypherParser#multiPartQuery.
	VisitMultiPartQuery(ctx *MultiPartQueryContext) interface{}

	// Visit a parse tree produced by CypherParser#updatingClause.
	VisitUpdatingClause(ctx *UpdatingClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#readingClause.
	VisitReadingClause(ctx *ReadingClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#matchClause.
	VisitMatchClause(ctx *MatchClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#unwindClause.
	VisitUnwindClause(ctx *UnwindClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#mergeClause.
	VisitMergeClause(ctx *MergeClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#mergeAction.
	VisitMergeAction(ctx *MergeActionContext) interface{}

	// Visit a parse tree produced by CypherParser#createClause.
	VisitCreateClause(ctx *CreateClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#setClause.
	VisitSetClause(ctx *SetClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#setItem.
	VisitSetItem(ctx *SetItemContext) interface{}

	// Visit a parse tree produced by CypherParser#deleteClause.
	VisitDeleteClause(ctx *DeleteClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#removeClause.
	VisitRemoveClause(ctx *RemoveClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#removeItem.
	VisitRemoveItem(ctx *RemoveItemContext) interface{}

	// Visit a parse tree produced by CypherParser#inQueryCall.
	VisitInQueryCall(ctx *InQueryCallContext) interface{}

	// Visit a parse tree produced by CypherParser#standaloneCall.
	VisitStandaloneCall(ctx *StandaloneCallContext) interface{}

	// Visit a parse tree produced by CypherParser#yieldItems.
	VisitYieldItems(ctx *YieldItemsContext) interface{}

	// Visit a parse tree produced by CypherParser#yieldItem.
	VisitYieldItem(ctx *YieldItemContext) interface{}

	// Visit a parse tree produced by CypherParser#withClause.
	VisitWithClause(ctx *WithClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#returnClause.
	VisitReturnClause(ctx *ReturnClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#returnBody.
	VisitReturnBody(ctx *ReturnBodyContext) interface{}

	// Visit a parse tree produced by CypherParser#returnItems.
	VisitReturnItems(ctx *ReturnItemsContext) interface{}

	// Visit a parse tree produced by CypherParser#returnItem.
	VisitReturnItem(ctx *ReturnItemContext) interface{}

	// Visit a parse tree produced by CypherParser#orderClause.
	VisitOrderClause(ctx *OrderClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#skipClause.
	VisitSkipClause(ctx *SkipClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#limitClause.
	VisitLimitClause(ctx *LimitClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#sortItem.
	VisitSortItem(ctx *SortItemContext) interface{}

	// Visit a parse tree produced by CypherParser#whereClause.
	VisitWhereClause(ctx *WhereClauseContext) interface{}

	// Visit a parse tree produced by CypherParser#pattern.
	VisitPattern(ctx *PatternContext) interface{}

	// Visit a parse tree produced by CypherParser#patternPart.
	VisitPatternPart(ctx *PatternPartContext) interface{}

	// Visit a parse tree produced by CypherParser#anonymousPatternPart.
	VisitAnonymousPatternPart(ctx *AnonymousPatternPartContext) interface{}

	// Visit a parse tree produced by CypherParser#patternElement.
	VisitPatternElement(ctx *PatternElementContext) interface{}

	// Visit a parse tree produced by CypherParser#nodePattern.
	VisitNodePattern(ctx *NodePatternContext) interface{}

	// Visit a parse tree produced by CypherParser#patternElementChain.
	VisitPatternElementChain(ctx *PatternElementChainContext) interface{}

	// Visit a parse tree produced by CypherParser#relationshipPattern.
	VisitRelationshipPattern(ctx *RelationshipPatternContext) interface{}

	// Visit a parse tree produced by CypherParser#relationshipDetail.
	VisitRelationshipDetail(ctx *RelationshipDetailContext) interface{}

	// Visit a parse tree produced by CypherParser#properties.
	VisitProperties(ctx *PropertiesContext) interface{}

	// Visit a parse tree produced by CypherParser#relationshipTypes.
	VisitRelationshipTypes(ctx *RelationshipTypesContext) interface{}

	// Visit a parse tree produced by CypherParser#nodeLabels.
	VisitNodeLabels(ctx *NodeLabelsContext) interface{}

	// Visit a parse tree produced by CypherParser#nodeLabel.
	VisitNodeLabel(ctx *NodeLabelContext) interface{}

	// Visit a parse tree produced by CypherParser#rangeLiteral.
	VisitRangeLiteral(ctx *RangeLiteralContext) interface{}

	// Visit a parse tree produced by CypherParser#minHops.
	VisitMinHops(ctx *MinHopsContext) interface{}

	// Visit a parse tree produced by CypherParser#maxHops.
	VisitMaxHops(ctx *MaxHopsContext) interface{}

	// Visit a parse tree produced by CypherParser#labelName.
	VisitLabelName(ctx *LabelNameContext) interface{}

	// Visit a parse tree produced by CypherParser#relTypeName.
	VisitRelTypeName(ctx *RelTypeNameContext) interface{}

	// Visit a parse tree produced by CypherParser#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by CypherParser#orExpr.
	VisitOrExpr(ctx *OrExprContext) interface{}

	// Visit a parse tree produced by CypherParser#xorExpr.
	VisitXorExpr(ctx *XorExprContext) interface{}

	// Visit a parse tree produced by CypherParser#andExpr.
	VisitAndExpr(ctx *AndExprContext) interface{}

	// Visit a parse tree produced by CypherParser#notExpr.
	VisitNotExpr(ctx *NotExprContext) interface{}

	// Visit a parse tree produced by CypherParser#comparisonExpr.
	VisitComparisonExpr(ctx *ComparisonExprContext) interface{}

	// Visit a parse tree produced by CypherParser#addOrSubtractExpr.
	VisitAddOrSubtractExpr(ctx *AddOrSubtractExprContext) interface{}

	// Visit a parse tree produced by CypherParser#multiplyDivideModuloExpr.
	VisitMultiplyDivideModuloExpr(ctx *MultiplyDivideModuloExprContext) interface{}

	// Visit a parse tree produced by CypherParser#powerOfExpr.
	VisitPowerOfExpr(ctx *PowerOfExprContext) interface{}

	// Visit a parse tree produced by CypherParser#unaryAddOrSubtractExpr.
	VisitUnaryAddOrSubtractExpr(ctx *UnaryAddOrSubtractExprContext) interface{}

	// Visit a parse tree produced by CypherParser#stringListNullOperatorExpr.
	VisitStringListNullOperatorExpr(ctx *StringListNullOperatorExprContext) interface{}

	// Visit a parse tree produced by CypherParser#listOperatorExpr.
	VisitListOperatorExpr(ctx *ListOperatorExprContext) interface{}

	// Visit a parse tree produced by CypherParser#stringOperatorExpr.
	VisitStringOperatorExpr(ctx *StringOperatorExprContext) interface{}

	// Visit a parse tree produced by CypherParser#nullOperatorExpr.
	VisitNullOperatorExpr(ctx *NullOperatorExprContext) interface{}

	// Visit a parse tree produced by CypherParser#propertyOrLabelsExpr.
	VisitPropertyOrLabelsExpr(ctx *PropertyOrLabelsExprContext) interface{}

	// Visit a parse tree produced by CypherParser#atom.
	VisitAtom(ctx *AtomContext) interface{}

	// Visit a parse tree produced by CypherParser#literal.
	VisitLiteral(ctx *LiteralContext) interface{}

	// Visit a parse tree produced by CypherParser#booleanLiteral.
	VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{}

	// Visit a parse tree produced by CypherParser#listLiteral.
	VisitListLiteral(ctx *ListLiteralContext) interface{}

	// Visit a parse tree produced by CypherParser#partialComparisonExpr.
	VisitPartialComparisonExpr(ctx *PartialComparisonExprContext) interface{}

	// Visit a parse tree produced by CypherParser#parenthesizedExpr.
	VisitParenthesizedExpr(ctx *ParenthesizedExprContext) interface{}

	// Visit a parse tree produced by CypherParser#relationshipsPattern.
	VisitRelationshipsPattern(ctx *RelationshipsPatternContext) interface{}

	// Visit a parse tree produced by CypherParser#filterExpr.
	VisitFilterExpr(ctx *FilterExprContext) interface{}

	// Visit a parse tree produced by CypherParser#idInColl.
	VisitIdInColl(ctx *IdInCollContext) interface{}

	// Visit a parse tree produced by CypherParser#functionInvocation.
	VisitFunctionInvocation(ctx *FunctionInvocationContext) interface{}

	// Visit a parse tree produced by CypherParser#functionName.
	VisitFunctionName(ctx *FunctionNameContext) interface{}

	// Visit a parse tree produced by CypherParser#explicitProcedureInvocation.
	VisitExplicitProcedureInvocation(ctx *ExplicitProcedureInvocationContext) interface{}

	// Visit a parse tree produced by CypherParser#implicitProcedureInvocation.
	VisitImplicitProcedureInvocation(ctx *ImplicitProcedureInvocationContext) interface{}

	// Visit a parse tree produced by CypherParser#procedureResultField.
	VisitProcedureResultField(ctx *ProcedureResultFieldContext) interface{}

	// Visit a parse tree produced by CypherParser#procedureName.
	VisitProcedureName(ctx *ProcedureNameContext) interface{}

	// Visit a parse tree produced by CypherParser#namespace.
	VisitNamespace(ctx *NamespaceContext) interface{}

	// Visit a parse tree produced by CypherParser#listComprehension.
	VisitListComprehension(ctx *ListComprehensionContext) interface{}

	// Visit a parse tree produced by CypherParser#patternComprehension.
	VisitPatternComprehension(ctx *PatternComprehensionContext) interface{}

	// Visit a parse tree produced by CypherParser#propertyLookup.
	VisitPropertyLookup(ctx *PropertyLookupContext) interface{}

	// Visit a parse tree produced by CypherParser#caseExpr.
	VisitCaseExpr(ctx *CaseExprContext) interface{}

	// Visit a parse tree produced by CypherParser#caseAlternatives.
	VisitCaseAlternatives(ctx *CaseAlternativesContext) interface{}

	// Visit a parse tree produced by CypherParser#variable.
	VisitVariable(ctx *VariableContext) interface{}

	// Visit a parse tree produced by CypherParser#numberLiteral.
	VisitNumberLiteral(ctx *NumberLiteralContext) interface{}

	// Visit a parse tree produced by CypherParser#mapLiteral.
	VisitMapLiteral(ctx *MapLiteralContext) interface{}

	// Visit a parse tree produced by CypherParser#parameter.
	VisitParameter(ctx *ParameterContext) interface{}

	// Visit a parse tree produced by CypherParser#propertyExpr.
	VisitPropertyExpr(ctx *PropertyExprContext) interface{}

	// Visit a parse tree produced by CypherParser#propertyKeyName.
	VisitPropertyKeyName(ctx *PropertyKeyNameContext) interface{}

	// Visit a parse tree produced by CypherParser#integerLiteral.
	VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{}

	// Visit a parse tree produced by CypherParser#doubleLiteral.
	VisitDoubleLiteral(ctx *DoubleLiteralContext) interface{}

	// Visit a parse tree produced by CypherParser#schemaName.
	VisitSchemaName(ctx *SchemaNameContext) interface{}

	// Visit a parse tree produced by CypherParser#reservedWord.
	VisitReservedWord(ctx *ReservedWordContext) interface{}

	// Visit a parse tree produced by CypherParser#symbolicName.
	VisitSymbolicName(ctx *SymbolicNameContext) interface{}

	// Visit a parse tree produced by CypherParser#leftArrowHead.
	VisitLeftArrowHead(ctx *LeftArrowHeadContext) interface{}

	// Visit a parse tree produced by CypherParser#rightArrowHead.
	VisitRightArrowHead(ctx *RightArrowHeadContext) interface{}

	// Visit a parse tree produced by CypherParser#dash.
	VisitDash(ctx *DashContext) interface{}
}
