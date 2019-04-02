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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var global *viper.Viper
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "adr",
	Version: `0.1
Copyright (C) 2019 James Rivett-Carnac
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.`,
	Short: "Manage ADRs for a project",
	Long: `Write, manage and list your  ADRs for a project.  For example:

initialize a new ADR directory:
	adr init docs/adrs

Add a new ADR:
	adr new Implement a Widget Factory

List the existing ADRs:
	adr list
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initGlobal)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.adr.yaml)")
	rootCmd.PersistentFlags().String("dir", "", "where to store the adrs")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose the outputs")
	viper.BindPFlag("ADRDir", rootCmd.PersistentFlags().Lookup("dir"))
	global = viper.New()
}

func initGlobal() {
	global.SetConfigName("config")
	global.AddConfigPath("$HOME/.config/adr")
	global.AddConfigPath("/etc/adr")
	global.SetDefault("template", template)
	global.SetDefault("baseTemplate", baseTemplate)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".adr" (without extension).
		viper.SetConfigName(".adr")
	}

	viper.AddConfigPath(".")
	viper.SetDefault("ADRDir", "docs/adr")
	viper.SetDefault("renderDir", "docs/")
	viper.SetDefault("Authors", []string{""})
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
