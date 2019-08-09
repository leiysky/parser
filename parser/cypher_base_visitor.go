// Generated from Cypher.g4 by ANTLR 4.7.

package parser // Cypher

import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseCypherVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseCypherVisitor) VisitCypher(ctx *CypherContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitStmt(ctx *StmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitQuery(ctx *QueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitRegularQuery(ctx *RegularQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitUnionClause(ctx *UnionClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitSingleQuery(ctx *SingleQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitSinglePartQuery(ctx *SinglePartQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitMultiPartQuery(ctx *MultiPartQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitMultiPartQueryPartial(ctx *MultiPartQueryPartialContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitUpdatingClause(ctx *UpdatingClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitReadingClause(ctx *ReadingClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitMatchClause(ctx *MatchClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitUnwindClause(ctx *UnwindClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitMergeClause(ctx *MergeClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitMergeAction(ctx *MergeActionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitCreateClause(ctx *CreateClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitSetClause(ctx *SetClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitSetItem(ctx *SetItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitDeleteClause(ctx *DeleteClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitRemoveClause(ctx *RemoveClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitRemoveItem(ctx *RemoveItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitInQueryCall(ctx *InQueryCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitStandaloneCall(ctx *StandaloneCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitYieldItems(ctx *YieldItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitYieldItem(ctx *YieldItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitWithClause(ctx *WithClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitReturnClause(ctx *ReturnClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitReturnBody(ctx *ReturnBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitReturnItems(ctx *ReturnItemsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitReturnItem(ctx *ReturnItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitOrderClause(ctx *OrderClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitSkipClause(ctx *SkipClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitLimitClause(ctx *LimitClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitSortItem(ctx *SortItemContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitWhereClause(ctx *WhereClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitPattern(ctx *PatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitPatternPart(ctx *PatternPartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitAnonymousPatternPart(ctx *AnonymousPatternPartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitPatternElement(ctx *PatternElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitNodePattern(ctx *NodePatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitPatternElementChain(ctx *PatternElementChainContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitRelationshipPattern(ctx *RelationshipPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitRelationshipDetail(ctx *RelationshipDetailContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitProperties(ctx *PropertiesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitRelationshipTypes(ctx *RelationshipTypesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitNodeLabels(ctx *NodeLabelsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitNodeLabel(ctx *NodeLabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitRangeLiteral(ctx *RangeLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitMinHops(ctx *MinHopsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitMaxHops(ctx *MaxHopsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitLabelName(ctx *LabelNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitRelTypeName(ctx *RelTypeNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitExpr(ctx *ExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitOrExpr(ctx *OrExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitXorExpr(ctx *XorExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitAndExpr(ctx *AndExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitNotExpr(ctx *NotExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitComparisonExpr(ctx *ComparisonExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitAddOrSubtractExpr(ctx *AddOrSubtractExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitMultiplyDivideModuloExpr(ctx *MultiplyDivideModuloExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitPowerOfExpr(ctx *PowerOfExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitUnaryAddOrSubtractExpr(ctx *UnaryAddOrSubtractExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitStringListNullOperatorExpr(ctx *StringListNullOperatorExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitListOperatorExpr(ctx *ListOperatorExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitStringOperatorExpr(ctx *StringOperatorExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitNullOperatorExpr(ctx *NullOperatorExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitPropertyOrLabelsExpr(ctx *PropertyOrLabelsExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitAtom(ctx *AtomContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitListLiteral(ctx *ListLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitPartialComparisonExpr(ctx *PartialComparisonExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitParenthesizedExpr(ctx *ParenthesizedExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitRelationshipsPattern(ctx *RelationshipsPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitFilterExpr(ctx *FilterExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitIdInColl(ctx *IdInCollContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitFunctionInvocation(ctx *FunctionInvocationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitFunctionName(ctx *FunctionNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitExplicitProcedureInvocation(ctx *ExplicitProcedureInvocationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitImplicitProcedureInvocation(ctx *ImplicitProcedureInvocationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitProcedureResultField(ctx *ProcedureResultFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitProcedureName(ctx *ProcedureNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitNamespace(ctx *NamespaceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitListComprehension(ctx *ListComprehensionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitPatternComprehension(ctx *PatternComprehensionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitPropertyLookup(ctx *PropertyLookupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitCaseExpr(ctx *CaseExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitCaseAlternatives(ctx *CaseAlternativesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitVariable(ctx *VariableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitNumberLiteral(ctx *NumberLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitMapLiteral(ctx *MapLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitParameter(ctx *ParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitPropertyExpr(ctx *PropertyExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitPropertyKeyName(ctx *PropertyKeyNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitDoubleLiteral(ctx *DoubleLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitSchemaName(ctx *SchemaNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitReservedWord(ctx *ReservedWordContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitSymbolicName(ctx *SymbolicNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitLeftArrowHead(ctx *LeftArrowHeadContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitRightArrowHead(ctx *RightArrowHeadContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCypherVisitor) VisitDash(ctx *DashContext) interface{} {
	return v.VisitChildren(ctx)
}
