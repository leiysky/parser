package parser

import "testing"

func TestParser(t *testing.T) {
	parser := New()
	parser.Parse(`
MATCH (n)
WHERE n.id = 1
RETURN n`)
}
