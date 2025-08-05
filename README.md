# go-code - AI Development Team CLI

A sophisticated CLI tool that orchestrates multiple specialized AI agents for software development using the Groq API. Think of it as having an entire development team available through your command line.

## ⚠️ **BETA SOFTWARE - USE WITH CAUTION**

**This is newly created software and may contain bugs.** The file generation feature is experimental and may:
- Overwrite existing files without warning
- Generate incorrect file structures  
- Create malformed code that needs manual fixing
- Fail to parse complex code blocks properly

**Always run in a separate directory and review generated code before using in production.**

## 🎯 What go-code Does vs. Doesn't Do

### ✅ **What go-code DOES:**
- Generates complete project file structures
- Creates working code files (JS, Python, Go, SQL, etc.)
- Sets up proper directory organization
- Provides multi-agent AI collaboration
- Writes package.json, requirements.txt, go.mod files
- Creates database schemas and configuration files

### ❌ **What go-code does NOT do:**
- Install npm packages or dependencies
- Run build or compilation commands  
- Execute database migrations or setup
- Start servers or applications
- Configure environment variables
- Deploy to production environments

**Think of go-code as a "file generator" - it creates the code, you run the setup.**

## 🚀 Features

- **9 Specialized AI Agents** - Each with distinct expertise and colored output
- **@agent Syntax** - Natural chat interface with auto-completion
- **Groq API Integration** - Fast, reliable AI responses with multiple model options
- **Configurable Permissions** - Control command execution and security
- **Session Management** - Persistent settings and context sharing
- **Colorized Output** - Beautiful terminal interface with agent-specific colors

## 🤖 Available Agents

| Agent | Icon | Specialization | Color |
|-------|------|----------------|-------|
| **@planner** | 🎯 | Project planning, task breakdown, architecture decisions | Blue |
| **@frontend** | 🎨 | React, Vue, Angular, HTML/CSS, UI/UX decisions | Magenta |
| **@backend** | ⚡ | APIs, databases, server logic, microservices | Green |
| **@devops** | 🚀 | CI/CD, Docker, Kubernetes, cloud deployment | Cyan |
| **@reviewer** | 🔍 | Code quality, security, best practices | Yellow |
| **@manager** | 👔 | Coordination, task delegation, progress tracking | Red |
| **@tools** | 🛠️ | Build tools, testing, debugging, utilities | White |
| **@research** | 📚 | Documentation lookup, best practices, research | Gray |
| **@security** | 🛡️ | Security audits, vulnerability assessment, secure coding | Bright Red |

## 📦 Installation

### Prerequisites
- Go 1.22 or higher
- Groq API key (get one at [console.groq.com](https://console.groq.com))

### Build from Source
```bash
git clone <repository-url>
cd go-code
go build -o go-code
./go-code init
```

### Set Up API Key
```bash
./go-code config set-key YOUR_GROQ_API_KEY
```

## 🛠️ Usage

### Initialize Configuration
```bash
go-code init
```

### List Available Agents
```bash
go-code agents
```

### 📁 After Project Generation

**go-code only creates files** - you need to manually set up the generated project:

```bash
# Navigate to generated project
cd generated-project/

# Install dependencies (for Node.js projects)
npm install
# or for Python projects
pip install -r requirements.txt
# or for Go projects  
go mod tidy

# Set up database (if applicable)
# - Create database
# - Run schema.sql file
# - Configure connection strings

# Start the application
npm start
# or
node app.js
# or  
go run main.go
```

### Chat with Agents
```bash
# Plan a project
go-code chat @planner "Plan a microservices architecture for an e-commerce platform"

# Get frontend help
go-code chat @frontend "Create a responsive React navbar component with dark mode"

# Backend assistance
go-code chat @backend "Design a REST API for user authentication with JWT"

# Security review
go-code chat @security "Review this authentication middleware for vulnerabilities"

# DevOps guidance
go-code chat @devops "Set up CI/CD pipeline for a Node.js application"
```

### Auto-Build Projects (NEW!)
```bash
# Automatically plan and generate a complete project
go-code build "a todo app with React frontend and Node.js backend"

# Generate a REST API with authentication
go-code build "blog management API with user authentication"

# Create a full-stack application
go-code build "e-commerce platform with product catalog and shopping cart"
```

**The `build` command will:**
1. 🎯 **Planner** creates a detailed execution plan
2. 📁 Creates proper project directory structure  
3. ⚡ **Backend/Frontend/Security** agents execute their assigned tasks
4. 📄 **Automatically writes generated code to files** in `generated-project/`
5. 🔄 Shares context between agents for coherent results

**⚠️ IMPORTANT: go-code only generates files - it does NOT:**
- Run `npm install` or package installation commands
- Execute setup scripts or build processes  
- Start servers or applications
- Install dependencies or libraries

**You must manually run setup commands after generation.**

### Configuration Management
```bash
# Show current configuration
go-code config show

# Set default model
go-code config set-model llama-3.1-70b-versatile

# Toggle command execution permissions
go-code config allow-commands
```

## ⚙️ Configuration

Configuration is stored in `~/.go-code/config.json`:

```json
{
  "groq_api_key": "your-api-key",
  "default_model": "llama-3.1-70b-versatile",
  "allowed_commands": ["npm", "go", "docker", "git"],
  "require_command_permission": true,
  "restrict_to_current_dir": true,
  "agent_preferences": {
    "planner": {
      "model": "llama-3.1-70b-versatile",
      "temperature": 0.7,
      "max_tokens": 4096
    },
    "frontend": {
      "model": "qwen2.5-coder-32b-instruct",
      "temperature": 0.3,
      "max_tokens": 4096
    },
    "backend": {
      "model": "qwen2.5-coder-32b-instruct",
      "temperature": 0.3,
      "max_tokens": 4096
    },
    "security": {
      "model": "llama-3.1-70b-versatile",
      "temperature": 0.2,
      "max_tokens": 4096
    }
  }
}
```

## 🎯 Available Models

- `llama-3.1-70b-versatile` - Best for planning, reasoning, and complex tasks
- `llama-3.1-8b-instant` - Fast responses for simple queries
- `qwen2.5-coder-32b-instruct` - Optimized for code generation
- `mixtral-8x7b-32768` - Good balance of speed and capability
- `gemma2-9b-it` - Efficient for general tasks
- `openai/gpt-oss-120b` - OpenAI's new open-source model (use with `--gpt-oss-120b` flag)

### Using OpenAI GPT-OSS-120B Model

To use the new OpenAI GPT-OSS-120B model, add the `--gpt-oss-120b` flag to any command:

```bash
# Chat with agents using the new OpenAI model
go-code --gpt-oss-120b chat @planner "Plan a microservices architecture"
go-code --gpt-oss-120b chat @frontend "Create a React component"

# Build projects with the new OpenAI model
go-code --gpt-oss-120b build "a todo app with React frontend"
```

The `--gpt-oss-120b` flag overrides all agent model configurations to use `openai/gpt-oss-120b` via Groq's API.

## 🔒 Security Features

### Defensive Security Only
The security agent is designed exclusively for defensive purposes:
- ✅ Security audits and vulnerability assessments
- ✅ Secure coding best practices
- ✅ Code review for security flaws
- ❌ No offensive security tools or techniques
- ❌ No assistance with bypassing security measures

### Command Execution Controls
- Permission-based system for running commands
- Configurable allowed command whitelist
- Directory restriction options
- Session-based permission caching

## 🏗️ Project Structure

```
go-code/
├── main.go                 # Application entry point
├── cmd/                    # Cobra commands
│   ├── root.go            # Root command and configuration
│   ├── init.go            # Initialization command
│   ├── chat.go            # Chat command
│   ├── agents.go          # Agent listing
│   └── config.go          # Configuration management
├── internal/
│   ├── agents/            # AI agent implementations
│   │   ├── base.go        # Base agent functionality
│   │   ├── planner.go     # Planning agent
│   │   ├── frontend.go    # Frontend agent
│   │   ├── backend.go     # Backend agent
│   │   ├── security.go    # Security agent
│   │   └── registry.go    # Agent management
│   ├── api/               # Groq API client
│   │   └── groq.go        # API client implementation
│   ├── config/            # Configuration management
│   │   └── manager.go     # Config loading/saving
│   └── ui/                # Terminal UI components
│       └── display.go     # Formatted output
└── pkg/
    └── models/            # Data models
        ├── agent.go       # Agent interfaces
        └── config.go      # Configuration models
```

## 🚧 Development Roadmap

### Current Status ✅
- [x] Core CLI framework with Cobra
- [x] Groq API client with error handling
- [x] Basic agent system (Planner, Frontend, Backend, Security)
- [x] Configuration management with Viper
- [x] Colorized terminal output
- [x] Command structure and help system

### In Progress 🚧
- [ ] Auto-completion system for @agent references
- [ ] Remaining agents (DevOps, Reviewer, Manager, Tools, Research)
- [ ] Command execution system with permissions

### Future Features 🔮
- [ ] Session management and context sharing
- [ ] Multi-agent collaboration workflows
- [ ] Internet research capability
- [ ] Project integration (Git, package.json detection)
- [ ] Agent-to-agent communication
- [ ] Shell completion (bash, zsh, fish)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://cobra.dev/) for CLI framework
- Powered by [Groq](https://groq.com/) for fast AI inference
- Uses [Viper](https://github.com/spf13/viper) for configuration management
- Terminal colors by [fatih/color](https://github.com/fatih/color)

---

**Made with ❤️ for developers who want AI assistance without leaving the terminal.**