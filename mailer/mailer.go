package mailer

import (
	"crypto/tls"
	"log"
	"strconv"

	"github.com/BogoCvetkov/go_mastercalss/config"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
	config *config.Config
}

func NewMailer(config *config.Config) *Mailer {
	host := config.SMTP_HOST
	port, err := strconv.Atoi(config.SMTP_PORT)
	user := config.SMTP_USER
	pass := config.SMTP_PASS

	if err != nil {
		log.Fatal("Wrong port variable")
	}

	d := gomail.NewDialer(host, port, user, pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &Mailer{dialer: d, config: config}
}

func (m *Mailer) NewMail(to string, sub string, content string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.config.SMTP_FROM)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", sub)
	msg.SetBody("text/html", content)

	if err := m.dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
