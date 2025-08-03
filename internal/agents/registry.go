package agents

import (
	"fmt"
	"strings"

	"go-code/internal/api"
	"go-code/pkg/models"
)

// Registry manages all available agents
type Registry struct {
	agents map[models.AgentType]models.Agent
	client *api.GroqClient
	config *models.Config
}

// NewRegistry creates a new agent registry
func NewRegistry(client *api.GroqClient, config *models.Config) *Registry {
	registry := &Registry{
		agents: make(map[models.AgentType]models.Agent),
		client: client,
		config: config,
	}

	registry.initializeAgents()
	return registry
}

// initializeAgents creates all available agents
func (r *Registry) initializeAgents() {
	// Get agent configurations
	plannerConfig := r.getAgentConfig(models.PlannerAgent)
	frontendConfig := r.getAgentConfig(models.FrontendAgent)
	backendConfig := r.getAgentConfig(models.BackendAgent)
	securityConfig := r.getAgentConfig(models.SecurityAgent)

	// Create agents
	r.agents[models.PlannerAgent] = NewPlannerAgent(r.client, plannerConfig)
	r.agents[models.FrontendAgent] = NewFrontendAgent(r.client, frontendConfig)
	r.agents[models.BackendAgent] = NewBackendAgent(r.client, backendConfig)
	r.agents[models.SecurityAgent] = NewSecurityAgent(r.client, securityConfig)

	// TODO: Add remaining agents (DevOps, Reviewer, Manager, Tools, Research)
}

// getAgentConfig returns the configuration for a specific agent
func (r *Registry) getAgentConfig(agentType models.AgentType) models.AgentConfig {
	if config, exists := r.config.AgentPreferences[agentType]; exists {
		return config
	}
	
	// Return default config
	defaultConfig := models.DefaultConfig()
	if config, exists := defaultConfig.AgentPreferences[agentType]; exists {
		return config
	}
	
	// Fallback config
	return models.AgentConfig{
		Model:       r.config.DefaultModel,
		Temperature: 0.7,
		MaxTokens:   4096,
	}
}

// GetAgent returns an agent by type
func (r *Registry) GetAgent(agentType models.AgentType) (models.Agent, error) {
	agent, exists := r.agents[agentType]
	if !exists {
		return nil, fmt.Errorf("agent type %s not found", agentType)
	}
	return agent, nil
}

// GetAgentByName returns an agent by name (fuzzy matching)
func (r *Registry) GetAgentByName(name string) (models.Agent, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	
	// Direct match
	for _, agent := range r.agents {
		if strings.ToLower(agent.Name()) == name {
			return agent, nil
		}
	}
	
	// Fuzzy match
	for _, agent := range r.agents {
		agentName := strings.ToLower(agent.Name())
		if strings.HasPrefix(agentName, name) || strings.Contains(agentName, name) {
			return agent, nil
		}
	}
	
	return nil, fmt.Errorf("no agent found matching '%s'", name)
}

// ListAgents returns all available agents
func (r *Registry) ListAgents() []models.Agent {
	agents := make([]models.Agent, 0, len(r.agents))
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}
	return agents
}

// GetAgentNames returns all agent names for auto-completion
func (r *Registry) GetAgentNames() []string {
	names := make([]string, 0, len(r.agents))
	for _, agent := range r.agents {
		names = append(names, strings.ToLower(agent.Name()))
	}
	return names
}

// CompleteAgentName provides auto-completion suggestions
func (r *Registry) CompleteAgentName(partial string) []string {
	partial = strings.ToLower(strings.TrimSpace(partial))
	if partial == "" {
		return r.GetAgentNames()
	}
	
	var matches []string
	for _, agent := range r.agents {
		name := strings.ToLower(agent.Name())
		if strings.HasPrefix(name, partial) {
			matches = append(matches, name)
		}
	}
	
	// If no prefix matches, try fuzzy matching
	if len(matches) == 0 {
		for _, agent := range r.agents {
			name := strings.ToLower(agent.Name())
			if strings.Contains(name, partial) {
				matches = append(matches, name)
			}
		}
	}
	
	return matches
}