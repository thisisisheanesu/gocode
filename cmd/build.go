package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go-code/internal/agents"
	"go-code/internal/api"
	"go-code/internal/config"
	"go-code/internal/orchestrator"
	"go-code/internal/ui"
	"go-code/pkg/models"
)

// buildCmd auto-coordinates agents to build a feature
var buildCmd = &cobra.Command{
	Use:   "build [description]",
	Short: "Auto-coordinate agents to build a feature",
	Long: `Automatically coordinate multiple agents to plan and build a complete feature.
The planner agent will create a detailed plan, then delegate tasks to appropriate 
specialized agents who will execute their parts automatically.

Examples:
  go-code build "a todo app with React frontend and Node.js backend"
  go-code build "user authentication system with JWT"
  go-code build "REST API for blog management"`,
	Args: cobra.MinimumNArgs(1),
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

		description := strings.Join(args, " ")

		// Create client and registry
		client := api.NewGroqClient(cfg.GroqAPIKey)
		registry := agents.NewRegistry(client, cfg)

		// Create orchestrator
		orch := orchestrator.New(registry, cfg)

		// Display starting message
		ui.DisplayInfo(fmt.Sprintf("🚀 Building: %s", description))
		fmt.Println()

		// Execute the build workflow
		if err := orch.ExecuteBuild(description); err != nil {
			ui.DisplayError(fmt.Errorf("build failed: %w", err))
			os.Exit(1)
		}

		ui.DisplaySuccess("Build completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}