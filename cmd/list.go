/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/yarbelk/adr/src/adr"
	"gitlab.com/yarbelk/adr/src/ioutils"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		listADRs(cmd, args)
	},
}

// adrF is an adr and a filename. because i'm lazy
type adrF struct {
	filename string
	adr      adr.ADR
}

func listADRs(cmd *cobra.Command, args []string) {
	fileDir := viper.GetString("ADRDir")
	files, err := ioutils.ReadDir(fileDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	adrs := make([]adrF, 0, len(files))
	for _, file := range files {
		inputFile := file.Name()
		adr := adr.ADR{}
		func() {
			f, err := os.OpenFile(path.Join(fileDir, inputFile), os.O_RDONLY, 0644)
			defer f.Close()
			if err != nil {
				log.Fatal("failed to open adr file", inputFile, err)
			}
			toml.DecodeReader(f, &adr)
		}()
		adrs = append(adrs, adrF{inputFile, adr})
	}
	sort.Slice(adrs, func(i, j int) bool { return adrs[i].adr.Number < adrs[j].adr.Number })
	for _, a := range adrs {
		fmt.Printf("%s\t%s\t%s\n", a.filename, a.adr.Title, a.adr.Status)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
