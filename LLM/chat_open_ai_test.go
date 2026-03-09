package LLM

import (
	"context"
	"fmt"
	"testing"

	"github.com/openai/openai-go/v3"
)

func TestChatOpenAI_Chat(t *testing.T) {
	ctx := context.Background()
	modeName := openai.ChatModelGPT3_5Turbo
	llm := NewChatOpenAI(ctx, modeName)
	prompt := "hello, what is your name?"
	result, tooCall := llm.Chat(prompt)
	if len(tooCall) != 0 {
		fmt.Println("tool call:", tooCall)
	}
	fmt.Println("result:", result)
}
