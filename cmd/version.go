package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	Version = "1.0.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  `Display the current version of the Gemi CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a styled version display
		versionStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#5F9EF3"))
		
		boxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5F9EF3")).
			Padding(1, 3).
			MarginTop(1).
			MarginBottom(1)
		
		versionInfo := fmt.Sprintf("Gemi CLI version %s", Version)
		fmt.Println(boxStyle.Render(versionStyle.Render(versionInfo)))
	},
}
