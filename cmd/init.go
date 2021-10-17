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

	"github.com/spf13/cobra"

	"github.com/jakew/cargo/internal/build"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [PATH]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Create sample Cargo scaffolding.",
	Long: `Create a sample Cargo scaffold at the path provided. If no
path is provided, the path of . is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rootPath := "./"
		if len(args) == 1 {
			rootPath = args[0]
		}

		log.Printf("initializing scaffolding in %s", rootPath)
		return build.WriteScaffolding(rootPath, DefaultCargo, DefaultTemplate, DefaultConfig, DefaultOutput)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}