package md2json

import (
	"encoding/json"
	"fmt"
	"testing"
)

var mds = `contacts
Name  | Age | Phone
------|-----|---------
Bob   | 31  | 555-1234
Alice | 27  | 555-4321
`

func Test_README(t *testing.T) {
	var mds = `contacts
Name  | Age | Phone
------|-----|---------
Bob   | 31  | 555-1234
Alice | 27  | 555-4321

contacts2
Name  | Age | Phone
------|-----|---------
Bob1   | 31  | 555-1234
Alice1 | 27  | 555-4321
`

	type Contact struct {
		Name  string `json:"Name"`
		Age   int    `json:"Age"`
		Phone string `json:"Phone"`
	}

	type Table struct {
		Type string    `json:"type"`
		Name string    `json:"name"`
		Data []Contact `json:"data"`
	}

	prefix := "Table No"

	tbOpts := TableOptions{
		PrefixTbName: prefix,
	}

	md := []byte(mds)
	jsonDocs := mdToJson(md, WithTableOptions(tbOpts))

	var p Table

	for i := 0; i < len(jsonDocs); i++ {
		err := json.Unmarshal([]byte(jsonDocs[i]), &p)
		if err != nil {
			t.Errorf("The JSON string (%s) failed to unmarshal \n", jsonDocs[i])
		}
		fmt.Println(p)
	}
}
