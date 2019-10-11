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
	"io/ioutil"
	"log"
	"os"
	"path"
	goTemplate "text/template"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/yarbelk/adr/src/adr"
)

// Render out all the ADRs
func Render(renderDir, baseTmpl, fileDir string) {
	files, err := ioutil.ReadDir(fileDir)
	tmpl := goTemplate.Must(goTemplate.New("adr").Parse(baseTmpl))

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		inputFile := file.Name()
		log.Println("looking at filename", inputFile)
		adr := adr.ADR{}
		func() {
			f, err := os.OpenFile(path.Join(fileDir, inputFile), os.O_RDONLY, 0644)
			defer f.Close()
			if err != nil {
				log.Fatal("failed to open adr file", inputFile, err)
			}
			toml.DecodeReader(f, &adr)
		}()
		outputFile := path.Join(renderDir, file.Name()+".md")

		f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal("Cant open file", err)
		}
		if err := tmpl.Execute(f, adr); err != nil {
			log.Fatal(err)
		}
	}
}

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "render all the ADRs to human readable format like markdown",
	Long: `render all (or one) ADR to markdown or some other format
	it will default to github flavoured markdown, and a simple layout.
	This can be overridden by setting the global config to having different
	base template for the rendering.`,
	Run: func(cmd *cobra.Command, args []string) {
		renderDir := viper.GetString("renderDir")
		baseTmpl := global.GetString("baseTemplate")
		fileDir := viper.GetString("ADRDir")
		Render(renderDir, baseTmpl, fileDir)
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
