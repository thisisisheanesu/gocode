package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go-code/internal/config"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize go-code with default configuration",
	Long: `Initialize go-code by creating a configuration file with default settings.
This command creates ~/.go-code/config.json with default agent configurations.

You'll need to set your Groq API key after initialization:
  go-code config set-key YOUR_API_KEY`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := config.NewManager()
		
		if err := manager.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ go-code initialized successfully!")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("1. Set your Groq API key:")
		fmt.Println("   go-code config set-key YOUR_API_KEY")
		fmt.Println()
		fmt.Println("2. List available agents:")
		fmt.Println("   go-code agents")
		fmt.Println()
		fmt.Println("3. Start chatting with an agent:")
		fmt.Println("   go-code chat @planner \"Help me plan a web application\"")
		fmt.Println()
		fmt.Println("Configuration saved to:", fmt.Sprintf("%s/.go-code/config.json", os.Getenv("HOME")))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}