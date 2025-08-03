package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go-code/internal/config"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage go-code configuration",
	Long:  `Configure go-code settings including API keys, models, and agent preferences.`,
}

// setKeyCmd sets the Groq API key
var setKeyCmd = &cobra.Command{
	Use:   "set-key [api-key]",
	Short: "Set Groq API key",
	Long:  `Set your Groq API key for accessing the Groq API.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := config.NewManager()
		if err := manager.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		apiKey := strings.TrimSpace(args[0])
		if apiKey == "" {
			fmt.Fprintf(os.Stderr, "API key cannot be empty\n")
			os.Exit(1)
		}

		if err := manager.SetGroqAPIKey(apiKey); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting API key: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Groq API key set successfully!")
	},
}

// setModelCmd sets the default model
var setModelCmd = &cobra.Command{
	Use:   "set-model [model]",
	Short: "Set default model",
	Long: `Set the default model for agents. Available models:
- llama-3.1-70b-versatile
- llama-3.1-8b-instant
- mixtral-8x7b-32768
- gemma2-9b-it
- qwen2.5-coder-32b-instruct`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := config.NewManager()
		if err := manager.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		model := strings.TrimSpace(args[0])
		if err := manager.SetDefaultModel(model); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting model: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Default model set to: %s\n", model)
	},
}

// allowCommandsCmd toggles command execution permissions
var allowCommandsCmd = &cobra.Command{
	Use:   "allow-commands",
	Short: "Toggle command execution permissions",
	Long:  `Toggle whether go-code agents can execute system commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := config.NewManager()
		if err := manager.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		config := manager.GetConfig()
		newSetting := !config.RequireCommandPermission
		
		if err := manager.SetAllowCommands(!newSetting); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating settings: %v\n", err)
			os.Exit(1)
		}

		if newSetting {
			fmt.Println("✅ Command execution enabled (agents can run system commands)")
		} else {
			fmt.Println("✅ Command execution disabled (agents will ask for permission)")
		}
	},
}

// showCmd shows current configuration
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current go-code configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := config.NewManager()
		if err := manager.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		config := manager.GetConfig()
		
		fmt.Println("🔧 go-code Configuration")
		fmt.Println("=" + strings.Repeat("=", 25))
		fmt.Println()
		
		// API Key
		if config.GroqAPIKey != "" {
			maskedKey := maskAPIKey(config.GroqAPIKey)
			color.Green("✅ Groq API Key: %s", maskedKey)
		} else {
			color.Red("❌ Groq API Key: Not set")
		}
		
		// Default Model
		fmt.Printf("🤖 Default Model: %s\n", config.DefaultModel)
		
		// Command Permissions
		if config.RequireCommandPermission {
			color.Yellow("🔒 Command Execution: Requires permission")
		} else {
			color.Green("🔓 Command Execution: Allowed")
		}
		
		// Working Directory
		if config.WorkingDirectory != "" {
			fmt.Printf("📁 Working Directory: %s\n", config.WorkingDirectory)
		}
		
		// Allowed Commands
		if len(config.AllowedCommands) > 0 {
			fmt.Printf("✅ Allowed Commands: %s\n", strings.Join(config.AllowedCommands, ", "))
		}
		
		fmt.Println()
		fmt.Printf("📄 Config file: %s/.go-code/config.json\n", os.Getenv("HOME"))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setKeyCmd)
	configCmd.AddCommand(setModelCmd)
	configCmd.AddCommand(allowCommandsCmd)
	configCmd.AddCommand(showCmd)
}

// maskAPIKey masks the API key for display
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}