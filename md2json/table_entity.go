package md2json

import (
	"github.com/panhongrainbow/goCodePebblez/bytez"
	"github.com/panhongrainbow/markdown/ast"
	"github.com/panhongrainbow/markdown/parser"
	"github.com/panhongrainbow/markdown/syncPool"
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

// Visit visits various nodes in the AST and responding based on the node type and entry status,
// the objective is to collect data by converting markdown to JSON.
func (visitor *JsonVisitor) Visit(node ast.Node, entering bool) (status ast.WalkStatus) {
	switch n := node.(type) {
	case *ast.Document:
		visitor.location.inDocument = entering
	case *ast.Paragraph:
		// Check if the paragraph content has a prefix matching the table name
		// and update the visitor's table name accordingly.
		if bytez.HasPrefix(n.Content, bytez.StringToReadOnlyBytes(visitor.Table.PrefixTbName)) {
			if visitor.Table.WipePrefix {
				visitor.tableName = string(n.Content[len(visitor.Table.PrefixTbName)+1:])
			} else {
				visitor.tableName = string(n.Content)
			}
		}
	case *ast.Table:
		if entering {
			// Start constructing the JSON document for the table.
			visitor.JsonDoc = append(visitor.JsonDoc, "{")
			visitor.JsonDoc = append(visitor.JsonDoc, quot+"type"+quot+colon+quot+"table"+quot+comma+
				quot+"name"+quot+colon+quot+visitor.tableName+quot+comma+
				quot+"data"+quot+colon)
			// visitor.Header = visitor.Header[:0]
			// syncPool.GlobalStringSlice.Put(&visitor.Header)
		}
		visitor.location.inTable = entering
		lastIndex := len(visitor.JsonDoc) - 1
		if visitor.JsonDoc[lastIndex] == "," {
			visitor.JsonDoc = visitor.JsonDoc[:lastIndex]
		}
		if !entering {
			// Finalize the JSON document for the table and save it.
			visitor.JsonDoc = append(visitor.JsonDoc, "}")
			visitor.JsonDocs = append(visitor.JsonDocs, strings.Join(visitor.JsonDoc, ""))
			// visitor.JsonDoc = visitor.JsonDoc[:0]

			syncPool.GlobalStringSlice.Put(&visitor.JsonDoc)
			visitor.JsonDoc = syncPool.GlobalStringSlice.Get()

			syncPool.GlobalStringSlice.Put(&visitor.Header)
			visitor.Header = syncPool.GlobalStringSlice.Get()
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
			// Start constructing the JSON array for the table body.
			visitor.JsonDoc = append(visitor.JsonDoc, "[")
		} else {
			// Finalize the JSON array for the table body.
			visitor.closeJSONObject("]")
		}
	case *ast.TableRow:
		visitor.location.inTableRow = entering
		visitor.columnIndex = 0
		if visitor.location.inTableBody {
			if entering {
				// Start constructing the JSON object for the table row.
				visitor.JsonDoc = append(visitor.JsonDoc, "{")
			} else {
				// Finalize the JSON object for the table row.
				visitor.closeJSONObject("}")
			}
		}
	case *ast.TableCell:
		// Handle table cell node.
		if visitor.location.inTableHeader {
			// Collect the content of table cells in the header.
			visitor.Header = append(visitor.Header, string(n.Content))
		}

		if visitor.location.inTableBody {
			if !bytez.IsAllBytesDigits(n.Content) {
				// Add quotation marks to non-numeric content.
				n.Content = bytez.Quotation(n.Content)
			}

			tableValue := string(n.Content)
			if len(n.Content) == 0 {
				// Handle empty cells with replacement.
				// Because JSON values cannot be left blank, they are replaced according to the configured values.
				// (JSON 值不能留空，根据配置进行替换)
				tableValue = quot + visitor.Table.ReplaceEmpty + quot
			}

			// Construct the JSON key-value pair for the cell content.
			visitor.JsonDoc = append(visitor.JsonDoc, "\""+visitor.Header[visitor.columnIndex]+"\":", tableValue, ",")

			visitor.columnIndex++
		}
		status = ast.Inquired
	case *ast.TableFooter:
		// Do nothing for table footer nodes.
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

func MdToJson(md []byte, optFuncs ...SetOptsFunc) (JsonDocs []string) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	p.Block(md)

	v := NewJsonVisitor(optFuncs...)

	ast.Walk(p.Doc, v)

	return v.JsonDocs
}
