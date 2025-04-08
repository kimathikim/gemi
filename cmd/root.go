package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	apiKey  string
	rootCmd = &cobra.Command{
		Use:   "gemi",
		Short: "Gemi is a beautiful CLI tool powered by Gemini AI",
		Long: `A beautiful CLI tool built with Cobra and enhanced with various libraries
to make it visually appealing and user-friendly. It uses the Gemini API
to provide interactive AI capabilities.`,
		Run: func(cmd *cobra.Command, args []string) {
			showWelcome()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Gemini API key (or set GEMINI_API_KEY env var)")

	// Add commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(modelsCmd)
}

func showWelcome() {
	// Display a welcome message with color
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 3).
		Render("Welcome to Gemi CLI")

	fmt.Println()
	fmt.Println(title)
	fmt.Println()

	// Show a spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Loading resources..."
	s.Color("cyan")
	s.Start()
	time.Sleep(1 * time.Second)
	s.Stop()

	// Show a progress bar
	fmt.Println()
	prog := progress.New(progress.WithDefaultGradient())
	fmt.Println("Initializing components:")

	for i := 0; i <= 100; i += 20 {
		prog.SetPercent(float64(i) / 100)
		fmt.Printf("\r%s", prog.View())
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println()
	fmt.Println()

	// Display available commands
	success := color.New(color.FgGreen).SprintFunc()
	info := color.New(color.FgCyan).SprintFunc()

	fmt.Println(success("✓ ") + "Gemi is ready to use!")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println(info("  gemi chat") + "      - Start an interactive chat with Gemini AI")
	fmt.Println(info("  gemi generate") + "  - Generate text with Gemini AI")
	fmt.Println(info("  gemi models") + "    - List available Gemini models")
	fmt.Println(info("  gemi version") + "   - Display version information")
	fmt.Println()

	// Check for API key
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			warning := color.New(color.FgYellow).SprintFunc()
			fmt.Println(warning("⚠ ") + "No API key found. Please set your Gemini API key using:")
			fmt.Println("  - The --api-key flag")
			fmt.Println("  - Or the GEMINI_API_KEY environment variable")
			fmt.Println()
		}
	}
}

func getApiKey() (string, error) {
	key := apiKey
	if key == "" {
		key = os.Getenv("GEMINI_API_KEY")
		if key == "" {
			return "", fmt.Errorf("no API key provided. Use --api-key flag or set GEMINI_API_KEY environment variable")
		}
	}
	return key, nil
}
