// ═══════════════════════════════════════════════════════════════
// Adaptador SMTP – Envío de correos electrónicos
// Implementa port.EmailSender usando gomail para SMTP
// ═══════════════════════════════════════════════════════════════
package email

import (
	"context"

	"github.com/cloudmart/notification-service/internal/domain/port"
	"gopkg.in/gomail.v2"
)

type smtpSender struct {
	dialer *gomail.Dialer
	from   string
}

func NewSMTPSender(host string, port int, username, password, from string) port.EmailSender {
	d := gomail.NewDialer(host, port, username, password)
	return &smtpSender{dialer: d, from: from}
}

func (s *smtpSender) SendEmail(ctx context.Context, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return s.dialer.DialAndSend(m)
}
