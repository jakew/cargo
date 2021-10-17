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
