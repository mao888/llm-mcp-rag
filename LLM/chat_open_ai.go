package LLM

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

type ChatOpenAI struct {
	Ctx          context.Context
	ModeName     string
	SystemPrompt string
	RagContext   string
	Tools        []mcp.Tool
	Messages     []openai.ChatCompletionMessageParamUnion
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
	if llm.SystemPrompt != "" {
		llm.Messages = append(llm.Messages, openai.SystemMessage(llm.SystemPrompt))
	}
	if llm.RagContext != "" {
		llm.Messages = append(llm.Messages, openai.UserMessage(llm.RagContext))
	}
	fmt.Printf("Successfully initialized ChatOpenAI llm: %s\n", llm.ModeName)
	return llm
}

func (c *ChatOpenAI) Chat(prompt string) (result string, toolCall []openai.ToolCallUnion) {
	if prompt != "" {
		c.Messages = append(c.Messages, openai.UserMessage(prompt))
	}
	toolParams := MCPTool2OpenAITool(c.Tools)
	stream := c.LLM.Chat.Completions.NewStreaming(c.Ctx, openai.ChatCompletionNewParams{
		Messages: c.Messages,
		Model:    c.ModeName,
		Seed:     openai.Int(0),
		Tools:    toolParams,
	})

	result = ""
	finished := false
	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if content, ok := acc.JustFinishedContent(); ok {
			finished = true
			result = content
		}

		if tool, ok := acc.JustFinishedToolCall(); ok {
			fmt.Println("Tool called:", tool)
			toolCall = append(toolCall, openai.ToolCallUnion{
				ID: tool.ID,
				Function: openai.FunctionToolCallFunction{
					Arguments: tool.Arguments,
					Name:      tool.Name,
				},
			})
		}

		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta.Content
			if !finished {
				result += delta
			}
		}
	}
	if stream.Err() != nil {
		panic(stream.Err())
	}
	return result, toolCall
}

func MCPTool2OpenAITool(tool []mcp.Tool) []openai.ChatCompletionToolUnionParam {
	openaiTools := make([]openai.ChatCompletionToolUnionParam, 0, len(tool))
	for _, t := range tool {
		openaiTools = append(openaiTools, openai.ChatCompletionToolUnionParam{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: shared.FunctionDefinitionParam{
					Name:        t.Name,
					Description: openai.String(t.Description),
					Parameters: openai.FunctionParameters{
						"type":       t.InputSchema.Type,
						"properties": t.InputSchema.Properties,
						"required":   t.InputSchema.Required,
					},
				},
			},
		})
	}
	return openaiTools
}
