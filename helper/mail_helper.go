package helper

import (
	"fmt"
	"net/smtp"

	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/config"
	"github.com/pajri/personal-backend/global"
)

type IEMail interface {
	SendMail(to []string, subject, body string) error
}

type Email struct {
	Mime string
}

func NewEmailHelper() IEMail {
	return Email{
		Mime: "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n",
	}
}

func (e Email) SendMail(to []string, subject, body string) error {
	auth := e.auth()
	address := e.smtpAddress()
	from := config.Config.SMTP.From
	message := fmt.Sprintf("%s\n%s\n%s", subject, e.Mime, body)

	err := smtp.SendMail(address, auth, from, to, []byte(message))
	if err != nil {
		return cerror.NewAndPrintWithTag("SMM00", err, global.FRIENDLY_MESSAGE)
	}

	return nil
}

func (e Email) smtpAddress() string {
	address := fmt.Sprintf("%s:%v", config.Config.SMTP.Host, config.Config.SMTP.Port)
	return address
}

func (e Email) auth() smtp.Auth {
	return smtp.PlainAuth("",
		config.Config.SMTP.Username,
		config.Config.SMTP.Password,
		config.Config.SMTP.Host)
}
