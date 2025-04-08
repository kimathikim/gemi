package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/vandi/gemi/internal/gemini"
	"github.com/vandi/gemi/internal/ui"
)

var (
	prompt        string
	outputFile    string
	stream        bool
	listModelsGen bool

	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate text with Gemini AI",
		Long:  `Generate text using Gemini AI based on a prompt.`,
		Run: func(cmd *cobra.Command, args []string) {
			// If --list-models flag is provided, list models and exit
			if listModelsGen {
				modelsCmd.Run(cmd, args)
				return
			}

			if prompt == "" {
				fmt.Println(ui.ErrorPrefix + "Prompt is required. Use --prompt or -p flag.")
				return
			}

			apiKey, err := getApiKey()
			if err != nil {
				fmt.Println(ui.ErrorPrefix + err.Error())
				return
			}

			client, err := gemini.NewClient(apiKey, modelName)
			if err != nil {
				fmt.Println(ui.ErrorPrefix + "Failed to initialize Gemini client: " + err.Error())
				return
			}
			defer client.Close()

			ctx := context.Background()

			// Show prompt with Markdown formatting using Glamour
			promptMd := "# Prompt\n\n```\n" + prompt + "\n```\n\n# Response\n"
			formattedPrompt, err := ui.RenderMarkdownWithGlamour(promptMd)
			if err != nil {
				fmt.Println(ui.ErrorPrefix + "Failed to render markdown: " + err.Error())
				fmt.Println("Prompt: " + prompt + "\n\nResponse:")
			} else {
				fmt.Println(formattedPrompt)
			}

			// Create a spinner
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Prefix = "Generating "
			s.Color("cyan")

			var result string

			if stream {
				// Create a custom writer that applies Markdown formatting using Glamour
				markdownWriter := &markdownStreamWriter{}

				if err := client.GenerateTextStream(ctx, prompt, markdownWriter); err != nil {
					fmt.Println("\n" + ui.ErrorPrefix + "Error generating response: " + err.Error())
					return
				}
				fmt.Println()
			} else {
				// Generate the response
				s.Start()
				result, err = client.GenerateText(ctx, prompt)
				s.Stop()

				if err != nil {
					fmt.Println(ui.ErrorPrefix + "Error generating response: " + err.Error())
					return
				}

				// Print the response with Markdown formatting using Glamour
				formattedResult, err := ui.RenderMarkdownWithGlamour(result)
				if err != nil {
					fmt.Println(ui.ErrorPrefix + "Failed to render markdown: " + err.Error())
					fmt.Println(result)
				} else {
					fmt.Println(formattedResult)
				}
			}

			// Save to file if requested
			if outputFile != "" && !stream {
				if err := os.WriteFile(outputFile, []byte(result), 0644); err != nil {
					fmt.Println(ui.ErrorPrefix + "Error saving to file: " + err.Error())
					return
				}
				fmt.Println(ui.SuccessPrefix + "Response saved to " + outputFile)
			}
		},
	}
)

// markdownStreamWriter is a custom io.Writer that applies Markdown formatting using Glamour to streamed content
type markdownStreamWriter struct {
	buffer strings.Builder
}

func (w *markdownStreamWriter) Write(p []byte) (n int, err error) {
	// Convert bytes to string
	text := string(p)

	// Apply Markdown formatting using Glamour
	// We accumulate the text first to handle multi-line Markdown elements
	w.buffer.WriteString(text)

	// Format and print the accumulated text using Glamour
	formatted, renderErr := ui.RenderMarkdownWithGlamour(w.buffer.String())

	// Clear the terminal line and reprint the entire formatted buffer
	// This ensures proper rendering of multi-line elements
	fmt.Print("\r\033[K") // Clear the current line

	if renderErr != nil {
		// If rendering fails, just print the plain text
		fmt.Print(w.buffer.String())
	} else {
		fmt.Print(formatted)
	}

	// Return the number of bytes written
	return len(p), nil
}

func init() {
	generateCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "The prompt to send to Gemini AI")
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Save the response to a file")
	generateCmd.Flags().BoolVarP(&stream, "stream", "s", false, "Stream the response as it's generated")
	generateCmd.Flags().StringVar(&modelName, "model", "gemini-1.5-pro-latest", "Gemini model to use")
	generateCmd.Flags().BoolVar(&listModelsGen, "list-models", false, "List available Gemini models")
}
