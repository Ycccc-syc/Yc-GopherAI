package user

import (
	myemail "GopherAI/common/email"
	myredis "GopherAI/common/redis"
	userdao "GopherAI/dao/user"
	"GopherAI/model"
	"GopherAI/utils/myjwt"
)

// ====== UserDAO 默认实现 ======

type userDAOImpl struct{}

func (d *userDAOImpl) IsExistUser(username string) (bool, *model.User) {
	return userdao.IsExistUser(username)
}

func (d *userDAOImpl) Register(username, email, password string) (*model.User, bool) {
	return userdao.Register(username, email, password)
}

// ====== TokenService 默认实现 ======

type jwtImpl struct{}

func (j *jwtImpl) GenerateToken(id int64, username string) (string, error) {
	return myjwt.GenerateToken(id, username)
}

// ====== CaptchaService 默认实现 ======

type captchaImpl struct{}

func (c *captchaImpl) SetForEmail(email, captcha string) error {
	return myredis.SetCaptchaForEmail(email, captcha)
}

func (c *captchaImpl) CheckForEmail(email, userInput string) (bool, error) {
	return myredis.CheckCaptchaForEmail(email, userInput)
}

// ====== EmailService 默认实现 ======

type emailImpl struct{}

func (e *emailImpl) SendCaptcha(to, code, msg string) error {
	return myemail.SendCaptcha(to, code, msg)
}
