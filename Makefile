.PHONY: parser

parser/cypher_parser.go parser/cypher_lexer.go parser/cypher_base_visitor.go parser/cypher_visitor.go: Cypher.g4
	make parser

parser: fmt
	java -jar antlr/antlr-4.7-complete.jar -Dlanguage=Go -visitor -no-listener -o parser Cypher.g4

fmt:
	@echo "gofmt (simplify)"
	@ gofmt -s -l -w . 2>&1 | awk '{print} END{if(NR>0) {exit 1}}'