package utils

import (
	"testing"

	"GopherAI/model"
	"github.com/cloudwego/eino/schema"
)

// 测试 MD5 加密
func TestMD5(t *testing.T) {
	// 已知输入和预期输出
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "5d41402abc4b2a76b9719d911017c592"},
		{"", "d41d8cd98f00b204e9800998ecf8427e"},
		{"GopherAI", "9b3fd2e653beeac4c869f3a60af8959f"},
	}

	for _, tt := range tests {
		result := MD5(tt.input)
		if result != tt.expected {
			t.Errorf("MD5(%q) = %q, 期望 %q", tt.input, result, tt.expected)
		}
	}
}

func TestGetRandomNumber(t *testing.T) {
	result := GetRandomNumbers(6)
	if len(result) != 6 {
		t.Errorf("GetRandomNumber(6) 返回长度 = %d,期望 6", len(result))
	}
	//检测是否全为数字
	for i, c := range result {
		if c < '0' || c > '9' {
			t.Errorf("第 %d个字符 %c 不是数字", i, c)
		}
	}
}

func TestGenerateUUID(t *testing.T) {
	id := GenerateUUID()
	if id == "" {
		t.Error("GenerateUUID() 返回空字符串")
	}
	// UUID 格式: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (36字符)
	if len(id) != 36 {
		t.Errorf("GenerateUUID() 长度 = %d, 期望 36", len(id))
	}
}

func TestConvertToModelMessage(t *testing.T) {
	sessionID := "test-session-id"
	userName := "test-user"
	msg := &schema.Message{
		Role:    schema.User,
		Content: "你好",
	}

	result := ConvertToModelMessage(sessionID, userName, msg)

	if result.SessionID != sessionID {
		t.Errorf("SessionID = %q, 期望 %q", result.SessionID, sessionID)
	}
	if result.UserName != userName {
		t.Errorf("UserName = %q, 期望 %q", result.UserName, userName)
	}
	if result.Content != msg.Content {
		t.Errorf("Content = %q, 期望 %q", result.Content, msg.Content)
	}
	if result.IsUser != false {
		t.Error("IsUser 应该为 false（默认值）")
	}
}

func TestConvertToSchemaMessages(t *testing.T) {
	msgs := []*model.Message{
		{IsUser: true, Content: "你好"},
		{IsUser: false, Content: "你好！有什么可以帮助你的吗？"},
	}

	result := ConvertToSchemaMessages(msgs)

	if len(result) != 2 {
		t.Fatalf("返回 %d 条消息, 期望 2", len(result))
	}

	if result[0].Role != schema.User {
		t.Errorf("第一条消息 Role = %v, 期望 User", result[0].Role)
	}
	if result[1].Role != schema.Assistant {
		t.Errorf("第二条消息 Role = %v, 期望 Assistant", result[1].Role)
	}
	if result[0].Content != "你好" || result[1].Content != "你好！有什么可以帮助你的吗？" {
		t.Error("Content 不匹配")
	}
}
