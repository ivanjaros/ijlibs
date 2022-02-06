package epp

import (
	"github.com/ivanjaros/jslibs/epp/utils"
)

// This is probably the RFC to use for reference: https://tools.ietf.org/html/rfc4933

const (
	CONTACT_STATUS_OK = "ok"

	CONSTACT_DISCLOSE_NAME         = "name"
	CONSTACT_DISCLOSE_ORGANIZATION = "org"
	CONSTACT_DISCLOSE_ADDRESS      = "addr"
	CONSTACT_DISCLOSE_TELEPHONE    = "voice"
	CONSTACT_DISCLOSE_FAX          = "fax"
	CONSTACT_DISCLOSE_EMAIL        = "email"
)

type ContactObject struct {
	ContactID  string                  `xml:"id"`
	PostalInfo ContactPostalInfoObject `xml:"postalInfo"`
	Voice      string                  `xml:"voice"`
	Fax        string                  `xml:"fax"`
	Email      string                  `xml:"email"`
	AuthInfo   ContactAuthObject       `xml:"authInfo"`
	Disclose   DiscloseContactObject   `xml:"disclose"`
}

type ContactInfoDataObject struct {
	ContactObject
	StorageID   string                 `xml:"roid"`
	Statuses    []InfoDataStatusObject `xml:"status"`
	Owner       string                 `xml:"clid"`
	Creator     string                 `xml:"crid"`
	CreatedAt   string                 `xml:"crDate,omitempty"`
	UpdatedBy   string                 `xml:"upID,omitempty"`
	UpdatedDate string                 `xml:"upDate,omitempty"`
}

type ContactAuthObject struct {
	Password string `xml:"pw,omitempty"`
	Other    string `xml:"ext,omitempty"`
}

type ContactPostalInfoObject struct {
	// omitempty is used only for responses via ContactInfoDataObject
	Type         string                     `xml:"type,attr,omitempty"`
	Name         string                     `xml:"name"`
	Organization string                     `xml:"org"`
	Address      ContactPostalAddressObject `xml:"addr"`
}

type ContactPostalAddressObject struct {
	Street      string `xml:"street"`
	City        string `xml:"city"`
	StateCode   string `xml:"sp"`
	PostalCode  string `xml:"pc"`
	CountryCode string `xml:"cc"`
}

type DiscloseContactObject struct {
	Flag         int       `xml:"flag,attr"`
	Name         *struct{} `xml:"name,omitempty"`
	Organization *struct{} `xml:"org,omitempty"`
	Address      *struct{} `xml:"addr,omitempty"`
	Voice        *struct{} `xml:"voice,omitempty"`
	Fax          *struct{} `xml:"fax,omitempty"`
	Email        *struct{} `xml:"email,omitempty"`
}

func (obj *DiscloseContactObject) GetFields() []string {
	fields := []string{}
	if obj.Name != nil {
		fields = append(fields, CONSTACT_DISCLOSE_NAME)
	}
	if obj.Organization != nil {
		fields = append(fields, CONSTACT_DISCLOSE_ORGANIZATION)
	}
	if obj.Address != nil {
		fields = append(fields, CONSTACT_DISCLOSE_ADDRESS)
	}
	if obj.Voice != nil {
		fields = append(fields, CONSTACT_DISCLOSE_TELEPHONE)
	}
	if obj.Fax != nil {
		fields = append(fields, CONSTACT_DISCLOSE_FAX)
	}
	if obj.Email != nil {
		fields = append(fields, CONSTACT_DISCLOSE_EMAIL)
	}
	return fields
}

func (obj *DiscloseContactObject) SetFields(fields ...string) {
	value := struct{}{}
	for _, field := range fields {
		switch field {
		case CONSTACT_DISCLOSE_NAME:
			if obj.Flag == 1 {
				obj.Name = &value
			} else {
				obj.Name = nil
			}
		case CONSTACT_DISCLOSE_ORGANIZATION:
			if obj.Flag == 1 {
				obj.Organization = &value
			} else {
				obj.Organization = nil
			}
		case CONSTACT_DISCLOSE_ADDRESS:
			if obj.Flag == 1 {
				obj.Address = &value
			} else {
				obj.Address = nil
			}
		case CONSTACT_DISCLOSE_TELEPHONE:
			if obj.Flag == 1 {
				obj.Voice = &value
			} else {
				obj.Voice = nil
			}
		case CONSTACT_DISCLOSE_FAX:
			if obj.Flag == 1 {
				obj.Fax = &value
			} else {
				obj.Fax = nil
			}
		case CONSTACT_DISCLOSE_EMAIL:
			if obj.Flag == 1 {
				obj.Email = &value
			} else {
				obj.Email = nil
			}
		}
	}
}

func (obj *ContactObject) Validate() (ErrorCode, error) {
	if obj == nil {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoObject
	}

	if utils.LengthRange(obj.ContactID, 3, 16) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrContactID
	}

	if utils.InArray(obj.PostalInfo.Type, []string{"int", "loc"}) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrPostalType
	}

	if utils.LengthRange(obj.PostalInfo.Name, 1, 255) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrPostalName
	}

	if utils.LengthRange(obj.PostalInfo.Organization, 0, 255) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrPostalOrg
	}

	if utils.LengthRange(obj.PostalInfo.Address.Street, 1, 255) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrPostalStreet
	}

	if utils.LengthRange(obj.PostalInfo.Address.City, 1, 255) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrPostalCity
	}

	if utils.LengthRange(obj.PostalInfo.Address.StateCode, 0, 255) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrPostalState
	}

	if utils.LengthRange(obj.PostalInfo.Address.PostalCode, 1, 16) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrPostalCode
	}

	if len(obj.PostalInfo.Address.CountryCode) != 2 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrPostalCountry
	}

	if PhoneRegex.MatchString(obj.Voice) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrVoice
	}

	if len(obj.Fax) > 0 && PhoneRegex.MatchString(obj.Fax) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrFax
	}

	// the spec does not validate the pattern(like presence of @), just length(1+)
	if len(obj.Email) < 1 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrEmail
	}

	if obj.AuthInfo.Password == "" && obj.AuthInfo.Other == "" {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoAuthInfo
	}

	if obj.AuthInfo.Password != "" && obj.AuthInfo.Other != "" {
		return STATUS_ERR_COMMAND_SYNTAX, ErrTooManyAuth
	}

	return STATUS_OK, nil
}

type GransyContactObject struct {
	IDNum     *string `xml:"idnum,omitempty"`
	VAT       *string `xml:"vat,omitempty"`
	Birthdate *string `xml:"birthdate,omitempty"`
}

func (e *GransyContactObject) Validate() (ErrorCode, error) {
	if e.IDNum != nil && *e.IDNum != "" && utils.LengthRange(*e.IDNum, 5, 20) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrInvalidExtension
	}
	if e.VAT != nil && *e.VAT != "" && utils.LengthRange(*e.VAT, 5, 20) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrInvalidExtension
	}
	if e.Birthdate != nil && *e.Birthdate != "" {
		if _, err := DateFromString(*e.Birthdate); err != nil {
			return STATUS_ERR_COMMAND_SYNTAX, ErrInvalidExtension
		}
	}
	return STATUS_OK, nil
}

func ContactDiscloseFields() []string {
	return []string{
		CONSTACT_DISCLOSE_NAME,
		CONSTACT_DISCLOSE_ORGANIZATION,
		CONSTACT_DISCLOSE_ADDRESS,
		CONSTACT_DISCLOSE_TELEPHONE,
		CONSTACT_DISCLOSE_FAX,
		CONSTACT_DISCLOSE_EMAIL,
	}
}
