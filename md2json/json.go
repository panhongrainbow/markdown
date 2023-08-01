package md2json

import (
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"strings"
	"unicode"
)

type JsonVisitor struct {
	JsonDocs    []string
	JsonDoc     []string
	Header      []string
	tableName   string
	columnIndex int
	location    struct {
		inDocument, inTable, inTableHeader, inTableBody, inTableRow bool
	}
	visitorEmbeddedOpts
}

const (
	quot  string = "\""
	colon string = ":"
	comma string = ","
)

func (visitor *JsonVisitor) Visit(node ast.Node, entering bool) (status ast.WalkStatus) {
	switch n := node.(type) {
	case *ast.Document:
		visitor.location.inDocument = entering
	case *ast.Paragraph:
		visitor.tableName = string(n.Content)
	case *ast.Table:
		if entering {
			visitor.JsonDoc = append(visitor.JsonDoc, "{")
			// visitor.JsonDoc = append(visitor.JsonDoc, "\""+visitor.tableName+"\":")
			visitor.JsonDoc = append(visitor.JsonDoc, quot+"type"+quot+colon+quot+"table"+quot+comma+
				quot+"name"+quot+colon+quot+visitor.tableName+quot+comma+
				quot+"data"+quot+colon)
			visitor.Header = visitor.Header[:0]
		}
		visitor.location.inTable = entering
		lastIndex := len(visitor.JsonDoc) - 1
		if visitor.JsonDoc[lastIndex] == "," {
			visitor.JsonDoc = visitor.JsonDoc[:lastIndex]
		}
		if !entering {
			// save
			visitor.JsonDoc = append(visitor.JsonDoc, "}")
			visitor.JsonDocs = append(visitor.JsonDocs, strings.Join(visitor.JsonDoc, ""))
			visitor.JsonDoc = visitor.JsonDoc[:0]
		}
	case *ast.TableHeader:
		visitor.location.inTableHeader = entering
	case *ast.TableBody:
		visitor.location.inTableBody = entering
		if entering && len(visitor.Header) == 0 {
			panic("Before entering the table body, the table header has not been created yet.")
			return
		}
		if entering {
			visitor.JsonDoc = append(visitor.JsonDoc, "[")
		} else {
			visitor.closeJSONObject("]")
		}
	case *ast.TableRow:
		visitor.location.inTableRow = entering
		visitor.columnIndex = 0
		if visitor.location.inTableBody {
			if entering {
				visitor.JsonDoc = append(visitor.JsonDoc, "{")
			} else {
				visitor.closeJSONObject("}")
			}
		}
	case *ast.TableCell:
		if visitor.location.inTableHeader {
			visitor.Header = append(visitor.Header, string(n.Content))
		}

		if visitor.location.inTableBody {
			if !IsAllDigits(n.Content) {
				n.Content = Quotation(n.Content)
			}

			visitor.JsonDoc = append(visitor.JsonDoc, "\""+visitor.Header[visitor.columnIndex]+"\":", string(n.Content), ",")

			visitor.columnIndex++
		}
		status = ast.Inquired
	case *ast.TableFooter:
		// do not thing!
	}
	return
}

//go:inline
func (visitor *JsonVisitor) closeJSONObject(closeTag string) {
	lastIndex := len(visitor.JsonDoc) - 1
	if visitor.JsonDoc[lastIndex] == "," {
		visitor.JsonDoc[lastIndex] = closeTag
		visitor.JsonDoc = append(visitor.JsonDoc, ",")
	} else {
		visitor.JsonDoc = append(visitor.JsonDoc, closeTag, ",")
	}
}

func IsAllDigits(data []byte) bool {
	for _, b := range data {
		if !unicode.IsDigit(rune(b)) {
			return false
		}
	}
	return true
}

func Quotation(slice []byte) (quoted []byte) {
	quoted = make([]byte, len(slice)+2, len(slice)+2)
	quoted[0] = 34
	quoted[len(quoted)-1] = 34
	copy(quoted[1:], slice)
	return
}

func mdToJson(md []byte, optFuncs ...SetOptsFunc) (JsonDocs []string) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	p.Block(md)

	v := NewJsonVisitor(optFuncs...)

	ast.Walk(p.Doc, v)

	return v.JsonDocs
}

func HasPrefix(a, b []byte) bool {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
