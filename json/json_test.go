package json

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

var mds = `Name  | Age | Phone
------|-----|---------
Bob   | 31  | 555-1234
Alice | 27  | 555-4321
`

func Test_README(t *testing.T) {
	type Person struct {
		Name  string `json:"Name"`
		Age   int    `json:"Age"`
		Phone string `json:"Phone"`
	}

	md := []byte(mds)
	jsonStr := mdToJson(md)

	str := strings.Join(jsonStr, "")

	var p []Person

	err := json.Unmarshal([]byte(str)[1:len(str)-1], &p)

	fmt.Println(str)

	fmt.Println(err)

	fmt.Println(p)
}
