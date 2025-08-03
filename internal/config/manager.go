package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go-code/pkg/models"
)

// Manager handles configuration loading and saving
type Manager struct {
	configPath string
	config     *models.Config
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".go-code")
	configPath := filepath.Join(configDir, "config.json")
	
	return &Manager{
		configPath: configPath,
		config:     models.DefaultConfig(),
	}
}

// Load loads the configuration from file
func (m *Manager) Load() error {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		// Create default config file
		return m.Save()
	}

	// Read config file
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse config
	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Merge with defaults to ensure all fields are present
	m.mergeWithDefaults(&config)
	m.config = &config

	return nil
}

// Save saves the configuration to file
func (m *Manager) Save() error {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func (m *Manager) GetConfig() *models.Config {
	return m.config
}

// SetGroqAPIKey sets the Groq API key
func (m *Manager) SetGroqAPIKey(apiKey string) error {
	m.config.GroqAPIKey = apiKey
	return m.Save()
}

// SetDefaultModel sets the default model
func (m *Manager) SetDefaultModel(model string) error {
	m.config.DefaultModel = model
	return m.Save()
}

// SetAllowCommands toggles command execution permissions
func (m *Manager) SetAllowCommands(allow bool) error {
	m.config.RequireCommandPermission = !allow
	return m.Save()
}

// AddAllowedCommand adds a command to the allowed list
func (m *Manager) AddAllowedCommand(command string) error {
	for _, cmd := range m.config.AllowedCommands {
		if cmd == command {
			return nil // Already exists
		}
	}
	m.config.AllowedCommands = append(m.config.AllowedCommands, command)
	return m.Save()
}

// RemoveAllowedCommand removes a command from the allowed list
func (m *Manager) RemoveAllowedCommand(command string) error {
	for i, cmd := range m.config.AllowedCommands {
		if cmd == command {
			m.config.AllowedCommands = append(m.config.AllowedCommands[:i], m.config.AllowedCommands[i+1:]...)
			break
		}
	}
	return m.Save()
}

// SetAgentConfig sets configuration for a specific agent
func (m *Manager) SetAgentConfig(agentType models.AgentType, config models.AgentConfig) error {
	if m.config.AgentPreferences == nil {
		m.config.AgentPreferences = make(map[models.AgentType]models.AgentConfig)
	}
	m.config.AgentPreferences[agentType] = config
	return m.Save()
}

// ValidateConfig validates the configuration
func (m *Manager) ValidateConfig() error {
	if m.config.GroqAPIKey == "" {
		return fmt.Errorf("Groq API key is required. Use 'go-code config set-key <key>' to set it")
	}

	if m.config.DefaultModel == "" {
		return fmt.Errorf("default model is required")
	}

	// Validate models in agent preferences
	for agentType, config := range m.config.AgentPreferences {
		if config.Model != "" && !isValidModel(config.Model) {
			return fmt.Errorf("invalid model '%s' for agent %s", config.Model, agentType)
		}
	}

	return nil
}

// mergeWithDefaults merges the loaded config with default values
func (m *Manager) mergeWithDefaults(config *models.Config) {
	defaults := models.DefaultConfig()

	if config.DefaultModel == "" {
		config.DefaultModel = defaults.DefaultModel
	}

	if len(config.AllowedCommands) == 0 {
		config.AllowedCommands = defaults.AllowedCommands
	}

	if config.AgentPreferences == nil {
		config.AgentPreferences = defaults.AgentPreferences
	} else {
		// Merge agent preferences with defaults
		for agentType, defaultConfig := range defaults.AgentPreferences {
			if _, exists := config.AgentPreferences[agentType]; !exists {
				config.AgentPreferences[agentType] = defaultConfig
			}
		}
	}

	if config.SessionPermissions == nil {
		config.SessionPermissions = make(map[string]bool)
	}
}

// isValidModel checks if a model is valid
func isValidModel(model string) bool {
	validModels := []string{
		// Production models
		"gemma2-9b-it",
		"llama-3.1-8b-instant",
		"llama-3.3-70b-versatile",
		"meta-llama/llama-guard-4-12b",
		// Preview models
		"deepseek-r1-distill-llama-70b",
		"meta-llama/llama-4-maverick-17b-128e-instruct",
		"meta-llama/llama-4-scout-17b-16e-instruct",
		"moonshotai/kimi-k2-instruct",
		"qwen/qwen3-32b",
	}

	for _, validModel := range validModels {
		if model == validModel {
			return true
		}
	}
	return false
}