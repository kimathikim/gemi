package gemini

import (
	"context"
	"fmt"
	"io"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Client wraps the Gemini API client
type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
	ctx    context.Context
}

// NewClient creates a new Gemini client
func NewClient(apiKey string, modelName string) (*Client, error) {
	if modelName == "" {
		modelName = "gemini-1.5-pro-latest"
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}

	model := client.GenerativeModel(modelName)
	model.Temperature = genai.Ptr[float32](0.7)

	return &Client{
		client: client,
		model:  model,
		ctx:    ctx,
	}, nil
}

// Close closes the client
func (c *Client) Close() error {
	return c.client.Close()
}

// GenerateText generates text from a prompt
func (c *Client) GenerateText(ctx context.Context, prompt string) (string, error) {
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %v", err)
	}

	return responseToString(resp), nil
}

// GenerateTextStream generates text from a prompt and streams the response
func (c *Client) GenerateTextStream(ctx context.Context, prompt string, writer io.Writer) error {
	iter := c.model.GenerateContentStream(ctx, genai.Text(prompt))

	for {
		resp, err := iter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to get next response: %v", err)
		}

		text := responseToString(resp)
		if _, err := fmt.Fprint(writer, text); err != nil {
			return fmt.Errorf("failed to write response: %v", err)
		}
	}

	return nil
}

// StartChat starts a new chat session
func (c *Client) StartChat() *genai.ChatSession {
	return c.model.StartChat()
}

// ListModels lists all available models
func (c *Client) ListModels() ([]*genai.ModelInfo, error) {
	iter := c.client.ListModels(c.ctx)
	var models []*genai.ModelInfo

	for {
		model, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list models: %v", err)
		}
		models = append(models, model)
	}

	return models, nil
}

// SwitchModel switches to a different model
func (c *Client) SwitchModel(modelName string) error {
	if modelName == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	c.model = c.client.GenerativeModel(modelName)
	c.model.Temperature = genai.Ptr[float32](0.7)
	return nil
}

// responseToString extracts text from a GenerateContentResponse
func responseToString(resp *genai.GenerateContentResponse) string {
	var result string
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				if text, ok := part.(genai.Text); ok {
					result += string(text)
				}
			}
		}
	}
	return result
}
