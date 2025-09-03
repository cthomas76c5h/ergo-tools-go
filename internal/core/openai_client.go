package core

import (
	"context"
	"encoding/json"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	c     *openai.Client
	model string
}

func NewOpenAI() *OpenAIClient {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		panic("OPENAI_API_KEY not set")
	}
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &OpenAIClient{c: openai.NewClient(key), model: model}
}

type ToolHandler func(name string, args map[string]any) (any, error)

func (o *OpenAIClient) RunWithTools(ctx context.Context, msgs []openai.ChatCompletionMessage, tools []openai.Tool, handle ToolHandler) (openai.ChatCompletionMessage, error) {
	resp, err := o.c.CreateChatCompletion(ctx, openai.ChatCompletionRequest{Model: o.model, Messages: msgs, Tools: tools, ToolChoice: "auto", Temperature: 0.2})
	if err != nil {
		return openai.ChatCompletionMessage{}, err
	}

	for {
		m := resp.Choices[0].Message
		if len(m.ToolCalls) == 0 {
			return m, nil
		}

		// append assistant message with tool calls
		msgs = append(msgs, m)

		// fulfill each tool call
		for _, tc := range m.ToolCalls {
			var args map[string]any
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				return openai.ChatCompletionMessage{}, err
			}
			out, err := handle(tc.Function.Name, args)
			if err != nil {
				return openai.ChatCompletionMessage{}, err
			}
			payload, _ := json.Marshal(out)
			msgs = append(msgs, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				ToolCallID: tc.ID,
				Name:       tc.Function.Name,
				Content:    string(payload),
			})
		}

		resp, err = o.c.CreateChatCompletion(ctx, openai.ChatCompletionRequest{Model: o.model, Messages: msgs, Tools: tools, ToolChoice: "auto", Temperature: 0.2})
		if err != nil {
			return openai.ChatCompletionMessage{}, err
		}
	}
}
