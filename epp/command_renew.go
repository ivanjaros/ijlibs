package epp

import (
	"github.com/ivanjaros/jslibs/epp/utils"
)

type RenewCommand struct {
	Domain RenewDomainObject `xml:"urn:ietf:params:xml:ns:domain-1.0 renew"`
}

type RenewDomainObject struct {
	Name           string       `xml:"name"`
	ExpirationDate string       `xml:"curExpDate"`
	Period         PeriodObject `xml:"period"`
}

func (cmd *RenewCommand) Validate() (ErrorCode, error) {
	if utils.LengthRange(cmd.Domain.Name, 1, 255) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDomainName
	}

	if cmd.Domain.Period.Value != 0 {
		if utils.InArray(cmd.Domain.Period.Unit, []string{EPP_PERIOD_YEAR, EPP_PERIOD_MONTH}) == false {
			return STATUS_ERR_COMMAND_SYNTAX, ErrPeriodType
		}

		// Gransy takes in only years.
		if cmd.Domain.Period.Unit != EPP_PERIOD_YEAR {
			return STATUS_ERR_COMMAND_SYNTAX, ErrPeriodLength
		}

		if cmd.Domain.Period.Value < 1 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrPeriodLength
		}
	}

	if _, err := DateFromString(cmd.Domain.ExpirationDate); err != nil {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDateFormat
	}

	return STATUS_OK, nil
}
