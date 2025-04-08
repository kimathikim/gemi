package ui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

var (
	// Colors
	PrimaryColor   = "#7D56F4"
	SecondaryColor = "#5F9EF3"
	AccentColor    = "#FF6B6B"
	SuccessColor   = "#10B981"
	WarningColor   = "#F59E0B"
	ErrorColor     = "#EF4444"
	TextColor      = "#FAFAFA"
	
	// Styles
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(TextColor)).
		Background(lipgloss.Color(PrimaryColor)).
		Padding(0, 3)
	
	SubtitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(SecondaryColor))
	
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(SecondaryColor)).
		Padding(1, 3)
	
	UserPromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(PrimaryColor)).
		Bold(true)
	
	AIResponseStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(SecondaryColor))
	
	// Color functions
	SuccessText = color.New(color.FgHiGreen).SprintFunc()
	InfoText    = color.New(color.FgHiCyan).SprintFunc()
	WarningText = color.New(color.FgHiYellow).SprintFunc()
	ErrorText   = color.New(color.FgHiRed).SprintFunc()
	
	// Prefixes
	SuccessPrefix = SuccessText("✓ ")
	InfoPrefix    = InfoText("ℹ ")
	WarningPrefix = WarningText("⚠ ")
	ErrorPrefix   = ErrorText("✗ ")
)

// RenderTitle renders a title with the title style
func RenderTitle(text string) string {
	return TitleStyle.Render(text)
}

// RenderBox renders text in a box
func RenderBox(text string) string {
	return BoxStyle.Render(text)
}

// RenderUserPrompt renders a user prompt
func RenderUserPrompt(text string) string {
	return UserPromptStyle.Render("You: ") + " " + text
}

// RenderAIResponse renders an AI response
func RenderAIResponse(text string) string {
	return AIResponseStyle.Render("Gemini: ") + " " + text
}

// RenderMarkdown renders text with Markdown styling for terminal display
func RenderMarkdown(text string) string {
	// Try to use Glamour first
	renderedText, err := RenderMarkdownWithGlamour(text)
	if err == nil {
		return renderedText
	}
	
	// Fall back to our custom renderer if Glamour fails
	// Apply styling to Markdown elements
	lines := strings.Split(text, "\n")
	
	// Track if we're inside a code block
	inCodeBlock := false
	
	// Process each line
	for i, line := range lines {
		// Handle code blocks
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			
			// Replace code block markers with terminal-friendly borders
			if inCodeBlock {
				// Start of code block
				codeLang := strings.TrimPrefix(line, "```")
				if codeLang != "" {
					lines[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("┌─── " + codeLang + " " + strings.Repeat("─", 50-len(codeLang)))
				} else {
					lines[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("┌" + strings.Repeat("─", 60))
				}
			} else {
				// End of code block
				lines[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("└" + strings.Repeat("─", 60))
			}
			continue
		}
		
		// If we're in a code block, style the code
		if inCodeBlock {
			lines[i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC")).Render("│ " + line)
			continue
		}
		
		// Style headings (outside code blocks)
		if strings.HasPrefix(line, "# ") {
			// H1 heading
			headingText := strings.TrimPrefix(line, "# ")
			lines[i] = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(PrimaryColor)).Render(headingText)
			// Add underline with ═ characters
			lines[i] += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color(PrimaryColor)).Render(strings.Repeat("═", len(headingText)))
		} else if strings.HasPrefix(line, "## ") {
			// H2 heading
			headingText := strings.TrimPrefix(line, "## ")
			lines[i] = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(SecondaryColor)).Render(headingText)
			// Add underline with ─ characters
			lines[i] += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color(SecondaryColor)).Render(strings.Repeat("─", len(headingText)))
		} else if strings.HasPrefix(line, "### ") {
			// H3 heading
			headingText := strings.TrimPrefix(line, "### ")
			lines[i] = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(AccentColor)).Render(headingText)
		}
		
		// Style lists
		if strings.HasPrefix(line, "* ") {
			listText := strings.TrimPrefix(line, "* ")
			lines[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(SecondaryColor)).Render("• ") + listText
		} else if strings.HasPrefix(line, "- ") {
			listText := strings.TrimPrefix(line, "- ")
			lines[i] = lipgloss.NewStyle().Foreground(lipgloss.Color(SecondaryColor)).Render("• ") + listText
		}
		
		// Style bold text with ** or __
		lines[i] = styleBoldText(lines[i])
		
		// Style inline code with backticks
		lines[i] = styleInlineCode(lines[i])
	}
	
	return strings.Join(lines, "\n")
}

// styleBoldText finds and styles bold text marked with ** or __
func styleBoldText(line string) string {
	// Handle **bold text**
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	line = boldRegex.ReplaceAllStringFunc(line, func(match string) string {
		// Extract the text between ** and **
		text := boldRegex.FindStringSubmatch(match)[1]
		return lipgloss.NewStyle().Bold(true).Render(text)
	})
	
	// Handle __bold text__
	boldRegex2 := regexp.MustCompile(`__([^_]+)__`)
	line = boldRegex2.ReplaceAllStringFunc(line, func(match string) string {
		// Extract the text between __ and __
		text := boldRegex2.FindStringSubmatch(match)[1]
		return lipgloss.NewStyle().Bold(true).Render(text)
	})
	
	return line
}

// styleInlineCode finds and styles inline code marked with backticks
func styleInlineCode(line string) string {
	codeRegex := regexp.MustCompile("`([^`]+)`")
	return codeRegex.ReplaceAllStringFunc(line, func(match string) string {
		// Extract the text between backticks
		text := codeRegex.FindStringSubmatch(match)[1]
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC")).Background(lipgloss.Color("#333333")).Padding(0, 1).Render(text)
	})
}
