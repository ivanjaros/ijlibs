package epp

import (
	"github.com/ivanjaros/jslibs/epp/utils"
)

type CheckCommand struct {
	Domain  *CheckDomainRequest  `xml:"urn:ietf:params:xml:ns:domain-1.0 check,omitempty"`
	Contact *CheckContactRequest `xml:"urn:ietf:params:xml:ns:contact-1.0 check,omitempty"`
	Host    *CheckHostRequest    `xml:"urn:ietf:params:xml:ns:host-1.0 check,omitempty"`
}

type CheckDomainRequest struct {
	Names []string `xml:"name"`
}

type CheckContactRequest struct {
	IDs []string `xml:"id"`
}

type CheckHostRequest struct {
	Names []string `xml:"name"`
}

func (cmd *CheckCommand) Validate() (ErrorCode, error) {
	types := 0

	if cmd.Domain != nil {
		types++

		if len(cmd.Domain.Names) == 0 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNoObjects
		}

		if len(cmd.Domain.Names) > 50 {
			return STATUS_ERR_PARAM_RANGE, ErrTooManyObjects
		}

		for _, domain := range cmd.Domain.Names {
			if domain == "" {
				return STATUS_ERR_COMMAND_SYNTAX, ErrNoObjectID
			}
			if DomainRegex.MatchString(domain) == false {
				return STATUS_ERR_COMMAND_SYNTAX, ErrDomainName
			}
		}
	}

	if cmd.Contact != nil {
		types++

		if len(cmd.Contact.IDs) == 0 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNoObjects
		}

		if len(cmd.Contact.IDs) > 50 {
			return STATUS_ERR_PARAM_RANGE, ErrTooManyObjects
		}

		for _, id := range cmd.Contact.IDs {
			if id == "" {
				return STATUS_ERR_COMMAND_SYNTAX, ErrNoObjectID
			}
			if utils.LengthRange(id, 3, 16) == false {
				return STATUS_ERR_COMMAND_SYNTAX, ErrContactID
			}
		}
	}

	if cmd.Host != nil {
		types++

		if len(cmd.Host.Names) == 0 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNoObjects
		}

		if len(cmd.Host.Names) > 50 {
			return STATUS_ERR_PARAM_RANGE, ErrTooManyObjects
		}

		for _, name := range cmd.Host.Names {
			if name == "" {
				return STATUS_ERR_COMMAND_SYNTAX, ErrNoObjectID
			}
			if HostnameRegex.MatchString(name) == false {
				return STATUS_ERR_COMMAND_SYNTAX, ErrName
			}
		}
	}

	if types == 0 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoObject
	}

	if types != 1 {
		return STATUS_ERR_INVALID_COMMAND, ErrTooManyCommands
	}

	return STATUS_OK, nil
}
