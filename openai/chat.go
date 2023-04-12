package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	chatUrl = `https://api.openai.com/v1/chat/completions`
)

type ChatReqInfo struct {
	Model            string         `json:"model"`                       // 模型ID gpt-3.5-turbo
	Messages         []MessagesData `json:"messages"`                    // 消息
	User             string         `json:"user,omitempty"`              // 用户
	MaxTokens        uint64         `json:"max_tokens,omitempty"`        // 最大令牌数(总令牌数) 4096
	Temperature      float32        `json:"temperature,omitempty"`       // 温度采样 0~2
	TopP             float32        `json:"top_p,omitempty"`             // 核心采样 0~1
	N                uint8          `json:"n,omitempty"`                 // 生成聊天补全数量
	Stream           bool           `json:"stream,omitempty"`            // 发送部分消息增量
	Stop             [4]string      `json:"stop,omitempty"`              // 停止生成令牌(生成时遇到即止)
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`  // 模型谈论新主题的可能性 -2.0~2.0
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"` // 模型逐字重复同一行的可能性 -2.0~2.0
}

type ChatRespInfo struct {
	ID      string       `json:"id"`      // 聊天ID chatcmpl-74R4hDTKxzh5EzElrhofHa5msxvKz
	Object  string       `json:"object"`  // 对象 chat.completion
	Created uint64       `json:"created"` // 创建时间 1681291743
	Model   string       `json:"model"`   // 模型ID gpt-3.5-turbo-0301
	Choices []ChatChoice `json:"choices"` // 回答
	Usage   ChatUsage    `json:"usage"`   // 用量
}

type MessagesData struct {
	Role    string `json:"role"`    // 角色 user/system/assistant
	Content string `json:"content"` // 内容
}

type ChatChoice struct {
	Messages     []MessagesData `json:"messages"`      // 消息
	FinishReason string         `json:"finish_reason"` // 完成原因(stop为回答完毕)
	Index        string         `json:"index"`         // 序列(第几个回答，与请求中N相关)
}

type ChatUsage struct {
	PromptTokens     uint16 `json:"prompt_tokens"`     // 提问令牌数
	CompletionTokens uint16 `json:"completion_tokens"` // 回答令牌数
	TotalTokens      uint16 `json:"total_tokens"`      // 总令牌数
}

// Chat 聊天补全
// model,content,apiKey 必传
func Chat(model, content, user, apiKey, orgID string) (data ChatRespInfo, err error) {
	if len(strings.TrimSpace(apiKey)) == 0 {
		err = fmt.Errorf("empty api_key")
		return
	}

	var dataByte []byte
	var req *http.Request
	var resp *http.Response

	reqData := ChatReqInfo{
		Model:            model,
		Messages:         []MessagesData{{Role: "user", Content: content}},
		User:             user,
		MaxTokens:        1000,
		Temperature:      1,
		TopP:             1,
		N:                1,
		Stream:           false,
		Stop:             [4]string{},
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	}
	if dataByte, err = json.Marshal(reqData); err != nil {
		return
	}
	if req, err = http.NewRequest("POST", chatUrl, bytes.NewBuffer(dataByte)); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	if len(strings.TrimSpace(orgID)) > 0 {
		req.Header.Set("OpenAI-Organization", orgID)
	}

	client := &http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			// Proxy: http.ProxyURL(fixedURL),
			DialContext: (&net.Dialer{
				Timeout: time.Second * 60, // 设置超时时间
			}).DialContext,
		},
	}
	if resp, err = client.Do(req); err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if dataByte, err = io.ReadAll(resp.Body); err != nil {
		return
	}
	if err = json.Unmarshal(dataByte, &data); err != nil {
		return
	}

	return
}
