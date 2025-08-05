package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var useGptOss120b bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-code",
	Short: "AI Development Team CLI - Orchestrate multiple AI agents for software development",
	Long: `go-code is a sophisticated CLI tool that orchestrates multiple specialized AI agents
for software development using the Groq API. Each agent has distinct roles and expertise:

🎯 @planner   - Project planning, task breakdown, architecture decisions
🎨 @frontend  - React, Vue, Angular, HTML/CSS, UI/UX decisions  
⚡ @backend   - APIs, databases, server logic, microservices
🚀 @devops    - CI/CD, Docker, Kubernetes, cloud deployment
🔍 @reviewer  - Code quality, security, best practices
👔 @manager   - Coordination, task delegation, progress tracking
🛠️ @tools     - Build tools, testing, debugging, utilities
📚 @research  - Documentation lookup, best practices, technology research
🛡️ @security  - Security audits, vulnerability assessment, secure coding

Use @agent syntax for auto-completion and direct agent communication.`,
	// Uncomment the following line if your bare application has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-code/config.json)")
	rootCmd.PersistentFlags().BoolVar(&useGptOss120b, "gpt-oss-120b", false, "Use the new OpenAI GPT-OSS-120B model")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".go-code" (without extension).
		viper.AddConfigPath(home + "/.go-code")
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// IsGptOss120bEnabled returns whether the --gpt-oss-120b flag is set
func IsGptOss120bEnabled() bool {
	return useGptOss120b
}