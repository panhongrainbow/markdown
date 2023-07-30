package json

import (
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"unicode"
)

type JsonVisitor struct {
	JsonDoc     []string
	Header      []string
	columnIndex int
	location    struct {
		inDocument    bool
		inTable       bool
		inTableHeader bool
		inTableBody   bool
		inTableRow    bool
	}
}

func (visitor *JsonVisitor) Visit(node ast.Node, entering bool) (status ast.WalkStatus) {

	switch n := node.(type) {
	case *ast.Document:
		visitor.location.inDocument = entering
	case *ast.Table:
		visitor.location.inTable = entering
	case *ast.TableHeader:
		visitor.location.inTableHeader = entering
	case *ast.TableBody:
		visitor.location.inTableBody = entering
		if entering == true &&
			len(visitor.Header) == 0 {
			panic("Before entering the table body, the table header has not been created yet.")
			return
		}
		switch entering {
		case true:
			visitor.JsonDoc = append(visitor.JsonDoc, "[")
		case false:
			if visitor.JsonDoc[len(visitor.JsonDoc)-1] == "," {
				visitor.JsonDoc[len(visitor.JsonDoc)-1] = "]"
			} else {
				visitor.JsonDoc = append(visitor.JsonDoc, "]")
			}
		}
	case *ast.TableRow:
		visitor.location.inTableRow = entering
		visitor.columnIndex = 0
		if visitor.location.inTableBody == true {
			switch entering {
			case true:
				visitor.JsonDoc = append(visitor.JsonDoc, "{")
			case false:
				if visitor.JsonDoc[len(visitor.JsonDoc)-1] == "," {
					visitor.JsonDoc[len(visitor.JsonDoc)-1] = "}"
					visitor.JsonDoc = append(visitor.JsonDoc, ",")
				} else {
					visitor.JsonDoc = append(visitor.JsonDoc, "}", ",")
				}

				// visitor.JsonDoc = append(visitor.JsonDoc, "}", ",")
			}
		}
	case *ast.TableCell:
		if visitor.location.inTableHeader == true {
			visitor.Header = append(visitor.Header, string(n.Content))
		}

		if visitor.location.inTableBody == true {

			value := string(n.Content)
			if !isAllDigits(n.Content) {
				value = "\"" + value + "\""
			}

			visitor.JsonDoc = append(visitor.JsonDoc, "\""+visitor.Header[visitor.columnIndex]+"\":", value, ",")

			visitor.columnIndex++
		}
		status = ast.Inquired
	case *ast.TableFooter:
		// do not thing !
	}

	return
}

func isAllDigits(data []byte) bool {
	for _, b := range data {
		if !unicode.IsDigit(rune(b)) {
			return false
		}
	}
	return true
}

func mdToJson(md []byte) (jsonStr []string) {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	p.Block(md)

	v := JsonVisitor{}

	ast.Walk(p.Doc, &v)

	return v.JsonDoc
}
