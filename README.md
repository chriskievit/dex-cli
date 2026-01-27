# DEX CLI - Azure DevOps CLI Tool

A secure Go-based CLI tool for managing Azure DevOps Git branches, work items, and pull requests.

[![Build](https://github.com/chriskievit/dex-cli/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/chriskievit/dex-cli/actions/workflows/build.yml)

## Features

- 🔐 **Secure Authentication**: Credentials stored in OS-native keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- 🌿 **Smart Branch Creation**: Create branches with automatic work item type detection
- 🔗 **Work Item Integration**: Automatic linking of branches and pull requests to work items
- 🚀 **Pull Request Management**: Create PRs with intelligent defaults
- 📝 **Work Item Viewing**: Display work item details directly from the CLI

## Installation

### Prerequisites

- Go 1.23 or higher
- Git installed and configured
- Azure DevOps account with a Personal Access Token (PAT)

### Build from Source

```bash
git clone https://github.com/chriskievit/dex-cli.git
cd dex-cli
go build -o dex
```

Move the binary to your PATH:

```bash
# macOS/Linux
sudo mv dex /usr/local/bin/

# Or add to your PATH
export PATH=$PATH:$(pwd)
```

## Configuration

### Authentication

First, authenticate with Azure DevOps:

```bash
dex auth login
```

You'll be prompted for:
- **Organization**: Your Azure DevOps organization name (e.g., `myorg` from `https://dev.azure.com/myorg`)
- **Personal Access Token (PAT)**: Your Azure DevOps PAT with appropriate permissions

Your PAT is stored securely in your system's keychain and never written to disk in plain text.

### Configuration File

The tool creates a configuration file at `~/.dex-cli/config.yaml`:

```yaml
organization: myorg
project: myproject
repository: myrepo
default_reviewer: ""
```

You can set configuration values using the `config set` commands, or edit the file directly.

## Usage

### Authentication Commands

```bash
# Login to Azure DevOps
dex auth login

# Check authentication status
dex auth status

# Logout (remove credentials)
dex auth logout
```

### Configuration Commands

```bash
# Show current configuration
dex config show

# Set project configuration value
dex config set project myproject

# Set repository configuration value
dex config set repo myrepository

# Set default reviewer configuration value
dex config set reviewer username@example.com
```

### Branch Management

Create a new branch linked to a work item:

```bash
# Create branch from default branch (main/master)
dex branch create 12345 add-login-feature

# Create branch from specific base branch
dex branch create 12345 fix-bug --from develop
```

**Branch Naming Convention**: `{work-item-type}/{work-item-id}/{description}`

Example: `user-story/12345/add-login-feature`

The work item type is automatically fetched from Azure DevOps.

### Pull Request Management

Create a pull request:

```bash
# Create PR from current branch
dex pr create --target main --title "Add login feature"

# Create PR with specific source branch
dex pr create --source feature/123/login --target main --title "Add login"

# Create draft PR
dex pr create --target main --title "WIP: New feature" --draft

# Specify work item manually
dex pr create --target main --title "Fix bug" --workitem 12345
```

**Smart Defaults**:
- Source branch defaults to your current Git branch
- Work item ID is automatically extracted from branch name if it follows the naming convention

### Work Item Commands

View work item details:

```bash
dex workitem show 12345
```

## Security Features

### Credential Storage

- **macOS**: Keychain Access
- **Windows**: Windows Credential Manager
- **Linux**: Secret Service (GNOME Keyring, KWallet)

Credentials are never:
- Stored in plain text files
- Logged to console or files
- Transmitted without HTTPS

### API Communication

- HTTPS only (TLS 1.2+)
- SSL certificate validation
- Proper timeout and retry logic
- No credential exposure in error messages

### Best Practices

1. **PAT Permissions**: Grant only necessary permissions to your PAT:
   - Code: Read & Write
   - Work Items: Read
   - Pull Requests: Read & Write

2. **PAT Expiration**: Set an expiration date for your PAT and rotate regularly

3. **Logout**: Use `dex auth logout` when done to remove credentials

## Examples

### Complete Workflow

```bash
# 1. Authenticate
dex auth login

# 2. Check work item
dex workitem show 12345

# 3. Create branch for the work item
dex branch create 12345 implement-new-feature

# 4. Make your changes, commit them
git add .
git commit -m "Implement new feature"

# 5. Push to remote
git push -u origin user-story/12345/implement-new-feature

# 6. Create pull request
dex pr create --target main --title "Implement new feature"
```

### Using Global Flags

```bash
# Override organization
dex --org myotherorg workitem show 12345

# Override project
dex --project myotherproject pr create --target main --title "Fix"

# Enable debug output
dex --debug branch create 12345 test-feature
```

## Troubleshooting

### "Not a git repository"

Make sure you're running commands from within a Git repository directory.

### "No credentials found"

Run `dex auth login` to authenticate first.

### "Failed to get repository"

Ensure your `~/.dex-cli/config.yaml` has the correct `repository` value set.

### "Invalid work item ID"

Verify the work item exists in your Azure DevOps project and you have access to it.

## Development

### Project Structure

```
dex-cli/
├── cmd/                 # Cobra commands
│   ├── root.go         # Root command
│   ├── auth.go         # Authentication commands
│   ├── branch.go       # Branch commands
│   ├── pr.go           # Pull request commands
│   └── workitem.go     # Work item commands
├── internal/
│   ├── auth/           # Credential management
│   ├── azdo/           # Azure DevOps API client
│   ├── config/         # Configuration management
│   └── git/            # Local Git operations
└── main.go             # Entry point
```

### Building

```bash
go build -o dex
```

### Testing

The project includes comprehensive unit tests for config, git modules, and command utility functions.

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests for a specific package
go test ./internal/config
go test ./internal/git
go test ./cmd
```

**Test Coverage:**
- `internal/config` - Configuration file operations
- `internal/git` - Git operations (using temporary repositories)
- `cmd/` - Command utility functions and handlers
  - `generateBranchDescription` - Work item title to branch name conversion
  - `isValidDescription` - Branch description validation
  - `extractWorkItemFromBranch` - Work item ID extraction from branch names
  - Config command handlers
  - Root command initialization

**Note:** Auth and Azure DevOps HTTP client modules are excluded from unit tests as they require external service dependencies.

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
