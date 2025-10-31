package ai

// GenerateRequest represents the request payload for AI generation
type GenerateRequest struct {
	Template string `json:"template" binding:"required"` // Template/prompt for generation
	Context  string `json:"context"`                     // Additional context for generation
}

// GenerateResponse represents the response from AI generation
type GenerateResponse struct {
	Content   string `json:"content"`    // Generated content
	Model     string `json:"model"`      // Model used for generation
	TokenUsed int    `json:"token_used"` // Tokens used for generation
}

// OpenAIRequest represents the request to OpenAI API
type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

// OpenAIMessage represents a message in OpenAI API
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response from OpenAI API
type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// GeminiRequest represents the request to Gemini API
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents content in Gemini API
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part in Gemini content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response from Gemini API
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}
