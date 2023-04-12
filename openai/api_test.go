package openai

import (
	"encoding/json"
	"testing"
)

var (
	testAPIKey = ""
	testOrgID  = ""
)

func TestChat(t *testing.T) {
	data, err := Chat("gpt-3.5-turbo", "how old are you", "user001", testAPIKey, testOrgID)
	if err != nil {
		t.Log(data)
		return
	}
	var dataByte []byte
	if dataByte, err = json.Marshal(data); err != nil {
		t.Log(data)
		return
	}
	t.Log(string(dataByte))
}

func TestCompletion(t *testing.T) {
	data, err := Completion("text-davinci-003", "how old are you", testAPIKey, testOrgID)
	if err != nil {
		t.Log(data)
		return
	}
	var dataByte []byte
	if dataByte, err = json.Marshal(data); err != nil {
		t.Log(data)
		return
	}
	t.Log(string(dataByte))
}

func TestEdit(t *testing.T) {
	data, err := Edit("text-davinci-edit-001", "how od ae y", "fix", testAPIKey, testOrgID)
	if err != nil {
		t.Log(data)
		return
	}
	var dataByte []byte
	if dataByte, err = json.Marshal(data); err != nil {
		t.Log(data)
		return
	}
	t.Log(string(dataByte))
}
