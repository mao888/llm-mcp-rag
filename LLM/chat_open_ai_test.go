package LLM

import (
	"context"
	"fmt"
	"testing"
)

func TestChatOpenAI_Chat(t *testing.T) {
	ctx := context.Background()
	//modeName := openai.ChatModelGPT4_1Mini
	modeName := "deepseek/deepseek-chat"
	llm := NewChatOpenAI(ctx, modeName)
	prompt := "hello, what is your name? and what can you do? use Chinese to answer me."
	result, tooCall := llm.Chat(prompt)
	if len(tooCall) != 0 {
		fmt.Println("tool call:", tooCall)
	}
	fmt.Println("result:", result)
}
