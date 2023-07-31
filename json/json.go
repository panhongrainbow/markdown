package json

import (
	"fmt"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"unicode"
)

type JsonVisitor struct {
	JsonDoc     []string
	Header      []string
	columnIndex int
	location    struct {
		inDocument, inTable, inTableHeader, inTableBody, inTableRow bool
	}
}

func NewJsonVisitor() *JsonVisitor {
	return &JsonVisitor{}
}

func (visitor *JsonVisitor) Visit(node ast.Node, entering bool) (status ast.WalkStatus) {
	switch n := node.(type) {
	case *ast.Document:
		visitor.location.inDocument = entering
		if entering {
			visitor.JsonDoc = append(visitor.JsonDoc, "{")
		} else {
			visitor.JsonDoc = append(visitor.JsonDoc, "}")
		}
	case *ast.Table:
		visitor.location.inTable = entering
		lastIndex := len(visitor.JsonDoc) - 1
		if visitor.JsonDoc[lastIndex] == "," {
			visitor.JsonDoc = visitor.JsonDoc[:lastIndex]
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
			if !isAllDigits(n.Content) {
				fmt.Println(222, string(n.Content))
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

func isAllDigits(data []byte) bool {
	for _, b := range data {
		if !unicode.IsDigit(rune(b)) {
			return false
		}
	}
	return true
}

func Quotation(slice []byte) []byte {
	newSlice := make([]byte, len(slice)+1)
	newSlice[0] = 34
	copy(newSlice[1:], slice)
	newSlice = append(newSlice, 34)
	return newSlice
}

func mdToJson(md []byte) (jsonStr []string) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	p.Block(md)

	v := JsonVisitor{}

	ast.Walk(p.Doc, &v)

	return v.JsonDoc
}
