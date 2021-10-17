/*
Copyright © 2021 Jake Winters <j@jakew.ca>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/jakew/cargo/internal/build"
)

// DefaultCargo is the default Cargo YAML filename.
const DefaultCargo = "cargo.yaml"

// DefaultTemplate is the default template name to look for.
const DefaultTemplate = "template.Dockerfile"

// DefaultConfig is the default config file to look for.
const DefaultConfig = "config.yaml"

// DefaultOutput is the default filename to use when writing output.
const DefaultOutput = "Dockerfile"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cargo [TEMPLATE]",
	Short: "Render provided template",
	Long: `Cargo pre-processes the Dockerfile template passed in using
Go templating and Sprig functions. You can specify a YAML 
configuration file using -c to load data for the template.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("running with args: %s", strings.Join(args, " "))

		cargo := DefaultCargo
		if len(args) > 0 {
			cargo = args[0]
		}

		return build.BuildCargo(cargo)
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once to
// the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
