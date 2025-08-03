package models

// Config represents the application configuration
type Config struct {
	GroqAPIKey              string                   `json:"groq_api_key"`
	DefaultModel            string                   `json:"default_model"`
	AllowedCommands         []string                 `json:"allowed_commands"`
	RequireCommandPermission bool                    `json:"require_command_permission"`
	RestrictToCurrentDir    bool                     `json:"restrict_to_current_dir"`
	AgentPreferences        map[AgentType]AgentConfig `json:"agent_preferences"`
	WorkingDirectory        string                   `json:"working_directory"`
	SessionPermissions      map[string]bool          `json:"session_permissions"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		DefaultModel:            "llama-3.1-70b-versatile",
		AllowedCommands:         []string{"npm", "go", "docker", "git"},
		RequireCommandPermission: true,
		RestrictToCurrentDir:    true,
		AgentPreferences: map[AgentType]AgentConfig{
			PlannerAgent: {
				Model:       "llama-3.1-70b-versatile",
				Temperature: 0.7,
				MaxTokens:   4096,
			},
			FrontendAgent: {
				Model:       "qwen2.5-coder-32b-instruct",
				Temperature: 0.3,
				MaxTokens:   4096,
			},
			BackendAgent: {
				Model:       "qwen2.5-coder-32b-instruct",
				Temperature: 0.3,
				MaxTokens:   4096,
			},
			DevOpsAgent: {
				Model:       "llama-3.1-70b-versatile",
				Temperature: 0.5,
				MaxTokens:   4096,
			},
			ReviewerAgent: {
				Model:       "llama-3.1-70b-versatile",
				Temperature: 0.3,
				MaxTokens:   4096,
			},
			ManagerAgent: {
				Model:       "llama-3.1-70b-versatile",
				Temperature: 0.6,
				MaxTokens:   4096,
			},
			ToolsAgent: {
				Model:       "qwen2.5-coder-32b-instruct",
				Temperature: 0.2,
				MaxTokens:   4096,
			},
			ResearchAgent: {
				Model:       "llama-3.1-70b-versatile",
				Temperature: 0.4,
				MaxTokens:   4096,
			},
			SecurityAgent: {
				Model:       "llama-3.1-70b-versatile",
				Temperature: 0.2,
				MaxTokens:   4096,
			},
		},
		SessionPermissions: make(map[string]bool),
	}
}