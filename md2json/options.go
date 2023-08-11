package md2json

type SetOptsFunc func(*JsonVisitor)

func NewJsonVisitor(funcs ...SetOptsFunc) (visitor *JsonVisitor) {
	visitor = new(JsonVisitor)

	for _, eachFunc := range funcs {
		eachFunc(visitor)
	}

	return
}

// visitorEmbeddedOpts is a collection of all the options and will be embedded into visitor
type visitorEmbeddedOpts struct {
	Table     TableOptions
	Paragraph ParagraphOptions
}

// Options for Paragraph node

type ParagraphOptions struct {
	empty string
}

func WithParagraphOptions(paragraph ParagraphOptions) SetOptsFunc {
	return func(mdOpts *JsonVisitor) {
		mdOpts.Paragraph = paragraph
	}
}

// Options for Table node

type TableOptions struct {
	PrefixTbName string
	ReplaceEmpty string
	WipePrefix   bool
}

func WithTableOptions(table TableOptions) SetOptsFunc {
	return func(mdOpts *JsonVisitor) {
		mdOpts.Table = table
	}
}
