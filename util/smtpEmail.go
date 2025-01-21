package util

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func SendEmail(host string, username string, password string, to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(host, 25, username, password)

	// 关闭ssl
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
