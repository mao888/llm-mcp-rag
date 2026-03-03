# llm-mcp-rag
基于Golang的LLM-MCP-RAG的小demo

## LLM
- 使用了函数式选项模式封装 LLM 客户端，
- 支持 SystemPrompt、RAG 上下文、Tools 等可选扩展参数，
- 同时内部使用 openai-go SDK 的 Option 模式初始化 client，
- 这样可以保持构造函数稳定，同时具备良好的扩展性。