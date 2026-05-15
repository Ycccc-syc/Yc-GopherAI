package code

import (
	"testing"
)

func TestCodeValue(t *testing.T) {
	if CodeSuccess != 1000 {
		t.Errorf("CodeSuccess = %d, 期望 1000", CodeSuccess)
	}
	if CodeInvalidParams != 2001 {
		t.Errorf("CodeInvalidParams = %d, 期望 2001", CodeInvalidParams)
	}
}

func TestCodeMsg(t *testing.T) {
	tests := []struct {
		code     Code
		expected string
	}{
		{CodeSuccess, "success"},
		{CodeInvalidParams, "请求参数错误"},
		{CodeUserNotExist, "用户不存在"},
		{CodeInvalidCaptcha, "验证码错误"},
		{CodeServerBusy, "服务繁忙"},
	}

	for _, tt := range tests {
		result := tt.code.Msg()
		if result != tt.expected {
			t.Errorf("%d.Msg() = %q, 期望 %q", tt.code, result, tt.expected)
		}
	}
}
