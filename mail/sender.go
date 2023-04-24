package mail

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddress, toEmailAddress string) EmailSender {

	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: toEmailAddress,
	}
}

func (g GmailSender) SendEmail(

	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,

) error {

	e := email.NewEmail()

	e.To = to
	e.Bcc = bcc
	e.Cc = cc
	e.Subject = subject
	e.From = fmt.Sprintf("%s <%s>", g.name, g.fromEmailAddress)
	e.HTML = []byte(content)

	for _, v := range attachFiles {
		_, err := e.AttachFile(v)
		if err != nil {
			return fmt.Errorf("fail to attach file %s %w", v, err)
		}
	}
	smtpAuth := smtp.PlainAuth("", g.fromEmailAddress, g.fromEmailPassword, smtpAuthAddress)

	return e.Send(smtpServerAddress, smtpAuth)

}
