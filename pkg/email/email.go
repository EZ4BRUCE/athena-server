package email

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
}

// 发送指定标题和内容的邮件至邮箱组
func (e *Mailer) SendMail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.From, "Athena 监控系统"))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetHeader("Cc", e.From)
	m.SetBody("text/html", body)
	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
	return dialer.DialAndSend(m)
}
