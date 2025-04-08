package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
)

// Default style to use when rendering markdown
const DefaultStyle = "dark"

// RenderMarkdownWithGlamour renders markdown text using Glamour
func RenderMarkdownWithGlamour(markdown string) (string, error) {
	// Check if a style is set in the environment
	style := os.Getenv("GLAMOUR_STYLE")
	if style == "" {
		style = DefaultStyle
	}

	// Create a renderer with the specified style
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),   // Automatically detect terminal style
		glamour.WithWordWrap(100), // Wrap at 100 characters
	)

	if err != nil {
		return "", fmt.Errorf("failed to create markdown renderer: %v", err)
	}

	// Render the markdown
	rendered, err := r.Render(markdown)
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %v", err)
	}

	return rendered, nil
}

// MustRenderMarkdown renders markdown text and panics on error
func MustRenderMarkdown(markdown string) string {
	rendered, err := RenderMarkdownWithGlamour(markdown)
	if err != nil {
		panic(err)
	}
	return rendered
}
