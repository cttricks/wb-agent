package controllers

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"wb-agent/utils"

	"github.com/sashabaranov/go-openai"
)

func ComposeSQL(apiKey, nl string) string {
	client := openai.NewClient(apiKey)
	ctx := context.Background()

	prompt, err := utils.SystemPrompt()
	if err != nil {
		fmt.Println("[ERROR] Unable to find Instructions.md\n", err.Error())
		return ""
	}

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: prompt},
			{Role: "user", Content: nl},
		},
		MaxTokens: 400,
	})
	if err != nil {
		fmt.Println("[ERROR] OpenAI response says:", err.Error())
		return ""
	}
	if len(resp.Choices) == 0 {
		return ""
	}
	raw := resp.Choices[0].Message.Content
	return sanitizeSQL(raw)
}

// Strip code fences and find the SQL statement (best-effort)
func sanitizeSQL(s string) string {
	s = strings.TrimSpace(s)

	// if wrapped in triple-backticks, extract content between first and last fence
	if strings.Contains(s, "```") {
		first := strings.Index(s, "```")
		last := strings.LastIndex(s, "```")
		if first != -1 && last != -1 && last > first {
			// cut out fences
			inside := s[first+3 : last]
			// if language given like ```sql\n..., remove leading token and newline
			inside = strings.TrimPrefix(inside, "sql\n")
			s = strings.TrimSpace(inside)
		}
	}

	// find first SQL keyword (SELECT/INSERT/UPDATE/DELETE/SHOW/CREATE/ALTER/DROP/WITH/EXPLAIN)
	re := regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE|SHOW|CREATE|ALTER|DROP|WITH|EXPLAIN)`)
	loc := re.FindStringIndex(s)
	if loc != nil {
		s = strings.TrimSpace(s[loc[0]:])
	}

	// ensure it ends with semicolon
	if !strings.HasSuffix(s, ";") {
		s = s + ";"
	}
	return s
}
