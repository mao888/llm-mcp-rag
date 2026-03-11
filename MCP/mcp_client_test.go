package MCP

import (
	"context"
	"fmt"
	"testing"
)

func TestMCPClient(t *testing.T) {
	var err error
	ctx := context.Background()
	client := NewMCPClient(ctx, "uvx", nil, []string{"mcp-server-fetch"})
	if err = client.Start(); err != nil {
		t.Fatalf("Failed to start MCP client: %v", err)
	}

	err = client.SetTools()
	if err != nil {
		t.Fatalf("Failed to set tools: %v", err)
	}

	tools := client.GetTools()
	if len(tools) == 0 {
		t.Fatalf("No tools were get from the MCP client")
	}
	fmt.Println("Tools received from MCP server:", tools)
}
