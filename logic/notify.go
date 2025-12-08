package logic

import (
	"acat/setting"
	"fmt"
	"github.com/go-mail/mail/v2"
	"math/rand"
)

// 生成6位数字验证码（保留原逻辑，仅移除 rand.Seed）
func GenerateCode() string {
	// Go 1.20+ 全局 rand 已自动 seed，无需手动调用 Seed
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// 发送邮件验证码
func SendEmailCode(to, code string) error {
	if setting.Conf.EmailSMTPPass == "" {
		return fmt.Errorf("SMTP 未配置")
	}

	m := mail.NewMessage()
	m.SetHeader("From", setting.Conf.EmailSMTPEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "【ACAT纳新系统】密码重置验证码")
	m.SetBody("text/plain", fmt.Sprintf("您的验证码是：%s，5分钟内有效。", code))

	d := mail.NewDialer(setting.Conf.EmailSMTPHost, 587, setting.Conf.EmailSMTPEmail, setting.Conf.EmailSMTPPass)
	return d.DialAndSend(m)
}

// PublicEmail 用来通知用户面试结果。以默认格式发送
func PublicEmail(to string, round int, name string) error {
	if setting.Conf.EmailSMTPPass == "" {
		return fmt.Errorf("SMTP 未配置")
	}
	m := mail.NewMessage()
	m.SetHeader("From", setting.Conf.EmailSMTPEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "【ACAT纳新系统】面试结果通知")
	m.SetBody("text/plain", fmt.Sprintf("同学 %s，你好！恭喜你通过第 %d 轮面试，请留意后续通知。", name, round))

	d := mail.NewDialer(setting.Conf.EmailSMTPHost, 587, setting.Conf.EmailSMTPEmail, setting.Conf.EmailSMTPPass)
	return d.DialAndSend(m)
}

// PublicCustomEmail 自定义发送文本
func PublicCustomEmail(to, content string) error {
	if setting.Conf.EmailSMTPPass == "" {
		return fmt.Errorf("SMTP 未配置")
	}
	m := mail.NewMessage()
	m.SetHeader("From", setting.Conf.EmailSMTPEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "【ACAT纳新系统】面试结果通知")
	m.SetBody("text/plain", content)
	d := mail.NewDialer(setting.Conf.EmailSMTPHost, 587, setting.Conf.EmailSMTPEmail, setting.Conf.EmailSMTPPass)
	return d.DialAndSend(m)
}
