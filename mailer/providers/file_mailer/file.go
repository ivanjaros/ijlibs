package file_mailer

import (
	"bytes"
	"errors"
	"github.com/ivanjaros/ijlibs/mailer"
	"github.com/ivanjaros/ijlibs/random"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/mail"
	"os"
	"strings"
)

const Provider = "file"

func New(from mail.Address, directory string) (*fileSender, error) {
	if FileExists(directory) == false {
		return nil, errors.New("directory does not exist")
	}
	return &fileSender{from, directory}, nil
}

type fileSender struct {
	from mail.Address
	dir  string
}

func (s fileSender) Provider() string {
	return Provider
}

func (s fileSender) Close() error {
	return nil
}

func (s *fileSender) Send(envelope mailer.Envelope) (response string, msgId string, err error) {
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

	id := random.String(24, random.AlphaSetLC, random.AlphaSetUC, random.NumSet)
	if err := ioutil.WriteFile(s.dir+"/"+id+".yaml", buff.Bytes(), 0644); err != nil {
		return err.Error(), "", err
	}

	return "", id, nil
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

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
