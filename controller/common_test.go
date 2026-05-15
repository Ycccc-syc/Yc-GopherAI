package controller

import (
	"GopherAI/common/code"
	"testing"
)

func TestResponse(t *testing.T) {
	res := new(Response)
	res.Success()

	if res.StatusCode != code.CodeSuccess {
		t.Errorf("StatusCode=%d,期望 %d", res.StatusCode, code.CodeSuccess)
	}
	if res.StatusMsg != "success" {
		t.Errorf("StatusMsg = %q, 期望 %q", res.StatusMsg, "success")
	}
}

func TestResponseCodeOf(t *testing.T) {
	res := new(Response)
	result := res.CodeOf(code.CodeUserNotExist)

	if result.StatusCode != code.CodeUserNotExist {
		t.Errorf("StatusCode = %d, 期望 %d", result.StatusCode, code.CodeUserNotExist)
	}
	if result.StatusMsg != "用户不存在" {
		t.Errorf("StatusMsg = %q, 期望 %q", result.StatusMsg, "用户不存在")
	}
	// 原对象也被修改了
	if res.StatusCode != code.CodeUserNotExist {
		t.Error("CodeOf 应该也修改了原对象")
	}
}
