package user

import (
	"GopherAI/common/code"
	"GopherAI/controller"
	serviceUser "GopherAI/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

var svc = serviceUser.NewService()

type (
	//这里的Username只能是账号登录，和我做的另一个项目有区别（邮箱账号均可)
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// omitempty当字段为空的时候，不返回这个东西
	LoginResponse struct {
		controller.Response
		Token string `json:"token,omitempty"`
	}
	//验证码由后端生成，存放到redis中，固然需要先发送一次请求CaptchaRequest,然后用返回的验证码
	//邮箱以及密码进行注册，后续再将账号进行返回
	RegisterRequest struct {
		Email    string `json:"email" binding:"required"`
		Captcha  string `json:"captcha"`
		Password string `json:"password"`
	}
	//注册成功之后，直接让其进行登录状态
	RegisterResponse struct {
		controller.Response
		Token string `json:"token,omitempty"`
	}

	CaptchaRequest struct {
		Email string `json:"email" binding:"required"`
	}

	CaptchaResponse struct {
		controller.Response
	}
)

// Login 用户登录
// @Summary      用户登录
// @Description  使用用户名和密码登录，返回 JWT Token
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "登录参数"
// @Success      200 {object} LoginResponse
// @Router       /user/login [post]

func Login(c *gin.Context) {

	req := new(LoginRequest)
	res := new(LoginResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	token, code_ := svc.Login(req.Username, req.Password)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)
}

// Register 用户注册
// @Summary      用户注册
// @Description  通过邮箱、验证码和密码注册账号，注册成功自动返回 JWT Token
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "注册参数"
// @Success      200 {object} RegisterResponse
// @Router       /user/register [post]
func Register(c *gin.Context) {

	req := new(RegisterRequest)
	res := new(RegisterResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	token, code_ := svc.Register(req.Email, req.Password, req.Captcha)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}

	res.Success()
	res.Token = token
	c.JSON(http.StatusOK, res)
}

// HandleCaptcha 发送验证码
// @Summary      发送邮箱验证码
// @Description  向指定邮箱发送验证码，用于注册流程
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        request body CaptchaRequest true "邮箱地址"
// @Success      200 {object} CaptchaResponse
// @Router       /user/captcha [post]
func HandleCaptcha(c *gin.Context) {
	req := new(CaptchaRequest)
	res := new(CaptchaResponse)
	//解析参数
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	//给service层进行处理
	code_ := svc.SendCaptcha(req.Email)
	if code_ != code.CodeSuccess {
		c.JSON(http.StatusOK, res.CodeOf(code_))
		return
	}
	//匿名字段，其实本身res.Success()调用就是res.Response.Success()
	//res.Response.Success()
	res.Success()
	c.JSON(http.StatusOK, res)
}
