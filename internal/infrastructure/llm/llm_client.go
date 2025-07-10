package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/sashabaranov/go-openai"
)

type LLMClient struct {
	client *openai.Client
	tools  []Tool
}

func NewLLMClient(tools []Tool) *LLMClient {
	return &LLMClient{
		client: openai.NewClient(os.Getenv("OPENAI_API_KEY")),
		tools:  tools,
	}
}

func (c *LLMClient) ProcessMessage(input string) (string, error) {
	toolDefinitions := c.buildToolDefinitions()

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "あなたは家計管理の専門家です。ユーザーの支出データを分析し、実用的で実現可能なアドバイスを提供してください。",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		},
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
			Tools:    toolDefinitions,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	message := resp.Choices[0].Message
	log.Println("============ choise message =============")
	spew.Dump(message)

	if len(message.ToolCalls) > 0 {
		return c.handleToolCalls(message, messages)
	}

	return message.Content, nil
}

func (c *LLMClient) buildToolDefinitions() []openai.Tool {
	var tools []openai.Tool

	for _, tool := range c.tools {
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        tool.Name(),
				Description: tool.Description(),
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"household_id": map[string]interface{}{
							"type":        "integer",
							"description": "家計簿ID",
						},
					},
					"required": []string{"household_id"},
				},
			},
		})
	}

	return tools
}

func (c *LLMClient) handleToolCalls(message openai.ChatCompletionMessage, messages []openai.ChatCompletionMessage) (string, error) {
	messages = append(messages, message)

	for _, toolCall := range message.ToolCalls {
		result, err := c.executeToolCall(toolCall)
		if err != nil {
			return "", fmt.Errorf("failed to execute tool call: %w", err)
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    result,
			ToolCallID: toolCall.ID,
		})
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to create final chat completion: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *LLMClient) executeToolCall(toolCall openai.ToolCall) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
		return "", fmt.Errorf("failed to unmarshal tool call arguments: %w", err)
	}

	log.Println("============ toolCall param =============")
	spew.Dump(params)

	for _, tool := range c.tools {
		if tool.Name() == toolCall.Function.Name {
			result, err := tool.Execute(params)
			if err != nil {
				return "", fmt.Errorf("tool execution failed: %w", err)
			}

			resultJSON, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal tool result: %w", err)
			}

			return string(resultJSON), nil
		}
	}

	return "", fmt.Errorf("tool not found: %s", toolCall.Function.Name)
}
