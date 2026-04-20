package utils

import (
	"Yc-GopherAI/model"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// （我自己的生成数字逻辑，源代码性能不如这个好）
func GetRandomNumbers(num int) string {
	if num <= 0 {
		return ""
	}

	numb := make([]byte, 0, num)

	for i := 0; i < num; i++ {
		// 0~9随机数
		numb = append(numb, '0'+byte(r.Intn(10)))

	}
	return string(numb)
}

// MD5 MD5加密
func MD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func GenerateUUID() string {
	return uuid.New().String()
}

// 将 schema 消息转换为数据库可存储的格式
func ConvertToModelMessage(sessionID string, userName string, msg *schema.Message) *model.Message {
	return &model.Message{
		SessionID: sessionID,
		UserName:  userName,
		Content:   msg.Content,
	}
}

// 将数据库消息转换为 schema 消息（供 AI 使用）
func ConvertToSchemaMessages(msgs []*model.Message) []*schema.Message {
	schemaMsgs := make([]*schema.Message, 0, len(msgs))
	for _, m := range msgs {
		role := schema.Assistant
		if m.IsUser {
			role = schema.User
		}
		schemaMsgs = append(schemaMsgs, &schema.Message{
			Role:    role,
			Content: m.Content,
		})
	}
	return schemaMsgs
}
