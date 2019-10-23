package serializer_test

import (
	"bytes"
	"fmt"
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
	created, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
	name := "%s %s"
	baseToml := `Title = "Test Title"
Number = 1
Authors = ["Test User"]
Created = 2019-01-01T00:00:00Z
Status = "%s"
Impact = "%s"
Text = "Good Text"
`
	statuses := []adr.Status{adr.Draft, adr.Approved, adr.Superceded}
	impacts := []adr.Impact{adr.Unknown, adr.Low, adr.Medium, adr.High}

	var tomlUnmarshal serializer.Unmarshaller
	for _, impact := range impacts {
		for _, status := range statuses {
			t.Run(fmt.Sprintf(name, impact, status), func(t *testing.T) {
				expected := adr.ADR{
					Title:   "Test Title",
					Number:  1,
					Authors: []string{"Test User"},
					Created: created,
					Status:  status,
					Impact:  impact,
					Text:    "Good Text",
				}
				reader := bytes.NewBufferString(fmt.Sprintf(baseToml, status, impact))
				tomlUnmarshal = serializer.NewUnmarshal(reader)
				newADR := adr.ADR{}
				if err := tomlUnmarshal.Unmarshal(&newADR); err != nil {
					t.Logf("%+v expected true and nil", err)
					t.FailNow()
				}
				assert.Equal(t, expected, newADR)
			})
		}
	}
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
