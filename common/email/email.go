package email

import (
	"Yc-GopherAI/config"
	"fmt"

	"gopkg.in/gomail.v2"
)

const (
	CodeMsg     = "GopherAI验证码如下(验证码仅限于2分钟有效): "
	UserNameMsg = "GopherAI的账号如下，请保留好，后续可以用账号/邮箱登录 "
)

func SendCaptcha(email, code, msg string) error {
	m := gomail.NewMessage()

	//发件人
	m.SetHeader("From", config.GetConfig().EmailConfig.Email)
	//收件人
	m.SetHeader("To", email)
	//主题
	m.SetHeader("Subject", "来自GoperAI的信息")
	//正文内容（纯文本形式，也可以用text/html)
	m.SetBody("text/plain", msg+" "+code)

	//配置SMTP服务器和授权码，587：是SMTP的明文/STATRTTLS端口号
	d := gomail.NewDialer("smtp.qq.com", 587, config.GetConfig().Email, config.GetConfig().EmailConfig.Authcode)

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("DiaAndSend err %v:\n", err)
		return err
	}
	fmt.Printf("send mail success/n")
	return nil
}
