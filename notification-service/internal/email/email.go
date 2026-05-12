package email

import (
	"bytes"
	"log/slog"
	"net/smtp"
	"strings"
	"text/template"
)

type JWTGenerator interface {
	GenerateToken(email string) (string, error)
}

type EmailSenderInterface interface {
	// vars: VerificationLink
	SendVerificationEmail(toAddr, subject, Template string, vars map[string]string) error
}

type EmailSender struct {
	FromEmail         string
	FromEmailPassword string
	FromEmailSMTP     string
	SMTPAddr          string
	log               *slog.Logger
}

func NewEmailSender(FromEmail, FromEmailPassword, FromEmailSMTP, SMTPAddr string, log *slog.Logger) EmailSender {
	return EmailSender{FromEmail: FromEmail, FromEmailPassword: FromEmailPassword,
		FromEmailSMTP: FromEmailSMTP, SMTPAddr: SMTPAddr, log: log}
}

func (e *EmailSender) sendHTMLEmail(to []string, subject, htmlbody string) error {
	auth := smtp.PlainAuth(
		"", e.FromEmail, e.FromEmailPassword, e.FromEmailSMTP,
	)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	message := "Subject: " + subject + "\n" + headers + "\n\n" + htmlbody
	return smtp.SendMail(e.SMTPAddr, auth, e.FromEmail, to, []byte(message))
}

// vars: VerificationLink
func (e *EmailSender) SendVerificationEmail(toAddr, subject, Template string, vars map[string]string) error {
	const op = "email.SendVerificationEmail"
	log := e.log.With(slog.String("op", op))
	log.Info("sending verification e-mail", slog.String("email", toAddr))
	tmpl, err := template.ParseFiles("./templates/" + Template + ".html")
	if err != nil {
		log.Error("template parsing error", slog.String("error", err.Error()))
		return err
	}

	// Converting comma separated strings
	to := strings.Split(toAddr, ",")

	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, vars)
	if err != nil {
		log.Error("executing template error", slog.String("error", err.Error()))
		return err
	}
	err = e.sendHTMLEmail(to, subject, rendered.String())
	if err != nil {
		log.Error("sending email error", slog.String("error", err.Error()))
		return err
	}
	return nil
}
