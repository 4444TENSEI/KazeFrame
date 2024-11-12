// 邮件发送服务
package service

import (
	"KazeFrame/internal/config"
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

// senderName 发件人昵称
// senderEmail 发件人邮箱
// smtpPassword 发件人授权码
// smtpServer 邮件服务地址
// smtpPort 邮件服务端口
// userEmail 收件人邮箱
// emailTitle 邮件标题
// emailContent 邮件内容

// 发送邮件通用服务
func SendMail(userEmail, emailTitle, emailContent string) error {
	// 从配置文件中获取邮件配置
	emailConfig := config.GetConfig().Email
	// 检查邮件服务是否启用
	if !emailConfig.Enable {
		return fmt.Errorf("邮件服务未启用")
	}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", emailConfig.SenderEmail, emailConfig.SenderName)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", emailTitle)
	m.SetBody("text/html", emailContent)
	d := gomail.NewDialer(emailConfig.SmtpServer, emailConfig.SmtpPort, emailConfig.SenderEmail, emailConfig.SenderPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		// 发送失败，返回错误
		return fmt.Errorf("发送邮件失败: %w", err)
	}
	// 邮件发送成功
	return nil
}
