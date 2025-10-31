package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	prompt := req.Template
	if req.Context != "" {
		prompt = fmt.Sprintf("%s\n\nAdditional Context:\n%s", req.Template, req.Context)
	}

	// Detect API type based on provider config
	isGemini := provider == "gemini"

	var reqBody []byte
	if isGemini {
		apiEndpoint = fmt.Sprintf("%s/models/%s:generateContent", apiEndpoint, model)

		// Create Gemini request
		geminiReq := GeminiRequest{
			Contents: []GeminiContent{
				{
					Parts: []GeminiPart{
						{
							Text: fmt.Sprintf("You are a helpful assistant that generates mock API response data. Generate realistic and valid JSON data based on this request:\n\n%s", prompt),
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
					Content: "You are a helpful assistant that generates mock API response data based on templates. Generate realistic and valid JSON data.",
				},
				{
					Role:    "user",
					Content: prompt,
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

	return &GenerateResponse{
		Content:   content,
		Model:     model,
		TokenUsed: tokenUsed,
	}, nil
}
