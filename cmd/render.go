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
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/jakew/cargo/internal/renderer"
)

var configFiles = []string{}
var output = ""

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render [TEMPLATE]",
	Short: "Render a specific Dockerfile template",
	Long: `Render a specific Dockerfile template. Multiple config files
can be set using multiple -c flags.`,
	Example: `render Dockerfile.template -c config.yaml`,
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := "template.Dockerfile"
		if len(args) > 0 {
			file = args[0]
		}

		var w io.Writer
		if output != "" {
			f, err := os.Create(output)
			if err != nil {
				return err
			}
			defer f.Close()
			w = f
		} else {
			w = cmd.OutOrStdout()
		}

		return renderer.PrintDockerfile(w, file, configFiles)
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringSliceVarP(&configFiles, "config", "c", []string{}, "Values that should be provided to the template.")
	renderCmd.MarkFlagFilename("config", "yaml", "yml")

	renderCmd.Flags().StringVarP(&output, "output", "o", "", "File to write the output to.")
}
