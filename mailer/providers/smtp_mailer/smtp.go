package smtp_mailer

import (
	"errors"
	"github.com/ivanjaros/ijlibs/mailer"
	"github.com/jordan-wright/email"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"time"
)

const Provider = "smtp"

// if poolCount is less than 1, no pool will be created
func New(from mail.Address, addr string, auth smtp.Auth, poolCount int) (*smtpSender, error) {
	sender := &smtpSender{from: from, address: addr, auth: auth}
	if poolCount > 0 {
		pool, err := email.NewPool(addr, poolCount, auth)
		if err != nil {
			return nil, err
		}
		sender.pool = pool
	}
	return sender, nil
}

type smtpSender struct {
	from    mail.Address
	address string
	auth    smtp.Auth
	pool    *email.Pool
}

// works only with pool
func (s *smtpSender) Close() error {
	if s.pool != nil {
		s.pool.Close() // at this time this methods does not return error
	}
	return nil
}

func (s smtpSender) Provider() string {
	return Provider
}

func (s *smtpSender) Send(e mailer.Envelope) (response string, msgId string, err error) {
	if e.Subject() == "" {
		err = errors.New("missing subject")
		return
	}
	if len(e.PublicRecipients()) == 0 {
		err = errors.New("no recipients set")
		return
	}
	if e.PlainBody() == nil && e.HTMLBody() == nil {
		err = errors.New("e-mail has no body")
		return
	}

	defer func() {
		if err := e.Close(); err != nil {
			log.Println(err)
		}
	}()

	m := email.NewEmail()

	m.Subject = e.Subject()

	if e.From() != nil {
		m.From = e.From().String()
	} else {
		m.From = s.from.String()
	}

	if e.Sender() != nil {
		m.Sender = e.Sender().String()
	}

	m.To = make([]string, len(e.PublicRecipients()))
	for k, v := range e.PublicRecipients() {
		m.To[k] = v.String()
	}

	m.Bcc = make([]string, len(e.PrivateRecipients()))
	for k, v := range e.PrivateRecipients() {
		m.Bcc[k] = v.String()
	}

	if e.ReplyTo() != nil {
		m.ReplyTo = []string{e.ReplyTo().String()}
	}

	if r := e.PlainBody(); r != nil {
		body, err := ioutil.ReadAll(r)
		if err != nil {
			return "", "", err
		}
		m.Text = body
	}

	if r := e.HTMLBody(); r != nil {
		body, err := ioutil.ReadAll(r)
		if err != nil {
			return "", "", err
		}
		m.HTML = body
	}

	m.Attachments = make([]*email.Attachment, 0, len(e.Attachments())+len(e.Assets()))
	for _, v := range e.Attachments() {
		att, err := m.Attach(v.Data(), v.FileName(), v.MimeType())
		if err != nil {
			return "", "", err
		}
		m.Attachments = append(m.Attachments, att)
	}
	for _, v := range e.Assets() {
		att, err := m.Attach(v.Data(), v.FileName(), v.MimeType())
		if err != nil {
			return "", "", err
		}
		att.HTMLRelated = true
		m.Attachments = append(m.Attachments, att)
	}

	for h, vals := range e.Headers() {
		for _, v := range vals {
			m.Headers.Add(h, v)
		}
	}

	if s.pool != nil {
		err = s.pool.Send(m, time.Second*5)
	} else {
		err = m.Send(s.address, s.auth)
	}

	return
}
