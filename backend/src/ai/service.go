package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	systemConfig "beo-echo/backend/src/systemConfigs"

	"github.com/rs/zerolog/log"
)

// AIService implements AI generation operations
type AIService struct {
	client *http.Client
}

// NewAIService creates a new AI service
func NewAIService() *AIService {
	return &AIService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Generate generates content based on the provided template and context
func (s *AIService) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	// Get AI configuration from system config
	apiKey, err := systemConfig.GetSystemConfigWithType[string](systemConfig.AI_API_KEY)
	if err != nil || apiKey == "" {
		return nil, fmt.Errorf("AI API key not configured")
	}

	apiEndpoint, err := systemConfig.GetSystemConfigWithType[string](systemConfig.AI_API_ENDPOINT)
	if err != nil || apiEndpoint == "" {
		return nil, fmt.Errorf("AI API endpoint not configured")
	}

	model, err := systemConfig.GetSystemConfigWithType[string](systemConfig.AI_MODEL)
	if err != nil || model == "" {
		model = "gemini-pro" // Default model
	}

	// Get AI provider type
	provider, err := systemConfig.GetSystemConfigWithType[string](systemConfig.AI_PROVIDER)
	if err != nil || provider == "" {
		provider = "gemini" // Default provider
	}

	// Build the prompt
	userMessage := req.Message
	if req.Context != "" {
		userMessage = fmt.Sprintf("%s\n\nCurrent editor content:\n```\n%s\n```", req.Message, req.Context)
	}

	contentType := req.ContentType
	if contentType == "" {
		contentType = "application/json"
	}

	// Detect API type based on provider config
	isGemini := provider == "gemini"

	responseFormatPrompt := "return ONLY the raw response body in the requested format data in http response" + contentType + " with no markdown, code blocks, explanations, or wrapping."

	// Build system prompt
	systemPrompt := `You are an AI assistant for a mock API service editor. Your role is to:

1. **Generate mock data** when users request it - ` + responseFormatPrompt + `
2. **Answer questions** about their API responses or mock data - Respond conversationally in user's language
3. **Provide help** and suggestions when asked

Examples:
- User: "generate user data" → Return: {"id": 1, "name": "John Doe", "email": "john@example.com"}
- User: "good morning" → Return: Good morning! How can I help you today?
- User: "what is this?" → Explain the current content conversationally
- User: "10 cities in indonesia" → Return: ["Jakarta", "Surabaya", "Bandung", "Medan", "Semarang", "Palembang", "Makassar", "Denpasar", "Yogyakarta", "Balikpapan"]

Remember: When generating data, return ONLY the JSON. No markdown, no code blocks, no explanations.`

	var reqBody []byte
	if isGemini {
		apiEndpoint = fmt.Sprintf("%s/models/%s:generateContent", apiEndpoint, model)

		// Create Gemini request
		geminiReq := GeminiRequest{
			Contents: []GeminiContent{
				{
					Parts: []GeminiPart{
						{
							Text: fmt.Sprintf("%s\n\nUser: %s", systemPrompt, userMessage),
						},
					},
				},
			},
		}
		reqBody, err = json.Marshal(geminiReq)
	} else {
		// Create OpenAI request
		openAIReq := OpenAIRequest{
			Model: model,
			Messages: []OpenAIMessage{
				{
					Role:    "system",
					Content: systemPrompt,
				},
				{
					Role:    "user",
					Content: userMessage,
				},
			},
		}
		reqBody, err = json.Marshal(openAIReq)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	if isGemini {
		// Gemini uses X-goog-api-key header
		httpReq.Header.Set("X-goog-api-key", apiKey)
	} else {
		// OpenAI uses Bearer token
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	}

	// Execute request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status_code", resp.StatusCode).
			Str("response_body", string(body)).
			Msg("AI API returned error")
		return nil, fmt.Errorf("AI API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response based on API type
	var content string
	var tokenUsed int
	var data string

	if isGemini {
		var geminiResp GeminiResponse
		if err := json.Unmarshal(body, &geminiResp); err != nil {
			return nil, fmt.Errorf("failed to parse Gemini response: %w", err)
		}

		if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
			return nil, fmt.Errorf("no response from Gemini AI")
		}

		content = geminiResp.Candidates[0].Content.Parts[0].Text
		tokenUsed = geminiResp.UsageMetadata.TotalTokenCount
	} else {
		var openAIResp OpenAIResponse
		if err := json.Unmarshal(body, &openAIResp); err != nil {
			return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
		}

		if len(openAIResp.Choices) == 0 {
			return nil, fmt.Errorf("no response from OpenAI")
		}

		content = openAIResp.Choices[0].Message.Content
		tokenUsed = openAIResp.Usage.TotalTokens
	}

	// Extract data and check if can apply
	extractedData, canApply := extractAndValidate(content)
	if canApply {
		data = extractedData
	}

	return &GenerateResponse{
		Content:   content,
		Model:     model,
		TokenUsed: tokenUsed,
		CanApply:  canApply,
		Data:      data,
	}, nil
}

// extractAndValidate checks if content is valid data format or contains code blocks
// Returns the extracted data and whether it can be applied
func extractAndValidate(content string) (string, bool) {
	// First, try direct validation for various formats
	if isValidDataFormat(content) {
		return content, true
	}

	// If not valid, check if content contains code block markers
	// Extract content from ```...```
	if strings.Contains(content, "```") {
		extracted := extractFromCodeBlock(content)
		if extracted != "" {
			return extracted, true
		}
	}

	return "", false
}

// isValidDataFormat checks if content is valid JSON, XML, or CSV
func isValidDataFormat(content string) bool {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return false
	}

	// Check for JSON
	var js json.RawMessage
	if err := json.Unmarshal([]byte(trimmed), &js); err == nil {
		return true
	}

	return false
}

// extractFromCodeBlock extracts content from markdown code block
// Removes ``` markers and language identifier (e.g., ```json)
func extractFromCodeBlock(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Start of code block
		if strings.HasPrefix(trimmed, "```") {
			if !inBlock {
				inBlock = true
				continue // Skip the opening ``` line
			} else {
				// End of code block
				break
			}
		}

		// Collect lines inside the block
		if inBlock {
			result = append(result, line)
		}
	}

	return strings.TrimSpace(strings.Join(result, "\n"))
}
