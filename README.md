# Gemi CLI

A beautiful command-line interface tool built with Go and Cobra, enhanced with various libraries to make it visually appealing and user-friendly. Gemi integrates with the Gemini API to provide interactive AI capabilities directly in your terminal.

## Features

- Interactive chat with Gemini AI
- Text generation with prompts
- Streaming responses
- List and switch between available Gemini models
- Colorful and styled output
- Progress bars and spinners
- Command-line flags and arguments
- Subcommands

## Libraries Used

- [Cobra](https://github.com/spf13/cobra) - A Commander for modern Go CLI interactions
- [Google Generative AI Go SDK](https://github.com/google/generative-ai-go) - Official Go SDK for Google's Generative AI models
- [Fatih/Color](https://github.com/fatih/color) - Color package for Go
- [Briandowns/Spinner](https://github.com/briandowns/spinner) - Go package with 70+ spinners for terminal
- [Charmbracelet/Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions for terminal applications
- [Charmbracelet/Bubbles](https://github.com/charmbracelet/bubbles) - TUI components for Bubble Tea
- [Charmbracelet/BubbleTea](https://github.com/charmbracelet/bubbletea) - TUI framework for Go

## Installation

```bash
# Clone the repository
git clone https://github.com/vandi/gemi.git
cd gemi

# Install dependencies
go mod tidy

# Build the application
go build -o gemi
```

## Usage

### API Key Setup

Before using Gemi, you need to obtain a Gemini API key from [Google AI Studio](https://makersuite.google.com/app/apikey).

You can provide your API key in two ways:

1. Using the `--api-key` flag:
   ```bash
   ./gemi --api-key "YOUR_API_KEY"
   ```

2. Setting the `GEMINI_API_KEY` environment variable:
   ```bash
   export GEMINI_API_KEY="YOUR_API_KEY"
   ```

### Commands

```bash
# Run the main command to see available options
./gemi

# Start an interactive chat with Gemini AI
./gemi chat

# Generate text with a prompt
./gemi generate --prompt "Write a short poem about coding"

# Stream the response as it's generated
./gemi generate --prompt "Explain quantum computing" --stream

# Save the response to a file
./gemi generate --prompt "Write a Python script to sort a list" --output script.py

# List available Gemini models
./gemi models

# Use a specific model
./gemi chat --model gemini-1.5-flash-latest
./gemi generate --model gemini-1.5-flash-latest --prompt "Summarize this concept"

# Display version information
./gemi version
```

### Chat Commands

While in chat mode, you can use the following commands:

- `/help` - Show available commands
- `/models` or `/list-models` - List available models
- `/model MODEL_NAME` - Switch to a different model
- `/quit` - Exit the chat (or use Ctrl+C)

## License

MIT
# gemi
