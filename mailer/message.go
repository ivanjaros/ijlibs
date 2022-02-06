package mailer

import (
	"bytes"
	"errors"
	"io"
	"net/mail"
	"strings"
)

// we have no reason to use private fields but the names collide with methods
// and figuring out some new names can result in bad DX so we'll just fall back to
// setters.
type Message struct {
	cc          []mail.Address
	bcc         []mail.Address
	replyTo     *mail.Address
	sender      *mail.Address
	from        *mail.Address
	subject     string
	plainBody   *nopCloser
	htmlBody    *nopCloser
	attachments []Asset
	assets      []Asset
	headers     map[string][]string
}

func (msg *Message) PublicRecipients() []mail.Address {
	return msg.cc
}

func (msg *Message) PrivateRecipients() []mail.Address {
	return msg.bcc
}

func (msg *Message) ReplyTo() *mail.Address {
	return msg.replyTo
}

func (msg *Message) Sender() *mail.Address {
	return msg.sender
}

func (msg *Message) From() *mail.Address {
	return msg.from
}

func (msg *Message) Subject() string {
	return msg.subject
}

func (msg *Message) PlainBody() io.ReadCloser {
	// this is required due to https://github.com/golang/go/issues/42663
	if msg.plainBody == nil {
		return &nopCloser{bytes.NewBuffer(nil)}
	}
	return msg.plainBody
}

func (msg *Message) HTMLBody() io.ReadCloser {
	// this is required due to https://github.com/golang/go/issues/42663
	if msg.htmlBody == nil {
		return &nopCloser{bytes.NewBuffer(nil)}
	}
	return msg.htmlBody
}

func (msg *Message) Attachments() []Asset {
	return msg.attachments
}

func (msg *Message) Assets() []Asset {
	return msg.assets
}

func (msg *Message) Headers() mail.Header {
	return msg.headers
}

func (msg *Message) Close() error {
	errs := make([]string, 0, len(msg.assets)+len(msg.attachments))
	for k := range msg.assets {
		if err := msg.assets[k].Data().Close(); err != nil {
			errs = append(errs, msg.assets[k].FileName()+": "+err.Error())
		}
	}
	for k := range msg.attachments {
		if err := msg.attachments[k].Data().Close(); err != nil {
			errs = append(errs, msg.assets[k].FileName()+": "+err.Error())
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.New("failed to close envelope:\n" + strings.Join(errs, "\n- "))
}

func (msg *Message) SetTo(to []mail.Address) {
	msg.cc = to
}

func (msg *Message) SetHiddenRecipients(to []mail.Address) {
	msg.bcc = to
}

func (msg *Message) SetReplyTo(to mail.Address) {
	msg.replyTo = &to
}

func (msg *Message) SetSender(s mail.Address) {
	msg.sender = &s
}

func (msg *Message) SetFrom(f mail.Address) {
	msg.from = &f
}

func (msg *Message) SetSubject(s string) {
	msg.subject = s
}

func (msg *Message) SetBody(plainBody []byte) {
	msg.plainBody = &nopCloser{bytes.NewBuffer(plainBody)}
}

// this method is useful when used with templates
func (msg *Message) GetTextWriter() io.Writer {
	if msg.plainBody == nil {
		msg.plainBody = &nopCloser{bytes.NewBuffer(nil)}
	}
	return msg.plainBody
}

func (msg *Message) SetHtml(html []byte) {
	msg.htmlBody = &nopCloser{bytes.NewBuffer(html)}
}

// this method is useful when used with templates
func (msg *Message) GetHtmlWriter() io.Writer {
	if msg.htmlBody == nil {
		msg.htmlBody = &nopCloser{bytes.NewBuffer(nil)}
	}
	return msg.htmlBody
}

func (msg *Message) SetAttachments(a []Asset) {
	msg.attachments = a
}

func (msg *Message) SetAssets(a []Asset) {
	msg.assets = a
}

func (msg *Message) SetHeaders(h map[string][]string) {
	msg.headers = h
}

type nopCloser struct {
	*bytes.Buffer
}

func (r *nopCloser) Close() error {
	return nil
}
