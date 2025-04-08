package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	// "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/vandi/gemi/internal/gemini"
	"github.com/vandi/gemi/internal/ui"
)

var (
	modelsCmd = &cobra.Command{
		Use:   "models",
		Short: "List available Gemini models",
		Long:  `List all available Gemini models that can be used with the chat and generate commands.`,
		Run: func(cmd *cobra.Command, args []string) {
			apiKey, err := getApiKey()
			if err != nil {
				fmt.Println(ui.ErrorPrefix + err.Error())
				return
			}

			// Create a client with any model (we'll just use it to list models)
			client, err := gemini.NewClient(apiKey, "gemini-1.5-pro-latest")
			if err != nil {
				fmt.Println(ui.ErrorPrefix + "Failed to initialize Gemini client: " + err.Error())
				return
			}
			defer client.Close()

			// Show a spinner while fetching models
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Prefix = "Fetching available models "
			s.Color("cyan")
			s.Start()

			// Get the list of models
			models, err := client.ListModels()
			s.Stop()

			if err != nil {
				fmt.Println(ui.ErrorPrefix + "Failed to list models: " + err.Error())
				return
			}

			// Sort models by name
			sort.Slice(models, func(i, j int) bool {
				return models[i].Name < models[j].Name
			})

			// Display the models in Markdown-friendly format
			title := ui.RenderTitle(" Available Gemini Models ")
			fmt.Println("\n" + title + "\n")

			// Group models by base model ID
			modelsByBase := make(map[string][]*struct {
				Name    string
				Version string
			})

			for _, model := range models {
				// Extract the model name from the full resource name
				// Format is typically "models/{model_name}"
				parts := strings.Split(model.Name, "/")
				modelName := parts[len(parts)-1]

				baseInfo := &struct {
					Name    string
					Version string
				}{
					Name:    modelName,
					Version: model.Version,
				}

				modelsByBase[model.BaseModelID] = append(modelsByBase[model.BaseModelID], baseInfo)
			}

			// Sort base model IDs for consistent output
			baseModelIDs := make([]string, 0, len(modelsByBase))
			for baseID := range modelsByBase {
				baseModelIDs = append(baseModelIDs, baseID)
			}
			sort.Strings(baseModelIDs)

			// Build a Markdown string
			var markdownOutput strings.Builder

			// Display models by base model ID in a Markdown-friendly format
			for i, baseID := range baseModelIDs {
				// Add a separator between base models except for the first one
				if i > 0 {
					markdownOutput.WriteString("---\n\n")
				}

				// Print header as a Markdown heading
				markdownOutput.WriteString("# " + baseID + "\n\n")

				// Sort models for consistent output
				models := modelsByBase[baseID]
				sort.Slice(models, func(i, j int) bool {
					return models[i].Name < models[j].Name
				})

				// Print models as a list with proper indentation
				for _, modelInfo := range models {
					markdownOutput.WriteString("* **" + modelInfo.Name + "** (version: " + modelInfo.Version + ")\n")
				}
				markdownOutput.WriteString("\n")
			}

			// Add usage instructions to the Markdown
			markdownOutput.WriteString("# Usage Instructions\n\n")
			markdownOutput.WriteString("To use a specific model:\n\n")
			markdownOutput.WriteString("```bash\n")
			markdownOutput.WriteString("gemi chat --model MODEL_NAME\n")
			markdownOutput.WriteString("gemi generate --model MODEL_NAME --prompt \"Your prompt\"\n")
			markdownOutput.WriteString("```\n\n")
			markdownOutput.WriteString("In chat mode, you can also switch models using:\n\n")
			markdownOutput.WriteString("```\n")
			markdownOutput.WriteString("/model MODEL_NAME\n")
			markdownOutput.WriteString("```\n")

			// Render the Markdown using Glamour
			renderedMarkdown, err := ui.RenderMarkdownWithGlamour(markdownOutput.String())
			if err != nil {
				fmt.Println(ui.ErrorPrefix + "Failed to render markdown: " + err.Error())
				return
			}

			fmt.Println(renderedMarkdown)
		},
	}
)

func init() {
	rootCmd.AddCommand(modelsCmd)
}
