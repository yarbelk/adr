// Copyright © 2019 James Rivett-Carnac
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yarbelk/adr/src/adr"
)

type authorsFlag []string
type adrFlag []int

func (i *authorsFlag) String() string {
	s := []string(*i)
	return fmt.Sprintf("%+v", s)
}

func (i *authorsFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *authorsFlag) Type() string {
	return "list"
}

func (i *adrFlag) String() string {
	s := []int(*i)
	return fmt.Sprintf("%+v", s)
}

func (i *adrFlag) Set(value string) error {
	v, err := strconv.Atoi(value)
	fmt.Println("related/supercedes", v, err)
	if err != nil {
		return err
	}
	*i = append(*i, v)
	return nil
}

func (i *adrFlag) Type() string {
	return "list"
}

var authors authorsFlag
var related adrFlag
var supercedes adrFlag

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "New ADR",
	Long: `Create a new ADR, with a number greater than the prior
	
You need to pass in the title as a single argument:
	adr new "This is the Title"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tmpl := global.GetString("template")
		filedir := viper.GetString("ADRDir")

		next := adr.NextNumber(filedir)
		f, err := ioutil.TempFile("", "")
		f.Write([]byte(tmpl))
		tempFile := f.Name()
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tempFile)
		f.Close()
		text := getText(tempFile)
		if text == tmpl {
			log.Println("No changes to text, ignoring new command")
			return
		}
		a := &adr.ADR{
			Title:      args[0],
			Number:     next,
			Authors:    authors,
			Created:    time.Now(),
			Status:     adr.Draft,
			Related:    related,
			Supercedes: supercedes,
			Text:       text,
		}
		for _, r := range related {
			if r >= next {
				log.Fatalf("Cannot be related to %d, is a future ADR", r)
			}
			fp := filepath.Join(filedir, adr.ADR{Number: r}.Filename())
			fmt.Println(fp)
			old, err := adr.Load(fp)
			if err != nil {
				log.Fatal("failed to load old ADR in update related", err)
			}
			old.RelatesTo(next)
			a.RelatesTo(old.Number)
			old.UpdateFile(filedir)
			if err != nil {
				log.Fatal(err)
			}
		}
		for _, s := range supercedes {
			if s >= next {
				log.Fatalf("Cannot be related to %d, is a future ADR", s)
			}
			fp := filepath.Join(filedir, adr.ADR{Number: s}.Filename())
			old, err := adr.Load(fp)
			a.Supercede(old)
			old.UpdateFile(filedir)
			if err != nil {
				log.Fatal("other here", err)
			}
		}
		if err := a.UpdateFile(filedir); err != nil {
			log.Fatal("Final Update", err)
		}
	},
}

// getText will explode on problems, like non-zero status
// codes and prevent anyting from happening.
func getText(filename string) string {
	editor := os.Getenv("EDITOR")
	if visual := os.Getenv("VISUAL"); visual != "" {
		editor = visual
	}
	log.Println("opening", editor, filename)
	cmd := exec.Command(editor, filename)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal("cmd error", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	newCmd.Flags().VarP(&authors, "authors", "a", "Authors, repeat this flag for multiple")
	newCmd.Flags().VarP(&related, "related", "r", "Any related stories")
	newCmd.Flags().VarP(&supercedes, "supercedes", "S", "Any related stories")
}
