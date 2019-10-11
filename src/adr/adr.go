package adr

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"time"
)

type Status string
type Impact string

const (
	Draft      Status = "DRAFT"
	Approved   Status = "Accepted"
	Superceded Status = "Superceded"
)

const (
	Unknown Impact = "Unknown"
	Low     Impact = "Low"
	Medium  Impact = "Medium"
	High    Impact = "High"
)

type ADR struct {
	Title        string
	Number       int
	Authors      []string
	Created      time.Time
	Status       Status
	Impact       Impact
	Related      []int
	SupercededBy []int
	Supercedes   []int
	Text         string
}

// Filename is the filename for this adr, in the format ADR-1234.  Its the 'asperational'
// or cannonical filename.
func (a ADR) Filename() string {
	return fmt.Sprintf("ADR-%04d", a.Number)
}

// NextNumber doesn't trust the count of files
func NextNumber(dir string) int {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	last := 0
	re := regexp.MustCompile("[[:digit:]]{4}")
	for _, file := range files {
		d := re.FindString(file.Name())
		if len(d) != 4 {
			continue
		}
		i, _ := strconv.Atoi(d)
		if i > last {
			last = i
		}
	}
	return last + 1
}
