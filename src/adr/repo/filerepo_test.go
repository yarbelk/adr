package repo

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/yarbelk/adr/src/adr"
)

var goodTomlADR = `Title = "Test Title"
Number = 1
Authors = ["Test User"]
Created = 2019-01-01T00:00:00Z
Status = "DRAFT"
Impact = "High"
Text = "Good Text"
`

func TestGet(t *testing.T) {
	created, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
	expected := &adr.ADR{
		Title:   "Test Title",
		Number:  1,
		Authors: []string{"Test User"},
		Created: created,
		Status:  adr.Draft,
		Impact:  adr.High,
		Text:    "Good Text",
	}
	testCase := []struct {
		testName      string
		input         int
		expected      *adr.ADR
		expectedError error
	}{
		{
			testName:      "GET AN ADR",
			input:         1,
			expected:      expected,
			expectedError: nil,
		},
		{
			testName:      "GET NOT FOUND ADR",
			input:         2,
			expected:      nil,
			expectedError: fmt.Errorf("Cannot find ADR-0002"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.testName, func(t *testing.T) {
			repo := FileRepo{
				adrs: map[int]adr.ADR{1: *expected},
			}
			actual, actualErr := repo.Get(tc.input)
			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedError, actualErr)
		})
	}
}
