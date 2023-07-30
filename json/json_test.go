package json

import (
	"fmt"
	"testing"
)

var mds = `Name  | Age | Phone
------|-----|---------
Bob   | 31  | 555-1234
Alice | 27  | 555-4321
`

func Test_README(t *testing.T) {
	md := []byte(mds)
	jsonStr := mdToJson(md)
	fmt.Println(jsonStr)
}
