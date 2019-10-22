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
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a new ",
	Long: `Create the adr docs dir (if it doesn't exist) and set up a config file:
`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("init called")
			fmt.Println("dir is", viper.GetString("ADRDir"))
		}
		updateOrCreateConf(viper.GetString("ADRDir"), viper.GetString("renderDir"))
		os.MkdirAll(viper.GetString("ADRDir"), os.ModePerm)
	},
}

func updateOrCreateConf(dir, renderDir string) {
	c := make(map[string]string)
	func() {
		f, err := os.OpenFile(".adr.toml", os.O_RDWR|os.O_CREATE, 0644)
		defer f.Close()
		if err != nil {
			fmt.Println("Failed to init the conf file: ", err)
			os.Exit(1)
		}
		m, e := toml.DecodeReader(f, &c)
		if verbose {
			fmt.Printf("%v\n, %v\n %v\n", c, m, e)
		}
	}()
	f, err := os.OpenFile(".adr.toml", os.O_RDWR|os.O_TRUNC, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println("Failed to init the conf file: ", err)
		os.Exit(1)
	}
	c["ADRDir"] = dir
	c["renderDir"] = dir
	f.Seek(0, 0)
	if err := toml.NewEncoder(f).Encode(c); err != nil {
		log.Fatal(err)
	}

}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
