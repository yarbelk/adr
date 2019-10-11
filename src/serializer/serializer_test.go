package serializer_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/yarbelk/adr/src/adr"
	"gitlab.com/yarbelk/adr/src/serializer"
)

var goodTomlADR = `Title = "Test Title"
Number = 1
Authors = ["Test User"]
Created = 2019-01-01T00:00:00Z
Status = "DRAFT"
Impact = "High"
Text = "Good Text"
`

func TestTOMLUnmarshal(t *testing.T) {
	var tomlUnmarshal serializer.Unmarshaller
	reader := bytes.NewBufferString(goodTomlADR)
	tomlUnmarshal = serializer.NewUnmarshal(reader)
	created, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
	expected := adr.ADR{
		Title:   "Test Title",
		Number:  1,
		Authors: []string{"Test User"},
		Created: created,
		Status:  adr.Draft,
		Impact:  adr.High,
		Text:    "Good Text",
	}
	newADR := adr.ADR{}
	if err := tomlUnmarshal.Unmarshal(&newADR); err != nil {
		t.Logf("%+v expected true and nil", err)
		t.FailNow()
	}
	assert.Equal(t, expected, newADR)
}

func TestTOMLMarshal(t *testing.T) {
	created, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
	testCase := []struct {
		testName      string
		input         adr.ADR
		expected      string
		expectedError error
	}{
		{
			testName: "GOT SUCCESS MARSHALLING ADR",
			input: adr.ADR{
				Title:   "Test Title",
				Number:  1,
				Authors: []string{"Test User"},
				Created: created,
				Status:  adr.Draft,
				Impact:  adr.High,
				Text:    "Good Text",
			},
			expected:      goodTomlADR,
			expectedError: nil,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.testName, func(t *testing.T) {
			var writer *bytes.Buffer = new(bytes.Buffer)
			tomlMarshal := serializer.NewMarshal(writer)
			err := tomlMarshal.Marshal(&tc.input)
			actual := writer.String()
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
