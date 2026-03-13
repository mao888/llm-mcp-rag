package MCP

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPClient struct {
	Ctx    context.Context
	Client *client.Client
	Cmd    string
	Tools  []mcp.Tool
	Args   []string
	Env    []string
}

func NewMCPClient(ctx context.Context, cmd string, env, args []string) *MCPClient {
	stdioTransport := transport.NewStdio(cmd, env, args...)
	client := client.NewClient(stdioTransport)
	return &MCPClient{
		Ctx:    ctx,
		Client: client,
		Cmd:    cmd,
		Args:   args,
		Env:    env,
	}
}

func (m *MCPClient) Start() error {
	var err error
	err = m.Client.Start(m.Ctx)
	if err != nil {
		return err
	}

	mcpInitReq := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "mcp-client-go",
				Version: "0.1.0",
			},
		},
	}
	if _, err = m.Client.Initialize(m.Ctx, mcpInitReq); err != nil {
		fmt.Println("mcp init error:", err)
		return err
	}
	return nil
}

func (m *MCPClient) SetTools() error {
	var (
		err   error
		tools *mcp.ListToolsResult
	)

	toolsReq := mcp.ListToolsRequest{}
	if tools, err = m.Client.ListTools(m.Ctx, toolsReq); err != nil {
		return err
	}
	m.Tools = tools.Tools
	return nil
}

func (m *MCPClient) Close() error {
	return m.Client.Close()
}

func (m *MCPClient) CallTool(name string, args any) (string, error) {
	var arguments map[string]any
	switch v := args.(type) {
	case map[string]any:
		arguments = v
	case string:
		err := json.Unmarshal([]byte(v), &arguments)
		if err != nil {
			return "", err
		}
	default:
		arguments = make(map[string]any)
	}

	res, err := m.Client.CallTool(m.Ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      name,
			Arguments: arguments,
		},
	})
	if err != nil {
		return "", err
	}
	return mcp.GetTextFromContent(res.Content), nil
}

func (m *MCPClient) GetTools() []mcp.Tool {
	return m.Tools
}
