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
	listModelUrl     = `https://api.openai.com/v1/models`
	retrieveModelUrl = `https://api.openai.com/v1/models/`
)

type ModelList struct {
	Data   []ModelData `json:"data"`
	Object string      `json:"object"`
}

type ModelData struct {
	ID         string           `json:"id"`       // 模型ID gpt-3.5-turbo
	Object     string           `json:"object"`   // 对象 model
	Created    uint64           `json:"created"`  // 创建时间 1677610602
	OwnedBy    string           `json:"owned_by"` // 归属 openai
	Permission []PermissionData `json:"permission"`
}

type PermissionData struct {
	ID                 string      `json:"id"`
	Object             string      `json:"object"`
	Created            uint64      `json:"created"`
	AllowCreateEngine  bool        `json:"allow_create_engine"`
	AllowSampling      bool        `json:"allow_sampling"`
	AllowLogprobs      bool        `json:"allow_logprobs"`
	AllowSearchIndices bool        `json:"allow_search_indices"`
	AllowView          bool        `json:"allow_view"`
	AllowFineTuning    bool        `json:"allow_fine_tuning"`
	Organization       string      `json:"organization"`
	Group              interface{} `json:"group"`
	IsBlocking         bool        `json:"is_blocking"`
}

// ListModel 获取模型列表
// apiKey 必传
func ListModel(apiKey, orgID string) (data ModelList, err error) {
	if len(strings.TrimSpace(apiKey)) == 0 {
		err = fmt.Errorf("empty api_key")
		return
	}

	var dataByte []byte
	var req *http.Request
	var resp *http.Response

	if req, err = http.NewRequest("POST", listModelUrl, bytes.NewBuffer(dataByte)); err != nil {
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

// RetrieveModel 获取指定模型
// apiKey,model 必传
func RetrieveModel(apiKey, orgID, model string) (data ModelData, err error) {
	if len(strings.TrimSpace(model)) == 0 {
		err = fmt.Errorf("empty model")
		return
	} else if len(strings.TrimSpace(apiKey)) == 0 {
		err = fmt.Errorf("empty api_key")
		return
	}

	var dataByte []byte
	var req *http.Request
	var resp *http.Response

	if req, err = http.NewRequest("POST", retrieveModelUrl+model, nil); err != nil {
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
