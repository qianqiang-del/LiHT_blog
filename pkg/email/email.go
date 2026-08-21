package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"blog/pkg/config"
)

// Sender 邮件发送器
type Sender struct {
	cfg *config.EmailConfig
}

// NewSender 创建邮件发送器
func NewSender(cfg *config.EmailConfig) *Sender {
	return &Sender{cfg: cfg}
}

// Send 发送邮件（SSL）
func (s *Sender) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.cfg.From, to, subject, body))

	// 建立 TLS 连接
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: s.cfg.Host,
	})
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %w", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return fmt.Errorf("创建 SMTP 客户端失败: %w", err)
	}
	defer client.Close()

	// 设置发件人和收件人
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("认证失败: %w", err)
	}

	if err = client.Mail(s.cfg.From); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败: %w", err)
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("准备发送失败: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件失败: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭写入失败: %w", err)
	}

	return client.Quit()
}

// SendCode 发送验证码邮件
func (s *Sender) SendCode(to, code string) error {
	subject := "【LiHT Blog】邮箱验证码"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">邮箱验证码</h2>
			<p>您好，您正在注册 LiHT Blog 账号。</p>
			<p>您的验证码是：</p>
			<div style="background: #f5f5f5; padding: 20px; text-align: center; font-size: 32px; font-weight: bold; color: #333; letter-spacing: 5px;">
				%s
			</div>
			<p style="color: #999; font-size: 14px;">验证码 1 分钟内有效，请勿泄露给他人。</p>
		</div>
	`, code)

	return s.Send(to, subject, body)
}
