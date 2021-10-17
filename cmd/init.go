package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/jakew/cargo/internal/renderer"
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
		return renderer.WriteScaffolding(rootPath, DefaultCargo, DefaultTemplate, DefaultConfig, DefaultOutput)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
