package helper

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/config"
	"github.com/pajri/personal-backend/global"
)

type IEMail interface {
	SendMail(to []string, subject, body string) error
}

type Email struct {
	Mime        string
	ContentType string
}

func NewEmailHelper() IEMail {
	return Email{
		Mime:        "MIME-version: 1.0",
		ContentType: "Content-Type: text/html; charset=UTF-8",
	}
}

func (e Email) SendMail(to []string, subject, body string) error {
	auth := e.auth()
	address := e.smtpAddress()
	from := config.Config.SMTP.From
	message := e.message(from, to, subject, body)

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

func (e Email) message(from string, to []string, subject, body string) string {
	toList := strings.Join(to, ",")

	message := fmt.Sprintf("%s\r\n%s\r\nFrom: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.Mime, e.ContentType, from, toList, subject, body)

	fmt.Println(message)
	return message
}
