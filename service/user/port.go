package user

import "GopherAI/model"

// UserDAO 用户数据访问接口
type UserDAO interface {
	IsExistUser(username string) (bool, *model.User)
	Register(username, email, password string) (*model.User, bool)
}

// TokenService JWT Token 生成接口
type TokenService interface {
	GenerateToken(id int64, username string) (string, error)
}

// CaptchaService 验证码存储和校验接口
type CaptchaService interface {
	SetForEmail(email, captcha string) error
	CheckForEmail(email, userInput string) (bool, error)
}

// EmailService 邮件发送接口
type EmailService interface {
	SendCaptcha(to, code, msg string) error
}
