package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go-code/internal/agents"
	"go-code/internal/api"
	"go-code/internal/config"
)

// agentsCmd lists all available agents
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "List all available agents",
	Long: `Display information about all available AI agents and their specializations.
Use @agent syntax to interact with specific agents.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := config.NewManager()
		if err := manager.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		cfg := manager.GetConfig()
		if cfg.GroqAPIKey == "" {
			fmt.Fprintf(os.Stderr, "❌ Groq API key not set. Use 'go-code config set-key YOUR_API_KEY' to set it.\n")
			os.Exit(1)
		}

		// Create client and registry
		client := api.NewGroqClient(cfg.GroqAPIKey)
		registry := agents.NewRegistry(client, cfg)
		
		fmt.Println("🤖 Available AI Development Agents")
		fmt.Println("=" + strings.Repeat("=", 35))
		fmt.Println()
		
		agentList := registry.ListAgents()
		for _, agent := range agentList {
			// Use agent's color for the output
			color := agent.Color()
			color.Printf("%s @%s\n", agent.Icon(), strings.ToLower(agent.Name()))
			fmt.Printf("   %s\n", agent.Role())
			fmt.Println()
		}
		
		fmt.Println("Usage Examples:")
		fmt.Println("  go-code chat @planner \"Plan a web application with authentication\"")
		fmt.Println("  go-code chat @frontend \"Create a React component for user login\"")
		fmt.Println("  go-code chat @backend \"Design a REST API for user management\"")
		fmt.Println("  go-code chat @security \"Review this authentication code\"")
		fmt.Println()
		fmt.Println("💡 Tip: Use tab completion for agent names after @")
	},
}

func init() {
	rootCmd.AddCommand(agentsCmd)
}