package epp

type CreateCommand struct {
	Contact *ContactObject `xml:"urn:ietf:params:xml:ns:contact-1.0 create,omitempty"`
	Domain  *DomainObject  `xml:"urn:ietf:params:xml:ns:domain-1.0 create,omitempty"`
	Host    *HostObject    `xml:"urn:ietf:params:xml:ns:host-1.0 create,omitempty"`
}

func (cmd *CreateCommand) Validate() (ErrorCode, error) {
	var objects int
	if cmd.Contact != nil {
		objects++
	}
	if cmd.Domain != nil {
		objects++
	}
	if cmd.Host != nil {
		objects++
	}
	if objects > 1 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrTooManyObjects
	}

	if cmd.Contact != nil {
		return cmd.Contact.Validate()
	}

	if cmd.Domain != nil {
		return cmd.Domain.ValidateCreate()
	}

	if cmd.Host != nil {
		return cmd.Host.ValidateCreate()
	}

	return STATUS_ERR_COMMAND_SYNTAX, ErrNoObject
}
