package user

import (
	"GopherAI/common/code"
	"GopherAI/model"
	"GopherAI/utils"
	userdao "GopherAI/dao/user"
)

// Service 用户服务，依赖通过接口注入
type Service struct {
	UserDAO  UserDAO
	TokenGen TokenService
	Captcha  CaptchaService
	Email    EmailService
}

func NewService() *Service {
	return &Service{
		UserDAO:  &userDAOImpl{},
		TokenGen: &jwtImpl{},
		Captcha:  &captchaImpl{},
		Email:    &emailImpl{},
	}
}

func (s *Service) Login(username, password string) (string, code.Code) {
	var userInformation *model.User
	var ok bool
	//1:判断用户是否存在
	if ok, userInformation = s.UserDAO.IsExistUser(username); !ok {
		return "", code.CodeUserNotExist
	}
	//2:判断用户是否密码账号正确
	if userInformation.Password != utils.MD5(password) {
		return "", code.CodeInvalidPassword
	}
	//3:返回一个Token
	token, err := s.TokenGen.GenerateToken(userInformation.ID, userInformation.Username)
	if err != nil {
		return "", code.CodeServerBusy
	}
	return token, code.CodeSuccess
}

func (s *Service) Register(email, password, captcha string) (string, code.Code) {
	var ok bool
	var userInformation *model.User

	//1:先判断用户是否已经存在了
	if ok, _ = s.UserDAO.IsExistUser(email); ok {
		return "", code.CodeUserExist
	}

	//2:从redis中验证验证码是否有效
	if ok, _ = s.Captcha.CheckForEmail(email, captcha); !ok {
		return "", code.CodeInvalidCaptcha
	}

	//3：生成11位的账号
	username := utils.GetRandomNumbers(11)

	//4：注册到数据库中
	if userInformation, ok = s.UserDAO.Register(username, email, password); !ok {
		return "", code.CodeServerBusy
	}

	//5：将账号一并发送到对应邮箱上去，后续需要账号登录
	if err := s.Email.SendCaptcha(email, username, userdao.UserNameMsg); err != nil {
		return "", code.CodeServerBusy
	}

	//6:生成Token
	token, err := s.TokenGen.GenerateToken(userInformation.ID, userInformation.Username)
	if err != nil {
		return "", code.CodeServerBusy
	}

	return token, code.CodeSuccess
}

// SendCaptcha 往指定邮箱发送验证码
func (s *Service) SendCaptcha(email_ string) code.Code {
	send_code := utils.GetRandomNumbers(6)
	//1:先存放到redis
	if err := s.Captcha.SetForEmail(email_, send_code); err != nil {
		return code.CodeServerBusy
	}
	//2:再进行远程发送
	if err := s.Email.SendCaptcha(email_, send_code, userdao.CodeMsg); err != nil {
		return code.CodeServerBusy
	}
	return code.CodeSuccess
}
