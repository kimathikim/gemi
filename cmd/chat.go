package cmd

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"github.com/vandi/gemi/internal/gemini"
	"github.com/vandi/gemi/internal/ui"
)

var (
	modelName  string
	listModels bool

	chatCmd = &cobra.Command{
		Use:   "chat",
		Short: "Start an interactive chat with Gemini AI",
		Long:  `Start an interactive chat session with Gemini AI in your terminal.`,
		Run: func(cmd *cobra.Command, args []string) {
			// If --list-models flag is provided, list models and exit
			if listModels {
				modelsCmd.Run(cmd, args)
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

			// Start the chat UI
			p := tea.NewProgram(initialChatModel(client, client.StartChat()))
			if _, err := p.Run(); err != nil {
				fmt.Println(ui.ErrorPrefix + "Error running chat: " + err.Error())
			}
		},
	}
)

func init() {
	chatCmd.Flags().StringVar(&modelName, "model", "gemini-1.5-pro-latest", "Gemini model to use")
	chatCmd.Flags().BoolVar(&listModels, "list-models", false, "List available Gemini models")
}

// Chat UI model
type chatModel struct {
	client       *gemini.Client
	chatSession  *genai.ChatSession
	messages     []message
	textInput    textinput.Model
	err          error
	width        int
	height       int
	currentModel string
}

type message struct {
	content string
	isUser  bool
}

func initialChatModel(client *gemini.Client, chatSession *genai.ChatSession) chatModel {
	ti := textinput.New()
	ti.Placeholder = "Type your message and press Enter (Ctrl+C to quit)"
	ti.Focus()
	ti.Width = 80

	return chatModel{
		client:       client,
		chatSession:  chatSession,
		textInput:    ti,
		messages:     []message{},
		width:        80,
		height:       24,
		currentModel: modelName,
	}
}

func (m chatModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textInput.Value() == "" {
				return m, nil
			}

			userInput := m.textInput.Value()
			m.messages = append(m.messages, message{content: userInput, isUser: true})
			m.textInput.Reset()

			// Check for special commands
			if strings.HasPrefix(userInput, "/model ") {
				// Command to change the model
				newModel := strings.TrimPrefix(userInput, "/model ")
				return m, func() tea.Msg {
					err := m.client.SwitchModel(newModel)
					if err != nil {
						return errorMsg{err}
					}

					// Create a new chat session with the new model
					m.chatSession = m.client.StartChat()
					m.currentModel = newModel

					return responseMsg{content: "Switched to model: " + newModel}
				}
			} else if userInput == "/models" || userInput == "/list-models" {
				// Command to list available models in Markdown format
				return m, func() tea.Msg {
					models, err := m.client.ListModels()
					if err != nil {
						return errorMsg{err}
					}

					var sb strings.Builder
					sb.WriteString("# Available Models\n\n")

					// Group models by base model ID for cleaner output
					modelsByBase := make(map[string][]string)
					for _, model := range models {
						parts := strings.Split(model.Name, "/")
						modelName := parts[len(parts)-1]
						modelsByBase[model.BaseModelID] = append(modelsByBase[model.BaseModelID], modelName)
					}

					// Sort base model IDs for consistent output
					baseModelIDs := make([]string, 0, len(modelsByBase))
					for baseID := range modelsByBase {
						baseModelIDs = append(baseModelIDs, baseID)
					}
					sort.Strings(baseModelIDs)

					for _, baseID := range baseModelIDs {
						sb.WriteString("## " + baseID + "\n\n")

						// Sort model names for consistent output
						modelNames := modelsByBase[baseID]
						sort.Strings(modelNames)

						for _, name := range modelNames {
							sb.WriteString("* **" + name + "**\n")
						}
						sb.WriteString("\n")
					}

					sb.WriteString("**Current model:** " + m.currentModel + "\n\n")
					sb.WriteString("To change models, type: `/model MODEL_NAME`")

					return responseMsg{content: sb.String()}
				}
			} else if userInput == "/help" {
				// Command to show help in Markdown format
				return m, func() tea.Msg {
					help := "# Available Commands\n\n" +
						"* **`/models`** or **`/list-models`** - List available models\n" +
						"* **`/model MODEL_NAME`** - Switch to a different model\n" +
						"* **`/help`** - Show this help message\n" +
						"* **`/quit`** or **`Ctrl+C`** - Exit the chat"
					return responseMsg{content: help}
				}
			} else if userInput == "/quit" {
				return m, tea.Quit
			} else {
				// Regular message to Gemini
				return m, func() tea.Msg {
					ctx := context.Background()
					resp, err := m.chatSession.SendMessage(ctx, genai.Text(userInput))
					if err != nil {
						return errorMsg{err}
					}

					var aiResponse string
					for _, candidate := range resp.Candidates {
						if candidate.Content != nil {
							for _, part := range candidate.Content.Parts {
								if text, ok := part.(genai.Text); ok {
									aiResponse += string(text)
								}
							}
						}
					}

					return responseMsg{content: aiResponse}
				}
			}
		}

	case responseMsg:
		m.messages = append(m.messages, message{content: msg.content, isUser: false})

	case errorMsg:
		m.err = msg.err

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m chatModel) View() string {
	var s strings.Builder

	// Title with current model
	title := ui.RenderTitle(" Gemini Chat - " + m.currentModel + " ")
	s.WriteString(title + "\n\n")

	// Messages
	if len(m.messages) == 0 {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("Start chatting with Gemini AI...") + "\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("Type /help to see available commands") + "\n\n")
	} else {
		// Calculate available height for messages
		availableHeight := m.height - 7 // Adjust based on other UI elements

		// If we have more messages than can fit, show only the most recent ones
		startIdx := 0
		if len(m.messages) > availableHeight/2 {
			startIdx = len(m.messages) - availableHeight/2
		}

		for i := startIdx; i < len(m.messages); i++ {
			msg := m.messages[i]
			if msg.isUser {
				s.WriteString(ui.RenderUserPrompt(msg.content) + "\n\n")
			} else {
				// Apply Markdown formatting to AI responses using Glamour
				formattedContent, err := ui.RenderMarkdownWithGlamour(msg.content)
				if err != nil {
					s.WriteString(ui.ErrorPrefix + "Failed to render markdown: " + err.Error() + "\n\n")
					s.WriteString(ui.AIResponseStyle.Render("Gemini: ") + "\n\n" + msg.content + "\n\n")
				} else {
					s.WriteString(ui.AIResponseStyle.Render("Gemini: ") + "\n\n" + formattedContent + "\n")
				}
			}
		}
	}

	// Error message
	if m.err != nil {
		s.WriteString(ui.ErrorPrefix + m.err.Error() + "\n\n")
	}

	// Input field
	s.WriteString(m.textInput.View() + "\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("Press Ctrl+C to quit") + "\n")

	return s.String()
}

// Message types for the tea.Program
type responseMsg struct {
	content string
}

type errorMsg struct {
	err error
}
