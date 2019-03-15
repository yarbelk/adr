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
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yarbelk/adr/src/adr"
)

type authorsFlag []string

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

var authors authorsFlag

const template = `# Background

Enter the background to his ADR here.  Setup the context
so we are on the same page.  Keep it simple and easy to
follow. Don't tell me the problem

# Complication

Now tell me where the problem/complication is that this ADR
is addressing.

# Options Considered

1. This was one
2. This was another

# Decision

What did we decided

# Outcome
`

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
		a := adr.ADR{
			Title:   args[0],
			Number:  next,
			Authors: authors,
			Created: time.Now(),
			Status:  adr.Draft,
			Text:    text,
		}
		fp := filepath.Join(filedir, a.Filename())
		f, err = os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0755)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
			fmt.Print(text)
			os.Exit(1)
		}
		f.Seek(0, 0)
		if err := toml.NewEncoder(f).Encode(a); err != nil {
			log.Fatal(err)
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
	newCmd.Flags().VarP(&authors, "authors", "", "Authors, repeat this flag for multiple")
}
