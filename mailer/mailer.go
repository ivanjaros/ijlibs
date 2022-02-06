package mailer

import (
	"errors"
	"io"
	"net/mail"
)

type Sender interface {
	// Send the e-mail and returns response message, message id, and error.
	// The response and id may be empty, do not rely on these information to be available, even on success.
	// Success is indicated by nil error.
	Send(envelope Envelope) (response string, msgId string, err error)
	// Identifies the underlying provider of the transport(ie. mailgun, sendgrid, amazon_ses,...)
	Provider() string
	// Signals to the sender that it will no longer be needed.
	Close() error
}

type Envelope interface {
	// CC, required
	PublicRecipients() []mail.Address
	// BCC, optional
	PrivateRecipients() []mail.Address
	// Optional, use in case when sender is a dummy address and recipient shouldn't reply to it directly.
	ReplyTo() *mail.Address
	// Optional, each mailer will apply "sender" according the actual e-mail address it uses to send the e-mails.
	// The message author can fill this but there is no guarantee it will be taken into account.
	Sender() *mail.Address
	// This will result in "on behalf of", in case the sender is on a different domain.
	// If no value is provided, the mailer will use its actual e-mail address it uses to send the e-mails.
	From() *mail.Address
	Subject() string
	// When extracting string, use strings.Builder and io.Copy. See the console mailer.
	PlainBody() io.ReadCloser
	// When extracting string, use strings.Builder and io.Copy. See the console mailer.
	HTMLBody() io.ReadCloser
	// Standalone attachments should be used to attach documents to e-mails.
	Attachments() []Asset
	// Inline attachments can be used for attaching images directly into e-mail body.
	Assets() []Asset
	Headers() mail.Header
	// When the sender is done processing the data, it calls this method to signalize
	// that the body, html body and attachments can be closed and won't be used again.
	// Additionally, each individual body or asset can be closed before or after if needed.
	Close() error
}

type Asset interface {
	FileName() string
	MimeType() string
	Data() io.ReadCloser
}

// some basic validation
func Validate(in Envelope) error {
	if len(in.PublicRecipients()) == 0 {
		return errors.New("no recipient provided")
	} else {
		for _, rec := range in.PublicRecipients() {
			if _, err := mail.ParseAddress(rec.String()); err != nil {
				return err
			}
		}
	}

	if len(in.PrivateRecipients()) > 0 {
		for _, rec := range in.PrivateRecipients() {
			if _, err := mail.ParseAddress(rec.String()); err != nil {
				return err
			}
		}
	}

	if rpl := in.ReplyTo(); rpl != nil {
		if _, err := mail.ParseAddress(rpl.String()); err != nil {
			return err
		}
	}

	if snd := in.Sender(); snd != nil {
		if _, err := mail.ParseAddress(snd.String()); err != nil {
			return err
		}
	}

	if frm := in.From(); frm != nil {
		if _, err := mail.ParseAddress(frm.String()); err != nil {
			return err
		}
	}

	// there is no actual maximum limit but after 78 characters there some consequences
	if ln := len(in.Subject()); ln == 0 || ln > 78 {
		return errors.New("invalid subject length")
	}

	if in.PlainBody() == nil && in.HTMLBody() == nil {
		return errors.New("empty content")
	}

	if att := in.Attachments(); len(att) > 0 {
		for _, a := range att {
			if a == nil {
				return errors.New("missing attachment")
			}
		}
	}

	if ass := in.Assets(); len(ass) > 0 {
		for _, a := range ass {
			if a == nil {
				return errors.New("missing attachment")
			}
		}
	}

	return nil
}
