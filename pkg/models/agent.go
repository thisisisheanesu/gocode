package models

import (
	"github.com/fatih/color"
)

// AgentType represents the different types of specialized agents
type AgentType string

const (
	PlannerAgent   AgentType = "planner"
	FrontendAgent  AgentType = "frontend"
	BackendAgent   AgentType = "backend"
	DevOpsAgent    AgentType = "devops"
	ReviewerAgent  AgentType = "reviewer"
	ManagerAgent   AgentType = "manager"
	ToolsAgent     AgentType = "tools"
	ResearchAgent  AgentType = "research"
	SecurityAgent  AgentType = "security"
)

// Agent represents a specialized AI agent
type Agent interface {
	Name() string
	Type() AgentType
	Color() *color.Color
	Icon() string
	Role() string
	Process(context string, message string) (*Response, error)
	GetSystemPrompt() string
}

// AgentConfig holds configuration for an agent
type AgentConfig struct {
	Model       string `json:"model"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// Response represents an agent's response
type Response struct {
	Content     string            `json:"content"`
	TokensUsed  int              `json:"tokens_used"`
	Model       string            `json:"model"`
	Agent       AgentType         `json:"agent"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Actions     []Action          `json:"actions,omitempty"`
}

// Action represents an action an agent wants to perform
type Action struct {
	Type        ActionType  `json:"type"`
	Command     string      `json:"command,omitempty"`
	FilePath    string      `json:"file_path,omitempty"`
	Content     string      `json:"content,omitempty"`
	URL         string      `json:"url,omitempty"`
	Description string      `json:"description"`
	RequiresPermission bool `json:"requires_permission"`
}

// ActionType represents the type of action
type ActionType string

const (
	CommandAction    ActionType = "command"
	FileWriteAction  ActionType = "file_write"
	FileReadAction   ActionType = "file_read"
	WebSearchAction  ActionType = "web_search"
	ResearchAction   ActionType = "research"
)