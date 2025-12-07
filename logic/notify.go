package logic

import (
	"acat/setting"
	"fmt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"time"
)

// 生成6位数字验证码
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// 发送邮件验证码
func SendEmailCode(to, code string) error {
	if setting.Conf.EmailSMTPPass == "" {
		return fmt.Errorf("SMTP 未配置")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", setting.Conf.EmailSMTPEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "【ACAT纳新系统】密码重置验证码")
	m.SetBody("text/plain", fmt.Sprintf("您的验证码是：%s，5分钟内有效。", code))

	d := gomail.NewDialer(setting.Conf.EmailSMTPHost, 587, setting.Conf.EmailSMTPEmail, setting.Conf.EmailSMTPPass)
	fmt.Println(setting.Conf.EmailSMTPHost, 587, setting.Conf.EmailSMTPEmail, setting.Conf.EmailSMTPPass)
	return d.DialAndSend(m)
}
