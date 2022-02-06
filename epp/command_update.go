package epp

import (
	"github.com/ivanjaros/jslibs/epp/utils"
	"github.com/pkg/errors"
	"net"
)

type UpdateCommand struct {
	Contact *UpdateContactObject `xml:"urn:ietf:params:xml:ns:contact-1.0 update,omitempty"`
	Domain  *UpdateDomainObject  `xml:"urn:ietf:params:xml:ns:domain-1.0 update,omitempty"`
	Host    *UpdateHostObject    `xml:"urn:ietf:params:xml:ns:host-1.0 update,omitempty"`
}

type UpdateContactObject struct {
	ContactID    string                     `xml:"id"`
	AddAction    *ContactChangeActionObject `xml:"add,omitempty"`
	ChangeAction *ContactChangeActionObject `xml:"chg,omitempty"`
	RemoveAction *ContactChangeActionObject `xml:"rem,omitempty"`
}

type ContactChangeActionObject struct {
	PostalInfo *ContactChangePostalInfoObject `xml:"postalInfo,omitempty"`
	Voice      *string                        `xml:"voice,omitempty"`
	Fax        *string                        `xml:"fax,omitempty"`
	Email      *string                        `xml:"email,omitempty"`
	AuthInfo   *ContactAuthObject             `xml:"authInfo,omitempty"`
	Disclose   *DiscloseContactObject         `xml:"disclose,omitempty"`
	Statuses   []InfoDataStatusObject         `xml:"status"`
}

type ContactChangePostalInfoObject struct {
	Type         *string                           `xml:"type,attr,omitempty"`
	Name         *string                           `xml:"name,omitempty"`
	Organization *string                           `xml:"org,omitempty"`
	Address      *ContactChangePostalAddressObject `xml:"addr,omitempty"`
}

type ContactChangePostalAddressObject struct {
	Street      *string `xml:"street,omitempty"`
	City        *string `xml:"city,omitempty"`
	StateCode   *string `xml:"sp,omitempty"`
	PostalCode  *string `xml:"pc,omitempty"`
	CountryCode *string `xml:"cc,omitempty"`
}

type UpdateDomainObject struct {
	Name         string              `xml:"name"`
	AddAction    *DomainUpdateObject `xml:"add,omitempty"`
	ChangeAction *DomainUpdateObject `xml:"chg,omitempty"`
	RemoveAction *DomainUpdateObject `xml:"rem,omitempty"`
}

type DomainUpdateObject struct {
	NameServers []DomainNameServerObject `xml:"ns"`
	Contacts    []DomainContactObject    `xml:"contact"`
	Statuses    []InfoDataStatusObject   `xml:"status"`
	Registrant  *string                  `xml:"registrant,omitempty"`
	AuthInfo    *DomainAuthObject        `xml:"authInfo,omitempty"`
}

type UpdateHostObject struct {
	Name         string            `xml:"name"`
	AddAction    *HostUpdateObject `xml:"add,omitempty"`
	RemoveAction *HostUpdateObject `xml:"rem,omitempty"`
	ChangeAction *HostUpdateObject `xml:"chg,omitempty"`
}

type HostUpdateObject struct {
	Name     *string                `xml:"name"`
	IPs      []string               `xml:"addr"`
	Statuses []InfoDataStatusObject `xml:"status"`
}

func (cmd *UpdateCommand) Validate() (ErrorCode, error) {
	count := 0

	if cmd.Contact != nil {
		count++

		code, err := cmd.Contact.Validate()
		if err != nil {
			return code, err
		}
	}

	if cmd.Host != nil {
		count++

		code, err := cmd.Host.Validate()
		if err != nil {
			return code, err
		}
	}

	if cmd.Domain != nil {
		count++

		code, err := cmd.Domain.Validate()
		if err != nil {
			return code, err
		}
	}

	if count == 0 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoObject
	}

	if count == 1 {
		return STATUS_OK, nil
	}

	return STATUS_ERR_COMMAND_SYNTAX, ErrTooManyObjects
}

func (obj *UpdateContactObject) Validate() (ErrorCode, error) {
	if obj.ChangeAction != nil {
		if code, err := obj.ChangeAction.Validate(); err != nil {
			return code, err
		}
	}

	if obj.AddAction != nil {
		if code, err := obj.AddAction.Validate(); err != nil {
			return code, err
		}
	}

	if obj.RemoveAction != nil {
		if code, err := obj.RemoveAction.Validate(); err != nil {
			return code, err
		}
	}

	changeStatuses := []string{}

	if obj.AddAction != nil && obj.AddAction.Statuses != nil {
		for _, s := range obj.AddAction.Statuses {
			changeStatuses = append(changeStatuses, s.Status)
		}
	}

	if obj.RemoveAction != nil && obj.RemoveAction.Statuses != nil {
		for _, s := range obj.RemoveAction.Statuses {
			changeStatuses = append(changeStatuses, s.Status)
		}
	}

	allowedStatuses := []string{
		STATUS_TYPE_CLIENT_DELETE_PROHIBITED,
		STATUS_TYPE_CLIENT_TRANSFER_PROHIBITED,
		STATUS_TYPE_CLIENT_UPDATE_PROHIBITED,
	}

	if diff := utils.ArrayDiff(changeStatuses, allowedStatuses); len(diff) > 0 {
		return STATUS_ERR_PARAM_RANGE, errors.Errorf("Cannot use status %s in update", diff[0])
	}

	return STATUS_OK, nil
}

func (obj *ContactChangeActionObject) Validate() (ErrorCode, error) {
	if obj.Email != nil {
		// the spec does not validate the pattern(like presence of @), just length(1+)
		if len(*obj.Email) < 1 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrEmail
		}
	}

	if obj.Voice != nil {
		if PhoneRegex.MatchString(*obj.Voice) == false || len(*obj.Voice) > 17 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrVoice
		}
	}

	if obj.Fax != nil {
		if PhoneRegex.MatchString(*obj.Fax) == false || len(*obj.Fax) > 17 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrFax
		}
	}

	if obj.AuthInfo != nil {
		if obj.AuthInfo.Other != "" {
			return STATUS_ERR_COMMAND_SYNTAX, ErrAuthType
		}
	}

	if obj.PostalInfo != nil {
		if obj.PostalInfo.Type != nil {
			if utils.InArray(*obj.PostalInfo.Type, []string{"int", "loc"}) == false {
				return STATUS_ERR_COMMAND_SYNTAX, ErrPostalType
			}
		}

		if obj.PostalInfo.Name != nil {
			if utils.LengthRange(*obj.PostalInfo.Name, 1, 255) == false {
				return STATUS_ERR_COMMAND_SYNTAX, ErrPostalName
			}
		}

		if obj.PostalInfo.Organization != nil {
			if utils.LengthRange(*obj.PostalInfo.Organization, 0, 255) == false {
				return STATUS_ERR_COMMAND_SYNTAX, ErrPostalOrg
			}
		}

		if obj.PostalInfo.Address != nil {
			if obj.PostalInfo.Address.Street != nil {
				if utils.LengthRange(*obj.PostalInfo.Address.Street, 1, 255) == false {
					return STATUS_ERR_COMMAND_SYNTAX, ErrPostalStreet
				}
			}

			if obj.PostalInfo.Address.City != nil {
				if utils.LengthRange(*obj.PostalInfo.Address.City, 1, 255) == false {
					return STATUS_ERR_COMMAND_SYNTAX, ErrPostalCity
				}
			}

			if obj.PostalInfo.Address.StateCode != nil {
				if utils.LengthRange(*obj.PostalInfo.Address.StateCode, 0, 255) == false {
					return STATUS_ERR_COMMAND_SYNTAX, ErrPostalState
				}
			}

			if obj.PostalInfo.Address.PostalCode != nil {
				if utils.LengthRange(*obj.PostalInfo.Address.PostalCode, 1, 16) == false {
					return STATUS_ERR_COMMAND_SYNTAX, ErrPostalCode
				}
			}

			if obj.PostalInfo.Address.CountryCode != nil {
				if len(*obj.PostalInfo.Address.CountryCode) != 2 {
					return STATUS_ERR_COMMAND_SYNTAX, ErrPostalCountry
				}
			}
		}
	}

	return STATUS_OK, nil
}

func (obj *UpdateHostObject) Validate() (ErrorCode, error) {
	if utils.LengthRange(obj.Name, 1, 255) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNSName
	}

	if obj.AddAction != nil {
		if code, err := obj.AddAction.Validate(); err != nil {
			return code, err
		}
	}

	if obj.ChangeAction != nil {
		if code, err := obj.ChangeAction.Validate(); err != nil {
			return code, err
		}
	}

	if obj.RemoveAction != nil {
		if code, err := obj.RemoveAction.Validate(); err != nil {
			return code, err
		}
	}

	return STATUS_OK, nil
}

func (obj *HostUpdateObject) Validate() (ErrorCode, error) {
	if obj.Name != nil && utils.LengthRange(*obj.Name, 1, 255) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNSName
	}

	if len(obj.IPs) > 0 {
		for _, ip := range obj.IPs {
			if net.ParseIP(ip) == nil {
				return STATUS_ERR_COMMAND_SYNTAX, ErrIP
			}
		}

		if utils.NumRange(len(obj.IPs), 0, 10) == false {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNSAddrCount
		}

		if len(obj.IPs) != len(utils.ArrayUnique(obj.IPs)) {
			return STATUS_ERR_COMMAND_SYNTAX, ErrIPDuplicate
		}
	}

	return STATUS_OK, nil
}

func (obj *UpdateDomainObject) Validate() (ErrorCode, error) {
	if DomainRegex.MatchString(obj.Name) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDomainName
	}

	if obj.AddAction != nil {
		if code, err := obj.AddAction.Validate(); err != nil {
			return code, err
		}
	}

	if obj.ChangeAction != nil {
		if code, err := obj.ChangeAction.Validate(); err != nil {
			return code, err
		}
	}

	if obj.RemoveAction != nil {
		if code, err := obj.RemoveAction.Validate(); err != nil {
			return code, err
		}
	}

	return STATUS_OK, nil
}

func (obj *DomainUpdateObject) Validate() (ErrorCode, error) {
	if len(obj.NameServers) > 0 {
		for _, ns := range obj.NameServers {
			// standard validation
			if code, err := ns.Validate(); err != nil {
				return code, err
			}

			// we do not support host attributes
			if len(ns.Host) > 0 {
				return STATUS_ERR_COMMAND_SYNTAX, ErrNsForbiddenAttr
			}

			if len(ns.Object) > 8 {
				return STATUS_ERR_COMMAND_SYNTAX, ErrNSMany
			}
		}
	}

	if len(obj.Contacts) > 0 {
		for _, contact := range obj.Contacts {
			if utils.InArray(contact.Type, DomainContactTypes()) == false {
				return STATUS_ERR_COMMAND_SYNTAX, ErrContactType
			}

			if contact.Value == "" {
				return STATUS_ERR_COMMAND_SYNTAX, ErrNoContactInfo
			}
		}
	}

	if len(obj.Statuses) > 0 {
		allowedStatuses := []string{
			STATUS_TYPE_CLIENT_HOLD,
			STATUS_TYPE_CLIENT_DELETE_PROHIBITED,
			STATUS_TYPE_CLIENT_TRANSFER_PROHIBITED,
			STATUS_TYPE_CLIENT_UPDATE_PROHIBITED,
			STATUS_TYPE_CLIENT_RENEW_PROHIBITED,
		}
		statuses := []string{}
		for _, status := range obj.Statuses {
			statuses = append(statuses, status.Status)
		}
		if len(utils.ArrayDiff(statuses, allowedStatuses)) > 0 {
			return STATUS_ERR_PARAM_RANGE, ErrStatus
		}
	}

	if obj.AuthInfo != nil {
		if obj.AuthInfo.Password == "" && obj.AuthInfo.Other == "" {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNoAuthInfo
		}

		if obj.AuthInfo.Password != "" && obj.AuthInfo.Other != "" {
			return STATUS_ERR_COMMAND_SYNTAX, ErrTooManyAuth
		}

		// only password is supported
		if obj.AuthInfo.Password == "" {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNoAuthInfo
		}
	}

	if obj.Registrant != nil && utils.LengthRange(*obj.Registrant, 1, 16) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrRegistrant
	}

	return STATUS_OK, nil
}
