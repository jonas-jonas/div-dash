package mail

import (
	"net/smtp"
	"strconv"
)

type MailService struct {
	smtpPassword string
	smtpServer   string
	smtpPort     int
	auth         smtp.Auth
}

func NewMailService(smtpPassword, smtpServer string, smtpPort int) *MailService {
	auth := smtp.PlainAuth("", "sender@div-dash.io", smtpPassword, smtpServer)
	return &MailService{
		smtpPassword,
		smtpServer,
		smtpPort,
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
