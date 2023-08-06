package md2json

import (
	"github.com/panhongrainbow/goCodePebblez/bytez"
	"github.com/panhongrainbow/markdown/ast"
	"github.com/panhongrainbow/markdown/parser"
	"strings"
)

const (
	quot  string = "\""
	colon string = ":"
	comma string = ","
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

func (visitor *JsonVisitor) Visit(node ast.Node, entering bool) (status ast.WalkStatus) {
	switch n := node.(type) {
	case *ast.Document:
		visitor.location.inDocument = entering
	case *ast.Paragraph:
		if bytez.HasPrefix(n.Content, bytez.StringToReadOnlyBytes(visitor.Table.PrefixTbName)) {
			visitor.tableName = string(n.Content)
		}
	case *ast.Table:
		if entering {
			visitor.JsonDoc = append(visitor.JsonDoc, "{")
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
			if !bytez.IsAllBytesDigits(n.Content) {
				n.Content = bytez.Quotation(n.Content)
			}

			tableValue := string(n.Content)
			if len(n.Content) == 0 {
				tableValue = quot + visitor.Table.ReplaceEmpty + quot
			}

			visitor.JsonDoc = append(visitor.JsonDoc, "\""+visitor.Header[visitor.columnIndex]+"\":", tableValue, ",")

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

func mdToJson(md []byte, optFuncs ...SetOptsFunc) (JsonDocs []string) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	p.Block(md)

	v := NewJsonVisitor(optFuncs...)

	ast.Walk(p.Doc, v)

	return v.JsonDocs
}
