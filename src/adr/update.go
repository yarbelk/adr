package adr

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func (a ADR) UpdateFile(filedir string) error {
	fp := filepath.Join(filedir, a.Filename())
	fmt.Printf("Update file '%s', with ADR '%+v\n", fp, a)
	f, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	f.Seek(0, 0)
	return toml.NewEncoder(f).Encode(a)
}

// Load an ADR from a filepath, or error out
func Load(filepath string) (*ADR, error) {
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	a := new(ADR)
	_, err = toml.DecodeReader(f, a)
	return a, err
}

// RelatesTo Setup the relates to, without dups
func (a *ADR) RelatesTo(n int) {
	if a.Related == nil {
		a.Related = []int{n}
	} else {
		for _, v := range a.Related {
			if v == n {
				return
			}
		}
		a.Related = append(a.Related, n)
	}
}

// Supercedes a prior ADR, will update both
func (a *ADR) Supercede(original *ADR) {
	if a.Supercedes == nil {
		a.Supercedes = []int{original.Number}
	} else {
		for _, v := range a.Supercedes {
			if v == original.Number {
				return
			}
		}
		a.Supercedes = append(a.Supercedes, original.Number)
	}
	if original.SupercededBy == nil {
		original.SupercededBy = []int{a.Number}
	} else {
		for _, v := range original.SupercededBy {
			if v == a.Number {
				return
			}
		}
		original.SupercededBy = append(original.SupercededBy, a.Number)
	}
}
