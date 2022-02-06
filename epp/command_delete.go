package epp

import (
	"github.com/ivanjaros/jslibs/epp/utils"
)

type DeleteCommand struct {
	Contact *DeleteContactObject `xml:"urn:ietf:params:xml:ns:contact-1.0 delete,omitempty"`
	Domain  *DeleteDomainObject  `xml:"urn:ietf:params:xml:ns:domain-1.0 delete,omitempty"`
	Host    *DeleteHostObject    `xml:"urn:ietf:params:xml:ns:host-1.0 delete,omitempty"`
}

type DeleteContactObject struct {
	ID string `xml:"id"`
}

type DeleteDomainObject struct {
	Name string `xml:"name"`
}

type DeleteHostObject struct {
	Name string `xml:"name"`
}

func (cmd *DeleteCommand) Validate() (ErrorCode, error) {
	count := 0

	if cmd.Contact != nil {
		count++
		if cmd.Contact.ID == "" {
			return STATUS_ERR_MISSING_PARAM, ErrContactID
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

	if count == 0 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoCommand
	}

	if count > 1 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrTooManyCommands
	}

	return STATUS_OK, nil
}
