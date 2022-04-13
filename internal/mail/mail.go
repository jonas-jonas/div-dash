package mail

import (
	"div-dash/internal/config"
	"net/smtp"
	"strconv"
)

type (
	mailServiceDependencies interface {
		config.ConfigProvider
	}

	MailServiceProvider interface {
		MailService() *MailService
	}
	MailService struct {
		smtpPassword string
		smtpServer   string
		smtpPort     int
		auth         smtp.Auth
	}
)

func NewMailService(m mailServiceDependencies) *MailService {
	config := m.Config().SMTP
	auth := smtp.PlainAuth("", "sender@div-dash.io", config.Password, config.Server) //TODO
	return &MailService{
		config.Password,
		config.Server,
		config.Port,
		auth,
	}
}

func (m *MailService) SendMail(to, from, subject, body string) error {
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + ".\r\n")
	return smtp.SendMail(m.smtpServer+":"+strconv.Itoa(m.smtpPort), m.auth, from, []string{to}, msg)
}
