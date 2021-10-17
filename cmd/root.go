// Package cmd handles the CLI commands.
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/jakew/cargo/internal/renderer"
)

// DefaultCargo is the default Cargo YAML filename.
const DefaultCargo = "cargo.yaml"

// DefaultTemplate is the default template name to look for.
const DefaultTemplate = "template.Dockerfile"

// DefaultConfig is the default config file to look for.
const DefaultConfig = "config.yaml"

// DefaultOutput is the default filename to use when writing output.
const DefaultOutput = "Dockerfile"

// manifestNames are the names of the manifests to render.
var manifestNames = []string{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cargo [TEMPLATE]",
	Short: "Render provided template",
	Long: `Cargo pre-processes the Dockerfile template passed in using
Go templating and Sprig functions. You can specify a YAML 
configuration file using -c to load data for the template.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cargo := DefaultCargo
		if len(args) > 0 {
			cargo = args[0]
		}

		return renderer.RenderCargo(cargo, manifestNames)
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

	rootCmd.Flags().StringSliceVarP(&manifestNames, "manifest", "m", []string{}, "Manifesets to render. (default: all)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
