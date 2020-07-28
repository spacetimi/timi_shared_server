package email_utils

import (
	"errors"
	"net/smtp"
	"strconv"

	"github.com/spacetimi/timi_shared_server/v2/utils/logger"
)

type Emailer struct {
	Account           EmailAccount
	SmtpHost          string
	SmtpServerAddress string
}

type Email struct {
	Subject string
	Body    string
}

type EmailAccount struct {
	EmailAddress string
	Password     string
}

func NewEmailer(account EmailAccount, smtpHost string, smtpPort int) *Emailer {
	mailer := &Emailer{}
	mailer.SmtpHost = smtpHost
	mailer.SmtpServerAddress = smtpHost + ":" + strconv.Itoa(smtpPort)
	mailer.Account = account

	return mailer
}

func (mailer *Emailer) SendEmail(to string, email Email) error {
	fromAddress := mailer.Account.EmailAddress
	password := mailer.Account.Password

	msg := "From: " + fromAddress + "\n" +
		"To: " + to + "\n" +
		"Subject: " + email.Subject + "\n\n" +
		email.Body

	auth := smtp.PlainAuth("", fromAddress, password, mailer.SmtpHost)

	err := smtp.SendMail(mailer.SmtpServerAddress, auth, fromAddress, []string{to}, []byte(msg))

	if err != nil {
		logger.LogError("error sending email" +
			"|to=" + to +
			"|from=" + fromAddress +
			"|error=" + err.Error())
		return errors.New("error sending email")
	}

	return nil
}
