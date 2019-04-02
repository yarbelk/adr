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
	"github.com/yarbelk/adr/src/serializer"
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
var render bool

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "New ADR",
	Long: `Create a new ADR, with a number greater than the prior
	
You need to pass in the title as a single argument:
	adr new "This is the Title"`,
	Args: cobra.ExactArgs(1),
	Run:  newADR,
}

func newADR(cmd *cobra.Command, args []string) {
	// This is setting up the 'vim buffer/file for the basic new adr'
	filedir := viper.GetString("ADRDir")

	next := adr.NextNumber(filedir)
	text, ok := getText()
	if !ok {
		log.Fatal("You need to actually edit the ADR text")
	}
	title := args[0]

	newADR := &adr.ADR{
		Title:      title,
		Number:     next,
		Authors:    authors,
		Created:    time.Now(),
		Status:     adr.Draft,
		Related:    related,
		Supercedes: supercedes,
		Text:       text,
	}
	updateAllAdrs(newADR, filedir, next)
}

func updateAllAdrs(newADR *adr.ADR, filedir string, next int) {
	// TODO: I don't like the stuff below here
	for _, r := range related {
		if r >= next {
			log.Fatalf("Cannot be related to %d, is a future ADR", r)
		}
		fp := filepath.Join(filedir, adr.ADR{Number: r}.Filename())
		if verbose {
			fmt.Println(fp)
		}

		f, err := os.OpenFile(fp, os.O_RDWR, 0644)
		defer f.Close()
		if err != nil {
			log.Fatalf("Impossible to open the adr %q, %v", fp, err)
		}
		oldADR := new(adr.ADR)
		if err := serializer.NewUnmarshal(f).Unmarshal(oldADR); err != nil {
			log.Fatal("failed to load old ADR in update related", err)
		}
		oldADR.RelatesTo(next)
		newADR.RelatesTo(oldADR.Number)
		f.Seek(0, 0)
		f.Truncate(0)
		if err := serializer.NewMarshal(f).Marshal(oldADR); err != nil {
			log.Fatal(err)
		}
	}
	for _, s := range supercedes {
		if s >= next {
			log.Fatalf("Cannot be related to %d, is a future ADR", s)
		}
		fp := filepath.Join(filedir, adr.ADR{Number: s}.Filename())
		f, err := os.OpenFile(fp, os.O_RDWR, 0644)
		defer f.Close()
		if err != nil {
			log.Fatalf("Impossible to open the adr %q, %v", fp, err)
		}
		oldADR := new(adr.ADR)
		if err := serializer.NewUnmarshal(f).Unmarshal(oldADR); err != nil {
			log.Fatal("failed to load old ADR in update related", err)
		}
		newADR.Supercede(oldADR)

		f.Seek(0, 0)
		f.Truncate(0)
		if err := serializer.NewMarshal(f).Marshal(oldADR); err != nil {
			log.Fatal(err)
		}
		if err != nil {
			log.Fatal("Cant write out the old adr in a supercede operation", err)
		}
	}
	fp := filepath.Join(filedir, newADR.Filename())
	f, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("failed to open file for new adr", err)
	}

	if err := serializer.NewMarshal(f).Marshal(newADR); err != nil {
		log.Fatal("Final Update", err)
	}
	// Note: fix this silly interface for this.
	if render {
		renderDir := viper.GetString("renderDir")
		baseTmpl := global.GetString("baseTemplate")
		fileDir := viper.GetString("ADRDir")
		Render(renderDir, baseTmpl, fileDir)
	}
}

// getText will explode on problems, like non-zero status
// codes and prevent anyting from happening.
func getText() (string, bool) {
	tmpl := global.GetString("template")
	f, err := ioutil.TempFile("", "")
	f.Write([]byte(tmpl))
	filename := f.Name()
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(filename)
	f.Close()
	editor := os.Getenv("EDITOR")
	if visual := os.Getenv("VISUAL"); visual != "" {
		editor = visual
	}
	log.Println("opening", editor, filename)
	cmd := exec.Command(editor, filename)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
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
	text := string(b)
	if text == tmpl {
		log.Println("No changes to text, ignoring new command")
		return "", false
	}
	return text, true
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
	newCmd.Flags().BoolVarP(&render, "render", "R", false, "also render the adrs")
}
