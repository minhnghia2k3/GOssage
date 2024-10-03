package mailer

import (
	"bytes"
	"fmt"
	"github.com/minhnghia2k3/GOssage/internal"
	"gopkg.in/gomail.v2"
	"html/template"
	"time"
)

type Mailer struct {
	fromEmail string
	host      string
	port      int
	username  string
	password  string
}

func NewMailer(fromEmail, host, username, password string, port int) *Mailer {
	return &Mailer{
		fromEmail: fromEmail,
		host:      host,
		port:      port,
		username:  username,
		password:  password,
	}
}

func (m *Mailer) Send(templateFile string, toEmail []string, data any) error {
	// parsing template
	tmpl, err := template.ParseFS(internal.FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", m.fromEmail)
	mail.SetHeader("To", toEmail...)
	mail.SetHeader("Subject", subject.String())
	mail.SetBody("text/plain", plainBody.String())
	mail.AddAlternative("text/html", htmlBody.String())

	d := gomail.NewDialer(m.host, m.port, m.username, m.password)

	// Send email with 3 retry times
	for i := 0; i < internal.MaxRetries; i++ {
		err = d.DialAndSend(mail)

		// If worked return result
		if err == nil {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	//After retries 3 times but fail
	return fmt.Errorf("failed to send email after %d attempt, error: %v", internal.MaxRetries, err)
}
