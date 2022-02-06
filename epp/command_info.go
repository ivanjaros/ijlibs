package epp

import (
	"github.com/ivanjaros/jslibs/epp/utils"
)

type InfoCommand struct {
	Contact   *ContactInfoRequestObject   `xml:"urn:ietf:params:xml:ns:contact-1.0 info,omitempty"`
	Domain    *DomainInfoRequestObject    `xml:"urn:ietf:params:xml:ns:domain-1.0 info,omitempty"`
	Host      *HostInfoRequestObject      `xml:"urn:ietf:params:xml:ns:host-1.0 info,omitempty"`
	Order     *OrderInfoRequestObject     `xml:"http://www.subreg.cz/epp/order-1.0 info,omitempty"`
	Registrar *RegistrarInfoRequestObject `xml:"urn:ietf:params:xml:ns:registrar-info-1.0 info,omitempty"`
}

type ContactInfoRequestObject struct {
	ID string `xml:"id"`
}

type DomainInfoRequestObject struct {
	Name     string           `xml:"name"`
	AuthInfo DomainAuthObject `xml:"authInfo"`
}

type OrderInfoRequestObject struct {
	ID int `xml:"id"`
}

type RegistrarInfoRequestObject struct {
	ID string `xml:"id"`
}

type HostInfoRequestObject struct {
	Name string `xml:"name"`
}

func (cmd *InfoCommand) Validate() (ErrorCode, error) {
	count := 0

	if cmd.Contact != nil {
		count++
		if utils.LengthRange(cmd.Contact.ID, 3, 16) == false {
			return STATUS_ERR_COMMAND_SYNTAX, ErrID
		}
	}

	if cmd.Domain != nil {
		count++
		if utils.LengthRange(cmd.Domain.Name, 1, 255) == false {
			return STATUS_ERR_COMMAND_SYNTAX, ErrDomainName
		}
	}

	if cmd.Host != nil {
		count++
		if utils.LengthRange(cmd.Host.Name, 1, 255) == false {
			return STATUS_ERR_COMMAND_SYNTAX, ErrDomainName
		}
	}

	if cmd.Order != nil {
		count++
		if cmd.Order.ID < 1 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrID
		}
	}

	if cmd.Registrar != nil {
		count++
		if utils.LengthRange(cmd.Registrar.ID, 3, 16) == false {
			return STATUS_ERR_COMMAND_SYNTAX, ErrID
		}
	}

	if count == 0 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoObject
	}

	if count > 1 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrTooManyObjects
	}

	return STATUS_OK, nil
}
