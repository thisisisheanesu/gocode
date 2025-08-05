package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go-code/internal/agents"
	"go-code/internal/api"
	"go-code/internal/config"
	"go-code/internal/ui"
	"go-code/pkg/models"
)

// chatCmd allows chatting with specific agents
var chatCmd = &cobra.Command{
	Use:   "chat [@agent] [message]",
	Short: "Chat with a specific AI agent",
	Long: `Chat with a specialized AI agent using @agent syntax.

Available agents:
  @planner   - Project planning and architecture
  @frontend  - UI/UX and frontend development  
  @backend   - APIs and server-side development
  @security  - Security audits and secure coding
  @devops    - Deployment and infrastructure
  @reviewer  - Code review and quality
  @manager   - Team coordination
  @tools     - Build tools and debugging
  @research  - Documentation and best practices

Examples:
  go-code chat @planner "Plan a microservices architecture"
  go-code chat @frontend "Create a responsive navbar component"
  go-code chat @security "Review this authentication function"`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		manager := config.NewManager()
		if err := manager.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		cfg := manager.GetConfig()
		
		// Override model if --gpt-oss-120b flag is set
		if IsGptOss120bEnabled() {
			// Create a copy of the config and override all agent models
			configCopy := *cfg
			agentPrefs := make(map[models.AgentType]models.AgentConfig)
			for agentType, agentConfig := range cfg.AgentPreferences {
				newConfig := agentConfig
				newConfig.Model = "openai/gpt-oss-120b"
				agentPrefs[agentType] = newConfig
			}
			configCopy.AgentPreferences = agentPrefs
			configCopy.DefaultModel = "openai/gpt-oss-120b"
			cfg = &configCopy
		}
		
		if err := manager.ValidateConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
			os.Exit(1)
		}

		// Parse agent and message
		agentName := strings.TrimPrefix(args[0], "@")
		message := strings.Join(args[1:], " ")

		if agentName == "" {
			fmt.Fprintf(os.Stderr, "Please specify an agent using @agent syntax\n")
			os.Exit(1)
		}

		if message == "" {
			fmt.Fprintf(os.Stderr, "Please provide a message\n")
			os.Exit(1)
		}

		// Create client and registry
		client := api.NewGroqClient(cfg.GroqAPIKey)
		registry := agents.NewRegistry(client, cfg)

		// Get agent
		agent, err := registry.GetAgentByName(agentName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			fmt.Fprintf(os.Stderr, "\nAvailable agents: %v\n", strings.Join(registry.GetAgentNames(), ", "))
			os.Exit(1)
		}

		// Display agent header
		ui.DisplayAgentHeader(agent)

		// Process message
		fmt.Println("💭 Thinking...")
		response, err := agent.Process("", message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing request: %v\n", err)
			os.Exit(1)
		}

		// Display response
		ui.DisplayAgentResponse(agent, response)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}