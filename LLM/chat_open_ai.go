package LLM

import (
	"context"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type ChatOpenAI struct {
	Ctx          context.Context
	ModeName     string
	SystemPrompt string
	RagContext   string
	Tools        []mcp.Tool
	LLM          openai.Client
}

type LLMOption func(*ChatOpenAI)

func WithSystemPrompt(systemPrompt string) LLMOption {
	return func(c *ChatOpenAI) {
		c.SystemPrompt = systemPrompt
	}
}

func WithRagContext(ragContext string) LLMOption {
	return func(c *ChatOpenAI) {
		c.RagContext = ragContext
	}
}

func WithTools(tools []mcp.Tool) LLMOption {
	return func(c *ChatOpenAI) {
		c.Tools = tools
	}
}

func NewChatOpenAI(ctx context.Context, modeName string, opts ...LLMOption) *ChatOpenAI {
	if modeName == "" {
		panic("modeName cannot be empty")
	}

	var (
		apiKey  = os.Getenv("OPENAI_API")
		baseURL = os.Getenv("OPENAI_API_BASE_URL")
	)
	if apiKey == "" {
		panic("OPENAI_API environment variable is not set")
	}

	options := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}
	if baseURL != "" {
		options = append(options, option.WithBaseURL(baseURL))
	}

	cli := openai.NewClient(options...)

	llm := &ChatOpenAI{
		Ctx:      ctx,
		ModeName: modeName,
		LLM:      cli,
	}

	for _, opt := range opts {
		opt(llm)
	}
	return llm
}
