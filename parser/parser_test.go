package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_Parse_SingleNumberExpression(t *testing.T) {
	src := "1 + 2 * 3 - 4"
	parser := New(src)
	rootNode := parser.Parse()
	assert.NotNil(t, rootNode)
}

