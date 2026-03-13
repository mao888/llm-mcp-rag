package Agent

import (
	"context"
	"fmt"

	"github.com/mao888/llm-mcp-rag/LLM"
	"github.com/mao888/llm-mcp-rag/MCP"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go/v3"
)

type Agent struct {
	MCPClient  []*MCP.MCPClient
	LLM        *LLM.ChatOpenAI
	Mode       string
	Ctx        context.Context
	RagContext string
}

func NewAgent(ctx context.Context, model string, mcpClients []*MCP.MCPClient) *Agent {
	// 1. Initialize MCP client and get all tools
	tools := make([]mcp.Tool, 0)
	for _, mcpClient := range mcpClients {
		err := mcpClient.Start()
		if err != nil {
			fmt.Println("Failed to start MCP client:", err)
			continue
		}
		err = mcpClient.SetTools()
		if err != nil {
			fmt.Println("Failed to set tools:", err)
			continue
		}
		tools = append(tools, mcpClient.GetTools()...)
	}
	// 2. Initialize LLM with all tools
	llm := LLM.NewChatOpenAI(ctx, model, LLM.WithTools(tools))
	// 3. Return Agent with LLM
	return &Agent{
		MCPClient: mcpClients,
		LLM:       llm,
		Mode:      model,
		Ctx:       ctx,
	}
}

func (a *Agent) StartAction(prompt string) string {
	if a.LLM == nil {
		fmt.Println("LLM is nil")
		return ""
	}
	response, toolCalls := a.LLM.Chat(prompt)
	for len(toolCalls) > 0 {
		// 找到这个tool是属于哪个MCP client
		for _, toolCall := range toolCalls { // 遍历所有的tool call, 这个tool call是llm返回的
			for _, mcpClient := range a.MCPClient { // 遍历所有的MCP client, 找到这个tool是属于哪个MCP client
				tools := mcpClient.GetTools() // 获取这个MCP client的所有tools
				for _, tool := range tools {  // 遍历这个MCP client的所有tools, 找到这个tool call对应的tool
					if tool.Name == toolCall.Function.Name {
						fmt.Println("Tool use:", toolCall.Function.Name)
						// 执行工具调用
						toolResult, err := mcpClient.CallTool(toolCall.Function.Name, toolCall.Function.Arguments)
						if err != nil {
							fmt.Println("Failed to execute tool:", err)
							continue
						}
						a.LLM.Messages = append(a.LLM.Messages, openai.ToolMessage(toolResult, toolCall.ID))
					}
				}
			}
		}
		response, toolCalls = a.LLM.Chat("") // 继续调用llm, 直到没有tool call为止
	}
	a.Close()
	return response
}

func (a *Agent) Close() {
	for _, mcpClient := range a.MCPClient {
		mcpClient.Close()
	}
}
