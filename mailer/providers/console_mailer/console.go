package console_mailer

import (
	"bytes"
	"github.com/ivanjaros/ijlibs/mailer"
	"github.com/ivanjaros/ijlibs/random"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"strings"
)

const Provider = "console"

// error is always nil but we implement in in case something changes in the future
func New(from mail.Address) (*consoleSender, error) {
	return &consoleSender{from}, nil
}

type consoleSender struct {
	from mail.Address
}

func (s consoleSender) Provider() string {
	return Provider
}

func (s consoleSender) Close() error {
	return nil
}

func (s *consoleSender) Send(envelope mailer.Envelope) (response string, msgId string, err error) {
	var e YamlEnvelope

	defer envelope.Close()

	cc := envelope.PublicRecipients()
	for k := range cc {
		e.CC = append(e.CC, cc[k].String())
	}

	bcc := envelope.PrivateRecipients()
	for k := range bcc {
		e.BCC = append(e.BCC, bcc[k].String())
	}

	if r := envelope.ReplyTo(); r != nil {
		e.ReplyTo = r.String()
	}

	if s := envelope.Sender(); s != nil {
		e.Sender = s.String()
	}

	e.From = s.from.String()
	if f := envelope.From(); f != nil {
		e.From = f.String()
	}

	e.Subject = envelope.Subject()

	if att := envelope.Attachments(); len(att) > 0 {
		e.Attachments = make([]string, len(att))
		for k, a := range att {
			data, _ := ioutil.ReadAll(a.Data())
			e.Attachments[k] = a.FileName() + "\n" + a.MimeType() + "\n" + string(data)
		}
	}

	if ass := envelope.Assets(); len(ass) > 0 {
		e.Attachments = make([]string, len(ass))
		for k, a := range ass {
			data, _ := ioutil.ReadAll(a.Data())
			e.Assets[k] = a.FileName() + "\n" + a.MimeType() + "\n" + string(data)
		}
	}

	e.Headers = envelope.Headers()

	var sb strings.Builder
	if _, err = io.Copy(&sb, envelope.PlainBody()); err == nil {
		e.PlainBody = sb.String()
	} else {
		return "", "", err
	}

	sb.Reset() // reuse
	if _, err = io.Copy(&sb, envelope.HTMLBody()); err == nil {
		e.HtmlBody = sb.String()
	} else {
		return "", "", err
	}

	buff := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buff).Encode(&e); err != nil {
		return "", "", err
	}

	log.Println("----------------------")
	log.Println("e-mail message preview")
	log.Println("----------------------")
	log.Println(buff.String())

	return "console mailer is only for testing and will not send out any e-mails", random.String(24, random.AlphaSetLC, random.AlphaSetUC, random.NumSet), nil
}

type YamlEnvelope struct {
	CC          []string            `yaml:"cc"`
	BCC         []string            `yaml:"bcc"`
	ReplyTo     string              `yaml:"reply-to"`
	Sender      string              `yaml:"sender"`
	From        string              `yaml:"from"`
	Subject     string              `yaml:"subject"`
	PlainBody   string              `yaml:"plain_body"`
	HtmlBody    string              `yaml:"html_body"`
	Attachments []string            `yaml:"attachments"`
	Assets      []string            `yaml:"assets"`
	Headers     map[string][]string `yaml:"headers"`
}
