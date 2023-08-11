package md2json

import (
	"encoding/json"
	"github.com/panhongrainbow/goCodePebblez/bytez"
	"github.com/panhongrainbow/markdown/syncPool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// mysticIsle serves as the test data,
// containing information about the residents in various regions and towns on the island.
var mysticIsle = `
# MysticIsle Residents Information

Welcome to MysticIsle, a place with diverse towns and vibrant communities.
Below is the detailed information about the residents of different towns on the island.

## Marigold Town

Config Table Marigold Town
Name   | Age | Phone    | Address    | Job
-------|-----|----------|------------|--------------
Charlie| 31  | 555-1234 | 123 Main St| Engineer
Eva    | 27  | 555-4321 | 456 Elm St | Teacher

(Robert's information is left blank temporarily as he is in the process of changing his phone number.)

## Oakwood City

Config Table Oakwood City
Name   | Age | Gender | Phone    | Address
-------|-----|--------|----------|--------------
Robert | 28  | male   |          | 789 Oak St
Alyssa | 29  | female | 555-8765 | 101 Pine St
`

// Test_Check_Json2Table validates table to json conversion with assertions.
func Test_Check_Json2Table(t *testing.T) {
	// Resident struct to hold resident information.
	type resident struct {
		Name    string `json:"Name"`
		Age     int    `json:"Age"`
		Gender  string `json:"Gender"`
		Phone   string `json:"Phone"`
		Address string `json:"Address"`
		Job     string `json:"Job"`
	}

	// Region struct to hold town information.
	type region struct {
		Type string     `json:"type"`
		Name string     `json:"name"`
		Data []resident `json:"data"`
	}

	// Define table options for conversion.
	tbOpts := TableOptions{
		PrefixTbName: "Config Table",
		ReplaceEmpty: "-",
		WipePrefix:   true,
	}

	// Convert Markdown to JSON using specified options.
	jsonDocs := MdToJson(
		bytez.StringToReadOnlyBytes(mysticIsle),
		WithTableOptions(tbOpts),
	)

	// Create a collection to hold region information.
	collections := make([]region, 0, 2)

	// Unmarshal JSON and populate the collections.
	for i := 0; i < len(jsonDocs); i++ {
		var jsonRow region
		err := json.Unmarshal([]byte(jsonDocs[i]), &jsonRow)
		require.NoError(t, err, "The JSON string (%s) failed to unmarshal", jsonDocs[i])

		collections = append(collections, jsonRow)
	}

	// Assertions to validate parsed JSON data for the first region.
	assert.Equal(t, "table", collections[0].Type)
	assert.Equal(t, "Marigold Town", collections[0].Name)

	// Assertions to validate parsed JSON data for the first resident.
	assert.Equal(t, "Charlie", collections[0].Data[0].Name)
	assert.Equal(t, 31, collections[0].Data[0].Age)
	assert.Equal(t, "", collections[0].Data[0].Gender)
	assert.Equal(t, "555-1234", collections[0].Data[0].Phone)
	assert.Equal(t, "123 Main St", collections[0].Data[0].Address)
	assert.Equal(t, "Engineer", collections[0].Data[0].Job)

	// Assertions to validate parsed JSON data for the second resident.
	assert.Equal(t, "Eva", collections[0].Data[1].Name)
	assert.Equal(t, 27, collections[0].Data[1].Age)
	assert.Equal(t, "", collections[0].Data[1].Gender)
	assert.Equal(t, "555-4321", collections[0].Data[1].Phone)
	assert.Equal(t, "456 Elm St", collections[0].Data[1].Address)
	assert.Equal(t, "Teacher", collections[0].Data[1].Job)

	// Assertions to validate parsed JSON data for the second region.
	assert.Equal(t, "table", collections[1].Type)
	assert.Equal(t, "Oakwood City", collections[1].Name)

	// Assertions to validate parsed JSON data for the third resident.
	assert.Equal(t, "Robert", collections[1].Data[0].Name)
	assert.Equal(t, 28, collections[1].Data[0].Age)
	assert.Equal(t, "male", collections[1].Data[0].Gender)
	assert.Equal(t, "-", collections[1].Data[0].Phone)
	assert.Equal(t, "789 Oak St", collections[1].Data[0].Address)
	assert.Equal(t, "", collections[1].Data[0].Job)

	assert.Equal(t, "Alyssa", collections[1].Data[1].Name)
	assert.Equal(t, 29, collections[1].Data[1].Age)
	assert.Equal(t, "female", collections[1].Data[1].Gender)
	assert.Equal(t, "555-8765", collections[1].Data[1].Phone)
	assert.Equal(t, "101 Pine St", collections[1].Data[1].Address)
	assert.Equal(t, "", collections[1].Data[1].Job)
}

func Benchmark_Check_Json2Table(b *testing.B) {
	triggerInit := syncPool.GlobalStringSlice.Get()
	syncPool.GlobalStringSlice.Put(&triggerInit)
	// Define table options for conversion.
	tbOpts := TableOptions{
		PrefixTbName: "Config Table",
		ReplaceEmpty: "-",
		WipePrefix:   true,
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Convert Markdown to JSON using specified options.
		_ = MdToJson(
			bytez.StringToReadOnlyBytes(mysticIsle),
			WithTableOptions(tbOpts),
		)
	}
}
