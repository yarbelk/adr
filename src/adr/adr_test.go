package adr_test

import (
	"testing"

	"gitlab.com/yarbelk/adr/src/adr"
)

func TestFilename(t *testing.T) {
	a := adr.ADR{Number: 1}
	if a.Filename() != "ADR-0001" {
		t.Logf("wanted 'ADR-0001', got '%s'", a.Filename())
		t.Fail()
	}
}

func TestSupercede(t *testing.T) {
	original := &adr.ADR{Number: 1}
	replacement := &adr.ADR{Number: 2}
	replacement.Supercede(original)
	if original.SupercededBy[0] != 2 {
		t.Logf("Expected SupercededBy '%d', got '%v'", 2, original.SupercededBy)
		t.Fail()
	}
	if replacement.Supercedes[0] != 1 {
		t.Logf("Expected Supercedes '%d', got '%v'", 1, replacement.Supercedes)
		t.Fail()
	}
}

func TestSupercededAgain(t *testing.T) {
	original := &adr.ADR{Number: 1, SupercededBy: []int{2}}
	replacement := &adr.ADR{Number: 3}
	replacement.Supercede(original)
	expected := []int{2, 3}
	for i, v := range expected {
		if original.SupercededBy[i] != v {
			t.Logf("Expected SupercededBy '%v', got '%v'", expected, original.SupercededBy)
			t.Fail()
		}
	}
}

func TestSupercededDup(t *testing.T) {
	original := &adr.ADR{Number: 2, SupercededBy: []int{3}}
	replacement := &adr.ADR{Number: 3}
	replacement.Supercede(original)
	expected := []int{3}
	for i, v := range expected {
		if original.SupercededBy[i] != v {
			t.Logf("Expected SupercededBy '%v', got '%v'", expected, original.SupercededBy)
			t.Fail()
		}
	}
}

func TestSupercedeAgain(t *testing.T) {
	original := &adr.ADR{Number: 2}
	replacement := &adr.ADR{Number: 3, Supercedes: []int{1}}
	replacement.Supercede(original)
	expected := []int{1, 2}
	for i, v := range expected {
		if replacement.Supercedes[i] != v {
			t.Logf("Expected Supercedes '%v', got '%v'", expected, replacement.Supercedes)
			t.Fail()
		}
	}
}

func TestSupercedeDuplicate(t *testing.T) {
	original := &adr.ADR{Number: 3}
	replacement := &adr.ADR{Number: 4, Supercedes: []int{3}}
	replacement.Supercede(original)
	expected := []int{3}
	for i, v := range expected {
		if replacement.Supercedes[i] != v {
			t.Logf("Expected Supercedes '%v', got '%v'", expected, replacement.Supercedes)
			t.Fail()
		}
	}
}
