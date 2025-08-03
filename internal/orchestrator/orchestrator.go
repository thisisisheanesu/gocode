package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"go-code/internal/agents"
	"go-code/internal/filewriter"
	"go-code/internal/ui"
	"go-code/pkg/models"
)

// Orchestrator coordinates multiple agents to execute complex tasks
type Orchestrator struct {
	registry   *agents.Registry
	config     *models.Config
	fileWriter *filewriter.FileWriter
}

// New creates a new orchestrator
func New(registry *agents.Registry, config *models.Config) *Orchestrator {
	// Create project directory
	cwd, _ := os.Getwd()
	projectDir := filepath.Join(cwd, "generated-project")
	
	return &Orchestrator{
		registry:   registry,
		config:     config,
		fileWriter: filewriter.New(projectDir),
	}
}

// Task represents a task that needs to be executed by an agent
type Task struct {
	ID          string
	AgentType   models.AgentType
	Description string
	Dependencies []string
	Status      string
	Result      *models.Response
}

// ExecuteBuild coordinates agents to build a complete feature
func (o *Orchestrator) ExecuteBuild(description string) error {
	startTime := time.Now()
	
	// Clear screen for clean output
	ui.ClearScreen()
	fmt.Printf("🚀 Building: %s\n\n", description)
	
	// Step 0: Create project structure
	ui.DisplayProgress("Creating project structure", 0, 10, startTime)
	if err := o.fileWriter.CreateProjectStructure(); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}
	ui.DisplayStageComplete("Creating project structure", 1, 10, startTime)

	// Step 1: Get planner to create the plan
	ui.DisplayProgress("Planning project", 1, 10, startTime)
	planner, err := o.registry.GetAgent(models.PlannerAgent)
	if err != nil {
		return fmt.Errorf("failed to get planner agent: %w", err)
	}

	planPrompt := fmt.Sprintf(`Create a detailed execution plan for: "%s"

Please structure your response as a numbered list of tasks that can be executed by specialized agents.
For each task, specify:
1. The task description
2. Which agent should handle it (backend, frontend, security, etc.)
3. Any dependencies on other tasks

Format your response like this:
1. [AGENT_TYPE] Task description
2. [AGENT_TYPE] Task description (depends on task 1)
...

Available agents: backend, frontend, security, planner
Focus on creating actionable, specific tasks that agents can execute independently.
When providing code, use proper code blocks with filenames where possible.`, description)

	planResponse, err := planner.Process("", planPrompt)
	if err != nil {
		return fmt.Errorf("failed to create plan: %w", err)
	}
	ui.DisplayStageComplete("Planning project", 2, 10, startTime)

	// Step 2: Parse the plan into executable tasks
	tasks := o.parsePlan(planResponse.Content)
	if len(tasks) == 0 {
		return fmt.Errorf("no executable tasks found in plan")
	}

	totalStages := len(tasks) + 2 // +2 for structure creation and planning

	// Step 3: Execute tasks in order
	for i, task := range tasks {
		stageName := fmt.Sprintf("%s: %s", task.AgentType, task.Description[:min(40, len(task.Description))])
		ui.DisplayProgress(stageName, i+3, totalStages, startTime)
		
		agent, err := o.registry.GetAgent(task.AgentType)
		if err != nil {
			ui.DisplayWarning(fmt.Sprintf("Skipping task - agent %s not available: %v", task.AgentType, err))
			continue
		}
		
		// Create context from previous results
		context := o.buildContext(tasks[:i])
		
		response, err := agent.Process(context, task.Description)
		if err != nil {
			ui.DisplayError(fmt.Errorf("task failed: %w", err))
			continue
		}

		task.Result = response
		task.Status = "completed"
		
		// Extract and write any code blocks to files
		o.writeGeneratedFiles(response.Content)
		
		ui.DisplayStageComplete(stageName, i+3, totalStages, startTime)
	}

	// Final results
	cwd, _ := os.Getwd()
	projectPath := filepath.Join(cwd, "generated-project")
	ui.DisplayFinalResults(totalStages, startTime, projectPath)

	return nil
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// parsePlan extracts executable tasks from the planner's response
func (o *Orchestrator) parsePlan(planContent string) []Task {
	var tasks []Task
	
	// Regex to match task lines like "1. [BACKEND] Create API endpoints" or "1. **[BACKEND]** Create API endpoints"
	taskRegex := regexp.MustCompile(`(?i)^\d+\.\s*\*?\*?\[(\w+)\]\*?\*?\s*(.+?)(?:\s*` + "```" + `|\s*$)`)
	
	lines := strings.Split(planContent, "\n")
	taskID := 1
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		matches := taskRegex.FindStringSubmatch(line)
		if len(matches) >= 3 {
			agentName := strings.ToLower(matches[1])
			description := strings.TrimSpace(matches[2])
			
			// Clean up description - remove any trailing punctuation or formatting
			description = strings.TrimSuffix(description, ".")
			description = strings.TrimSuffix(description, ":")
			
			// Map agent names to types
			agentType := o.mapAgentName(agentName)
			if agentType == "" {
				continue
			}
			
			tasks = append(tasks, Task{
				ID:          fmt.Sprintf("task_%d", taskID),
				AgentType:   agentType,
				Description: description,
				Status:      "pending",
			})
			taskID++
		}
	}
	
	return tasks
}

// mapAgentName maps agent names from plan to agent types
func (o *Orchestrator) mapAgentName(name string) models.AgentType {
	switch strings.ToLower(name) {
	case "backend":
		return models.BackendAgent
	case "frontend":
		return models.FrontendAgent
	case "security":
		return models.SecurityAgent
	case "planner":
		return models.PlannerAgent
	default:
		return ""
	}
}

// buildContext creates context from previous task results
func (o *Orchestrator) buildContext(previousTasks []Task) string {
	if len(previousTasks) == 0 {
		return ""
	}
	
	var contextParts []string
	contextParts = append(contextParts, "Previous task results:")
	
	for _, task := range previousTasks {
		if task.Result != nil && task.Status == "completed" {
			contextParts = append(contextParts, fmt.Sprintf("\n- %s (%s): %s", 
				task.Description, task.AgentType, 
				o.truncateText(task.Result.Content, 200)))
		}
	}
	
	return strings.Join(contextParts, "\n")
}

// truncateText truncates text to a maximum length
func (o *Orchestrator) truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}

// writeGeneratedFiles extracts code blocks and writes them to files
func (o *Orchestrator) writeGeneratedFiles(content string) {
	codeBlocks := o.fileWriter.ExtractCodeBlocks(content)
	
	for filename, code := range codeBlocks {
		// Skip empty code blocks
		if strings.TrimSpace(code) == "" {
			continue
		}
		
		if err := o.fileWriter.WriteFile(filename, code); err != nil {
			fmt.Printf("⚠️  Failed to write %s: %v\n", filename, err)
		}
	}
}

