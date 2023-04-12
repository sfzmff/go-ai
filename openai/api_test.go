package openai

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

const (
	testAPIKey = ""
	testOrgID  = ""
	proxyUrl   = ""
)

func TestChat(t *testing.T) {
	fixedURL, _ := url.Parse(proxyUrl) // 须使用代理
	data, err := Chat("gpt-3.5-turbo", "how old are you?", "user001", testAPIKey, testOrgID, http.ProxyURL(fixedURL))
	if err != nil {
		t.Log(err.Error())
		return
	}
	var dataByte []byte
	if dataByte, err = json.Marshal(data); err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(string(dataByte))
}

func TestCompletion(t *testing.T) {
	fixedURL, _ := url.Parse(proxyUrl) // 须使用代理
	data, err := Completion("text-davinci-003", "how old are you?", testAPIKey, testOrgID, http.ProxyURL(fixedURL))
	if err != nil {
		t.Log(err.Error())
		return
	}
	var dataByte []byte
	if dataByte, err = json.Marshal(data); err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(string(dataByte))
}

func TestEdit(t *testing.T) {
	fixedURL, _ := url.Parse(proxyUrl) // 须使用代理
	data, err := Edit("text-davinci-edit-001", "how old ae yo?", "fix the spelling mistakes", testAPIKey, testOrgID, http.ProxyURL(fixedURL))
	if err != nil {
		t.Log(err.Error())
		return
	}
	var dataByte []byte
	if dataByte, err = json.Marshal(data); err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(string(dataByte))
}
