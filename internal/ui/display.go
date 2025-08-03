package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"go-code/pkg/models"
)

// DisplayAgentHeader shows the agent information before processing
func DisplayAgentHeader(agent models.Agent) {
	fmt.Println()
	
	// Agent info with color and icon
	agentColor := agent.Color()
	agentColor.Printf("%s %s Agent\n", agent.Icon(), agent.Name())
	
	// Role description in gray
	gray := color.New(color.FgHiBlack)
	gray.Printf("   %s\n", agent.Role())
	
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()
}

// DisplayAgentResponse shows the agent's response with formatting
func DisplayAgentResponse(agent models.Agent, response *models.Response) {
	fmt.Println()
	
	// Response header
	agentColor := agent.Color()
	agentColor.Printf("%s %s:\n", agent.Icon(), agent.Name())
	fmt.Println()
	
	// Response content
	fmt.Println(response.Content)
	fmt.Println()
	
	// Metadata
	displayResponseMetadata(response)
	fmt.Println()
}

// displayResponseMetadata shows token usage and model info
func displayResponseMetadata(response *models.Response) {
	gray := color.New(color.FgHiBlack)
	
	metadata := []string{}
	
	if response.Model != "" {
		metadata = append(metadata, fmt.Sprintf("Model: %s", response.Model))
	}
	
	if response.TokensUsed > 0 {
		metadata = append(metadata, fmt.Sprintf("Tokens: %d", response.TokensUsed))
	}
	
	if len(metadata) > 0 {
		gray.Printf("📊 %s\n", strings.Join(metadata, " • "))
	}
}

// DisplayThinking shows a thinking animation
func DisplayThinking(agent models.Agent) {
	agentColor := agent.Color()
	thinkingChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	
	for i := 0; i < 10; i++ {
		char := thinkingChars[i%len(thinkingChars)]
		agentColor.Printf("\r%s %s is thinking... %s", agent.Icon(), agent.Name(), char)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Print("\r" + strings.Repeat(" ", 50) + "\r") // Clear line
}

// DisplayError shows an error message
func DisplayError(err error) {
	red := color.New(color.FgRed, color.Bold)
	red.Printf("❌ Error: %v\n", err)
}

// DisplaySuccess shows a success message
func DisplaySuccess(message string) {
	green := color.New(color.FgGreen, color.Bold)
	green.Printf("✅ %s\n", message)
}

// DisplayWarning shows a warning message
func DisplayWarning(message string) {
	yellow := color.New(color.FgYellow, color.Bold)
	yellow.Printf("⚠️  %s\n", message)
}

// DisplayInfo shows an info message
func DisplayInfo(message string) {
	blue := color.New(color.FgBlue)
	blue.Printf("ℹ️  %s\n", message)
}

// DisplayAgentList shows a formatted list of agents
func DisplayAgentList(agents []models.Agent) {
	fmt.Println("🤖 Available Agents:")
	fmt.Println()
	
	for _, agent := range agents {
		agentColor := agent.Color()
		agentColor.Printf("  %s @%s\n", agent.Icon(), strings.ToLower(agent.Name()))
		
		gray := color.New(color.FgHiBlack)
		gray.Printf("     %s\n", agent.Role())
		fmt.Println()
	}
}

// DisplayHelp shows help information
func DisplayHelp() {
	fmt.Println("🚀 go-code - AI Development Team CLI")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init                    Initialize configuration")
	fmt.Println("  agents                  List available agents")
	fmt.Println("  chat @agent message     Chat with an agent")
	fmt.Println("  config                  Manage configuration")
	fmt.Println()
	fmt.Println("Agent Examples:")
	fmt.Println("  go-code chat @planner \"Plan a web app\"")
	fmt.Println("  go-code chat @frontend \"Create React component\"")
	fmt.Println("  go-code chat @backend \"Design REST API\"")
	fmt.Println("  go-code chat @security \"Review this code\"")
	fmt.Println()
	fmt.Println("For more help: go-code --help")
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// DisplayProgress shows progress with stage replacement
func DisplayProgress(stage string, current, total int, startTime time.Time) {
	// Clear the line and move cursor to beginning
	fmt.Print("\r\033[K")
	
	elapsed := time.Since(startTime)
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	
	// Progress bar
	progress := float64(current) / float64(total) * 100
	progressBar := strings.Repeat("█", int(progress/5)) + strings.Repeat("░", 20-int(progress/5))
	
	cyan.Printf("🔄 Stage %d/%d: %s ", current, total, stage)
	fmt.Printf("[%s] %.0f%% ", progressBar, progress)
	yellow.Printf("(⏱️ %s)", elapsed.Round(time.Second))
}

// DisplayStageComplete shows stage completion
func DisplayStageComplete(stage string, current, total int, startTime time.Time) {
	elapsed := time.Since(startTime)
	green := color.New(color.FgGreen, color.Bold)
	
	fmt.Print("\r\033[K")
	green.Printf("✅ Stage %d/%d: %s completed (⏱️ %s)\n", current, total, stage, elapsed.Round(time.Second))
}

// DisplayFinalResults shows final completion with total time
func DisplayFinalResults(totalStages int, startTime time.Time, projectPath string) {
	totalTime := time.Since(startTime)
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan)
	
	fmt.Println()
	fmt.Println(strings.Repeat("═", 60))
	green.Printf("🎉 BUILD COMPLETED SUCCESSFULLY! \n")
	fmt.Println(strings.Repeat("═", 60))
	fmt.Printf("📊 Total stages: %d\n", totalStages)
	fmt.Printf("⏱️  Total time: %s\n", totalTime.Round(time.Second))
	fmt.Printf("📁 Project location: %s\n", projectPath)
	fmt.Println()
	cyan.Println("Manual setup required:")
	fmt.Println("  cd generated-project")
	fmt.Println("  npm install           # Install dependencies") 
	fmt.Println("  # Set up database (if needed)")
	fmt.Println("  npm start            # Start application")
	fmt.Println()
	fmt.Println("💡 go-code only generates files - you must run setup commands manually")
}