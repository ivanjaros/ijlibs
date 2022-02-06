package epp

import (
	"github.com/ivanjaros/jslibs/epp/utils"
)

const (
	EPP_TRANSFER_OP_APPROVE              = "approve"
	EPP_TRANSFER_OP_CANCEL               = "cancel"
	EPP_TRANSFER_OP_QUERY                = "query"
	EPP_TRANSFER_OP_REJECT               = "reject"
	EPP_TRANSFER_OP_REQUEST              = "request"
	EPP_TRANSFER_STATUS_PENDING          = "pending"
	EPP_TRANSFER_STATUS_CLIENT_APPROVED  = "clientApproved"
	EPP_TRANSFER_STATUS_CLIENT_CANCELLED = "clientCancelled"
	EPP_TRANSFER_STATUS_CLIENT_REJECTED  = "clientRejected"
	EPP_TRANSFER_STATUS_SERVER_APPROVED  = "serverApproved"
	EPP_TRANSFER_STATUS_SERVER_CANCELLED = "serverCancelled"
)

type TransferCommand struct {
	Operation string                       `xml:"op,attr"`
	Domain    *DomainTransferRequestObject `xml:"urn:ietf:params:xml:ns:domain-1.0 transfer"`
}

type DomainTransferRequestObject struct {
	Name     string           `xml:"name"`
	Period   PeriodObject     `xml:"period"`
	AuthInfo DomainAuthObject `xml:"authInfo"`
}

func (cmd *TransferCommand) Validate() (ErrorCode, error) {
	if utils.InArray(cmd.Operation, []string{
		EPP_TRANSFER_OP_APPROVE,
		EPP_TRANSFER_OP_CANCEL,
		EPP_TRANSFER_OP_QUERY,
		EPP_TRANSFER_OP_REJECT,
		EPP_TRANSFER_OP_REQUEST,
	}) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrOperation
	}

	if cmd.Domain == nil {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoObject
	}

	if utils.LengthRange(cmd.Domain.Name, 1, 255) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDomainName
	}

	// Only "request" operation uses period object so we can ignore it for other operations
	if cmd.Operation == EPP_TRANSFER_OP_REQUEST {
		if utils.InArray(cmd.Domain.Period.Unit, []string{EPP_PERIOD_YEAR, EPP_PERIOD_MONTH}) == false {
			return STATUS_ERR_COMMAND_SYNTAX, ErrPeriodType
		}

		// Gransy takes in only years.
		if cmd.Domain.Period.Unit != EPP_PERIOD_YEAR {
			return STATUS_ERR_COMMAND_SYNTAX, ErrPeriodLength
		}

		// No value is a valid value in this case. For such scenario, the period will be retrieved from the TLD.
		// Therefore we only check if the value is not a negative number, which is an invalid value.
		if cmd.Domain.Period.Value < 0 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrPeriodLength
		}
	}

	if cmd.Domain.AuthInfo.Password == "" && cmd.Domain.AuthInfo.Other == "" {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoAuthInfo
	}

	if cmd.Domain.AuthInfo.Password != "" && cmd.Domain.AuthInfo.Other != "" {
		return STATUS_ERR_COMMAND_SYNTAX, ErrTooManyAuth
	}

	// Gransy takes in only passwords
	if cmd.Domain.AuthInfo.Password == "" {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoAuthInfo
	}

	return STATUS_OK, nil
}
