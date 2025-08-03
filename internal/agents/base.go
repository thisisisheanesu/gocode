package agents

import (
	"fmt"

	"github.com/fatih/color"
	"go-code/internal/api"
	"go-code/pkg/models"
)

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	agentType    models.AgentType
	name         string
	icon         string
	role         string
	color        *color.Color
	systemPrompt string
	client       *api.GroqClient
	config       models.AgentConfig
}

// NewBaseAgent creates a new base agent
func NewBaseAgent(
	agentType models.AgentType,
	name, icon, role, systemPrompt string,
	colorAttr color.Attribute,
	client *api.GroqClient,
	config models.AgentConfig,
) *BaseAgent {
	return &BaseAgent{
		agentType:    agentType,
		name:         name,
		icon:         icon,
		role:         role,
		color:        color.New(colorAttr),
		systemPrompt: systemPrompt,
		client:       client,
		config:       config,
	}
}

// Name returns the agent's name
func (a *BaseAgent) Name() string {
	return a.name
}

// Type returns the agent's type
func (a *BaseAgent) Type() models.AgentType {
	return a.agentType
}

// Color returns the agent's color
func (a *BaseAgent) Color() *color.Color {
	return a.color
}

// Icon returns the agent's icon
func (a *BaseAgent) Icon() string {
	return a.icon
}

// Role returns the agent's role description
func (a *BaseAgent) Role() string {
	return a.role
}

// GetSystemPrompt returns the agent's system prompt
func (a *BaseAgent) GetSystemPrompt() string {
	return a.systemPrompt
}

// Process sends a message to the agent and returns the response
func (a *BaseAgent) Process(context string, message string) (*models.Response, error) {
	fullMessage := message
	if context != "" {
		fullMessage = fmt.Sprintf("Context: %s\n\nUser Request: %s", context, message)
	}

	return a.client.ProcessAgentRequest(a.agentType, a.systemPrompt, fullMessage, a.config)
}