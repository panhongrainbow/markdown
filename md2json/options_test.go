package md2json

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Check_md2JsonOptions(t *testing.T) {

	empty := "empty"

	paraOpts := ParagraphOptions{
		empty: empty,
	}

	prefix := "Table No"

	tbOpts := TableOptions{
		PrefixTbName: prefix,
	}

	visitor := NewJsonVisitor(WithTableOptions(tbOpts), WithParagraphOptions(paraOpts))

	require.Equal(t, prefix, visitor.Table.PrefixTbName)
	require.Equal(t, empty, visitor.Paragraph.empty)
}
