parser:
	java -jar antlr/antlr-4.7-complete.jar -Dlanguage=Go -visitor -no-listener -o . Cypher.g4