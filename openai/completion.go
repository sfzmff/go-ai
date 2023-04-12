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
	completionUrl = `https://api.openai.com/v1/completions`
)

type CompletionReqInfo struct {
	Model       string      `json:"model"`       // 模型ID gpt-3.5-turbo
	Prompt      string      `json:"prompt"`      // 消息
	MaxTokens   uint64      `json:"max_tokens"`  // 最大令牌数(总令牌数) 4096
	Temperature float32     `json:"temperature"` // 温度采样 0~2
	TopP        float32     `json:"top_p"`       // 核心采样 0~1
	N           uint8       `json:"n"`           // 生成聊天补全数量
	Stream      bool        `json:"stream"`      // 发送部分消息增量
	Stop        [4]string   `json:"stop"`        // 停止生成令牌(生成时遇到即止)
	Logprobs    interface{} `json:"logprobs"`    // 最可能令牌的对数概率
}

type CompletionRespInfo struct {
	ID      string             `json:"id"`      // 聊天ID cmpl-74VDsZIpsz5lyveWiQmOD9xxfqt6x
	Object  string             `json:"object"`  // 对象 text_completion
	Created uint64             `json:"created"` // 创建时间 1681307688
	Model   string             `json:"model"`   // 模型ID text-davinci-003
	Choices []CompletionChoice `json:"choices"` // 回答
	Usage   CompletionUsage    `json:"usage"`   // 用量
}

type CompletionChoice struct {
	Text         string      `json:"text"`          // 文本
	Index        string      `json:"index"`         // 序列(第几个回答，与请求中N相关)
	FinishReason string      `json:"finish_reason"` // 完成原因(stop为回答完毕)
	Logprobs     interface{} `json:"logprobs"`      // 最可能令牌的对数概率
}

type CompletionUsage struct {
	PromptTokens     uint16 `json:"prompt_tokens"`     // 提问令牌数
	CompletionTokens uint16 `json:"completion_tokens"` // 回答令牌数
	TotalTokens      uint16 `json:"total_tokens"`      // 总令牌数
}

// Completion 补全
// apiKey 必传
func Completion(model, prompt, user, apiKey, orgID string) (data CompletionRespInfo, err error) {
	if len(strings.TrimSpace(apiKey)) == 0 {
		err = fmt.Errorf("empty api_key")
		return
	}

	var dataByte []byte
	var req *http.Request
	var resp *http.Response

	reqData := CompletionReqInfo{
		Model:       model,
		Prompt:      prompt,
		MaxTokens:   1000,
		Temperature: 1,
		TopP:        1,
		N:           1,
		Stream:      false,
		Stop:        [4]string{},
		Logprobs:    nil,
	}
	if dataByte, err = json.Marshal(reqData); err != nil {
		return
	}
	if req, err = http.NewRequest("POST", completionUrl, bytes.NewBuffer(dataByte)); err != nil {
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
