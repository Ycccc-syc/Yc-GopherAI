package user

import (
	"errors"
	"testing"

	"GopherAI/common/code"
	"GopherAI/model"
)

var errSome = errors.New("some error")

// ====== Mock ======

type mockUserDAO struct {
	users      map[string]*model.User
	registerOK bool
}

func (m *mockUserDAO) IsExistUser(username string) (bool, *model.User) {
	u, ok := m.users[username]
	return ok, u
}

func (m *mockUserDAO) Register(username, email, password string) (*model.User, bool) {
	if !m.registerOK {
		return nil, false
	}
	return &model.User{
		Username: username,
		Email:    email,
		Password: password,
	}, true
}

type mockTokenGen struct {
	token string
	err   error
}

func (m *mockTokenGen) GenerateToken(_ int64, _ string) (string, error) {
	return m.token, m.err
}

type mockCaptcha struct {
	codes map[string]string
}

func (m *mockCaptcha) SetForEmail(email, captcha string) error {
	if m.codes == nil {
		m.codes = make(map[string]string)
	}
	m.codes[email] = captcha
	return nil
}

func (m *mockCaptcha) CheckForEmail(email, userInput string) (bool, error) {
	code, ok := m.codes[email]
	return ok && code == userInput, nil
}

type mockEmail struct {
	err error
}

func (m *mockEmail) SendCaptcha(_, _, _ string) error {
	return m.err
}

// ====== Login 测试 ======

func TestLogin_Success(t *testing.T) {
	svc := &Service{
		UserDAO: &mockUserDAO{
			users: map[string]*model.User{
				"testuser": {Username: "testuser", Password: "e10adc3949ba59abbe56e057f20f883e"},
			},
		},
		TokenGen: &mockTokenGen{token: "fake-jwt-token"},
	}

	token, c := svc.Login("testuser", "123456")
	if c != code.CodeSuccess {
		t.Errorf("期望 CodeSuccess, 得到 %d", c)
	}
	if token != "fake-jwt-token" {
		t.Errorf("期望 fake-jwt-token, 得到 %s", token)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	svc := &Service{
		UserDAO:  &mockUserDAO{users: map[string]*model.User{}},
		TokenGen: &mockTokenGen{},
	}

	_, c := svc.Login("nonexistent", "123456")
	if c != code.CodeUserNotExist {
		t.Errorf("期望 CodeUserNotExist, 得到 %d", c)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := &Service{
		UserDAO: &mockUserDAO{
			users: map[string]*model.User{
				"testuser": {Username: "testuser", Password: "e10adc3949ba59abbe56e057f20f883e"},
			},
		},
		TokenGen: &mockTokenGen{},
	}

	_, c := svc.Login("testuser", "wrongpass")
	if c != code.CodeInvalidPassword {
		t.Errorf("期望 CodeInvalidPassword, 得到 %d", c)
	}
}

// ====== Register 测试 ======

func TestRegister_UserExist(t *testing.T) {
	svc := &Service{
		UserDAO: &mockUserDAO{
			users: map[string]*model.User{
				"old@test.com": {Email: "old@test.com"},
			},
		},
		Captcha: &mockCaptcha{},
		Email:   &mockEmail{},
	}

	_, c := svc.Register("old@test.com", "123456", "000000")
	if c != code.CodeUserExist {
		t.Errorf("期望 CodeUserExist, 得到 %d", c)
	}
}

func TestRegister_InvalidCaptcha(t *testing.T) {
	svc := &Service{
		UserDAO: &mockUserDAO{users: map[string]*model.User{}},
		Captcha: &mockCaptcha{codes: map[string]string{}},
		Email:   &mockEmail{},
	}

	_, c := svc.Register("new@test.com", "123456", "wrong")
	if c != code.CodeInvalidCaptcha {
		t.Errorf("期望 CodeInvalidCaptcha, 得到 %d", c)
	}
}

func TestRegister_DaoFail(t *testing.T) {
	svc := &Service{
		UserDAO: &mockUserDAO{users: map[string]*model.User{}, registerOK: false},
		Captcha: &mockCaptcha{codes: map[string]string{"new@test.com": "000000"}},
		Email:   &mockEmail{},
	}

	_, c := svc.Register("new@test.com", "123456", "000000")
	if c != code.CodeServerBusy {
		t.Errorf("期望 CodeServerBusy, 得到 %d", c)
	}
}

func TestRegister_EmailFail(t *testing.T) {
	svc := &Service{
		UserDAO: &mockUserDAO{users: map[string]*model.User{}, registerOK: true},
		Captcha: &mockCaptcha{codes: map[string]string{"new@test.com": "000000"}},
		Email:   &mockEmail{err: errSome},
	}

	_, c := svc.Register("new@test.com", "123456", "000000")
	if c != code.CodeServerBusy {
		t.Errorf("期望 CodeServerBusy, 得到 %d", c)
	}
}

func TestRegister_TokenFail(t *testing.T) {
	svc := &Service{
		UserDAO:  &mockUserDAO{users: map[string]*model.User{}, registerOK: true},
		Captcha:  &mockCaptcha{codes: map[string]string{"new@test.com": "000000"}},
		Email:    &mockEmail{},
		TokenGen: &mockTokenGen{err: errSome},
	}

	_, c := svc.Register("new@test.com", "123456", "000000")
	if c != code.CodeServerBusy {
		t.Errorf("期望 CodeServerBusy, 得到 %d", c)
	}
}

func TestRegister_Success(t *testing.T) {
	svc := &Service{
		UserDAO:  &mockUserDAO{users: map[string]*model.User{}, registerOK: true},
		Captcha:  &mockCaptcha{codes: map[string]string{"new@test.com": "000000"}},
		Email:    &mockEmail{},
		TokenGen: &mockTokenGen{token: "jwt-for-new-user"},
	}

	token, c := svc.Register("new@test.com", "123456", "000000")
	if c != code.CodeSuccess {
		t.Errorf("期望 CodeSuccess, 得到 %d", c)
	}
	if token != "jwt-for-new-user" {
		t.Errorf("期望 jwt-for-new-user, 得到 %s", token)
	}
}

// ====== SendCaptcha 测试 ======

func TestSendCaptcha_Success(t *testing.T) {
	svc := &Service{
		Captcha: &mockCaptcha{},
		Email:   &mockEmail{},
	}

	c := svc.SendCaptcha("test@test.com")
	if c != code.CodeSuccess {
		t.Errorf("期望 CodeSuccess, 得到 %d", c)
	}
}
