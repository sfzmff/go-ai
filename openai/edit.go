package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	editUrl = `https://api.openai.com/v1/edits`
)

type EditReqInfo struct {
	Model       string  `json:"model"`                 // 模型ID text-davinci-edit-001
	Input       string  `json:"input"`                 // 输入(待修改文本)
	Instruction string  `json:"instruction"`           // 指示
	MaxTokens   uint64  `json:"max_tokens,omitempty"`  // 最大令牌数(总令牌数) 4096
	Temperature float32 `json:"temperature,omitempty"` // 温度采样 0~2
	TopP        float32 `json:"top_p,omitempty"`       // 核心采样 0~1
	N           uint8   `json:"n,omitempty"`           // 生成聊天补全数量
}

type EditRespInfo struct {
	Object  string       `json:"object"`  // 对象 edit
	Created uint64       `json:"created"` // 创建时间 1681309029
	Choices []EditChoice `json:"choices"` // 回答
	Usage   EditUsage    `json:"usage"`   // 用量
}

type EditChoice struct {
	Text  string `json:"text"`  // 文本
	Index uint8  `json:"index"` // 序列(第几个回答，与请求中N相关)
}

type EditUsage struct {
	PromptTokens     uint16 `json:"prompt_tokens"`     // 提问令牌数
	CompletionTokens uint16 `json:"completion_tokens"` // 回答令牌数
	TotalTokens      uint16 `json:"total_tokens"`      // 总令牌数
}

// Edit 修改
// model,instruction,apiKey 必传
func Edit(model, input, instruction, apiKey, orgID string, proxy func(*http.Request) (*url.URL, error)) (data EditRespInfo, err error) {
	if len(strings.TrimSpace(model)) == 0 {
		err = fmt.Errorf("empty model")
		return
	} else if len(strings.TrimSpace(input)) == 0 {
		err = fmt.Errorf("empty input")
		return
	} else if len(strings.TrimSpace(instruction)) == 0 {
		err = fmt.Errorf("empty instruction")
		return
	} else if len(strings.TrimSpace(apiKey)) == 0 {
		err = fmt.Errorf("empty api_key")
		return
	}

	var dataByte []byte
	var req *http.Request
	var resp *http.Response

	reqData := EditReqInfo{
		Model:       model,
		Input:       input,
		Instruction: instruction,
		// MaxTokens:   1000,
		// Temperature: 1,
		// TopP:        1,
		// N:           1,
	}
	if dataByte, err = json.Marshal(reqData); err != nil {
		return
	}
	if req, err = http.NewRequest("POST", editUrl, bytes.NewBuffer(dataByte)); err != nil {
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
			Proxy: proxy,
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
