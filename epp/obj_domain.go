package epp

import (
	"github.com/ivanjaros/jslibs/epp/utils"
	"time"
)

const (
	DOMAIN_STATUS_OK                = "ok"
	EPP_PERIOD_YEAR                 = "y"
	EPP_PERIOD_MONTH                = "m"
	EPP_DOMAIN_CONTACT_TYPE_ADMIN   = "admin"
	EPP_DOMAIN_CONTACT_TYPE_TECH    = "tech"
	EPP_DOMAIN_CONTACT_TYPE_BILLING = "billing"
	DOMAIN_RGP_REPORT               = "report"
	DOMAIN_RGP_REQUEST              = "request"

	EXP_POLICY_AUTORENEW = "autorenew"
	EXP_POLICY_RGP       = "rgp"
)

type DomainObject struct {
	Name        string                   `xml:"name"`
	Period      PeriodObject             `xml:"period"`
	NameServers []DomainNameServerObject `xml:"ns"`
	Registrant  string                   `xml:"registrant"`
	Contact     []DomainContactObject    `xml:"contact"`
	AuthInfo    DomainAuthObject         `xml:"authInfo"`
}

type DomainInfoDataObject struct {
	Name           string                 `xml:"name"`
	StorageID      string                 `xml:"roid"`
	Statuses       []InfoDataStatusObject `xml:"status"`
	Owner          string                 `xml:"clid"`
	Creator        string                 `xml:"crid"`
	UpdatedBy      string                 `xml:"upID,omitempty"`
	CreatedDate    string                 `xml:"crDate"`
	UpdatedDate    string                 `xml:"upDate,omitempty"`
	ExpirationDate string                 `xml:"exDate,omitempty"`
	TransferDate   string                 `xml:"trdate,omitempty"`
	RegistrantID   string                 `xml:"registrant,omitempty"`
	Contacts       []DomainContactObject  `xml:"contact,omitempty"`
	NameServers    DomainNameServerObject `xml:"ns"`
	AuthInfo       DomainAuthObject       `xml:"authInfo"`
}

type DomainInfoRGPExtension struct {
	Status DomainRGPStatusExtension `xml:"rgpStatus"`
}

type DomainRGPStatusExtension struct {
	Value string `xml:"s,attr"`
}

type UpdateDomainRGPExtension struct {
	Restore UpdateDomainRGPRestoreExtension `xml:"restore"`
}

type UpdateDomainRGPRestoreExtension struct {
	Operation string                     `xml:"op,attr"`
	Report    *DomainRestoreReportObject `xml:"report,omitempty"`
}

type DomainRestoreReportObject struct {
	PreData       InnerXML                       `xml:"preData"`
	PostData      InnerXML                       `xml:"postData"`
	DeleteTime    string                         `xml:"delTime"`
	RestoreTime   string                         `xml:"resTime"`
	RestoreReason string                         `xml:"resReason"`
	Statement     []DomainRestoreStatementObject `xml:"statement"`
	Other         InnerXML                       `xml:"other"`
}

type DomainRestoreStatementObject struct {
	Language  string `xml:"lang,attr" json:"lang"`
	Statement string `xml:",chardata" json:"statement"`
}

// https://tools.ietf.org/html/rfc3915#section-4.2.5
func (obj *UpdateDomainRGPExtension) Validate() (ErrorCode, error) {
	if utils.InArray(obj.Restore.Operation, []string{DOMAIN_RGP_REPORT, DOMAIN_RGP_REQUEST}) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrRestoreOp
	}

	if obj.Restore.Operation == DOMAIN_RGP_REQUEST && obj.Restore.Report != nil {
		return STATUS_ERR_COMMAND_SYNTAX, ErrRestoreReportPresent
	}

	if obj.Restore.Operation == DOMAIN_RGP_REPORT && obj.Restore.Report == nil {
		return STATUS_ERR_COMMAND_SYNTAX, ErrRestoreReportNotPresent
	}

	if rep := obj.Restore.Report; rep != nil {
		if _, err := time.Parse(time.RFC3339, rep.DeleteTime); err != nil {
			return STATUS_ERR_COMMAND_SYNTAX, ErrRestoreDelTime
		}
		if _, err := time.Parse(time.RFC3339, rep.RestoreTime); err != nil {
			return STATUS_ERR_COMMAND_SYNTAX, ErrRestoreTime
		}
		if len(rep.Statement) != 2 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrRestoreStatement
		}
	}

	return STATUS_OK, nil
}

type PeriodObject struct {
	Value int    `xml:",chardata"`
	Unit  string `xml:"unit,attr"`
}

type DomainAuthObject struct {
	Password string `xml:"pw,omitempty"`
	Other    string `xml:"ext,omitempty"`
}

type DomainNameServerObject struct {
	// hostAttr is never used but we'll leave it here for correct
	// RFC compatibility and future use.
	Host   []DomainNameServerHostObject `xml:"hostAttr,omitempty"`
	Object []string                     `xml:"hostObj"`
}

func (obj *DomainNameServerObject) Validate() (ErrorCode, error) {
	if len(obj.Host) == 0 && len(obj.Object) == 0 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoNS
	}

	if len(obj.Object) > 0 {
		for _, obj := range obj.Object {
			if utils.LengthRange(obj, 1, 255) == false {
				return STATUS_ERR_COMMAND_SYNTAX, ErrNS
			}
		}

		if len(obj.Object) > len(utils.ArrayUnique(obj.Object)) {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNSDuplicate
		}
	}

	if len(obj.Host) > 0 {
		for _, host := range obj.Host {
			if host.Name == "" {
				return STATUS_ERR_COMMAND_SYNTAX, ErrNSNoName
			}
			if HostnameRegex.MatchString(host.Name) == false {
				return STATUS_ERR_COMMAND_SYNTAX, ErrNSName
			}
			if host.Address != "" && Ipv4Regex.MatchString(host.Address) == false {
				return STATUS_ERR_COMMAND_SYNTAX, ErrNSAddress
			}
		}
	}

	return STATUS_OK, nil
}

type DomainNameServerHostObject struct {
	Name    string `xml:"hostName"`
	Address string `xml:"hostAddr"`
}

type DomainContactObject struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

func (obj *DomainObject) ValidateCreate() (ErrorCode, error) {
	if obj == nil {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoObject
	}

	if utils.LengthRange(obj.Name, 1, 255) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDomainName
	}

	if DomainRegex.MatchString(obj.Name) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDomainName
	}

	if utils.InArray(obj.Period.Unit, []string{EPP_PERIOD_YEAR, EPP_PERIOD_MONTH}) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrPeriodType
	}

	// Gransy takes in only years.
	if obj.Period.Unit != EPP_PERIOD_YEAR {
		return STATUS_ERR_COMMAND_SYNTAX, ErrPeriodLength
	}

	// according to spec, this is optional and should equal to server's
	// setting but we require it, probably constraint from previous PHP version.
	if obj.Period.Value < 1 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrPeriodLength
	}

	if utils.LengthRange(obj.Registrant, 1, 16) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrRegistrant
	}

	if obj.AuthInfo.Password == "" && obj.AuthInfo.Other == "" {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoAuthInfo
	}

	if obj.AuthInfo.Password != "" && obj.AuthInfo.Other != "" {
		return STATUS_ERR_COMMAND_SYNTAX, ErrTooManyAuth
	}

	// Only password is supported.
	if obj.AuthInfo.Password == "" {
		return STATUS_ERR_COMMAND_SYNTAX, ErrAuthType
	}

	// according to spec this is optional by we are requiring it.
	// probably constraint from the previous PHP version.
	if len(obj.Contact) > 0 {
		cTypes := map[string]int{
			EPP_DOMAIN_CONTACT_TYPE_ADMIN:   0,
			EPP_DOMAIN_CONTACT_TYPE_TECH:    0,
			EPP_DOMAIN_CONTACT_TYPE_BILLING: 0,
		}
		for _, cType := range obj.Contact {
			if utils.InArray(cType.Type, DomainContactTypes()) == false {
				return STATUS_ERR_COMMAND_SYNTAX, ErrContactType
			}
			if cType.Value == "" {
				return STATUS_ERR_COMMAND_SYNTAX, ErrNoContactInfo
			}
			cTypes[cType.Type]++
			if cTypes[cType.Type] > 1 {
				return STATUS_ERR_COMMAND_SYNTAX, ErrContactTypeMany
			}
		}
	} else {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoContactInfo
	}

	for _, ns := range obj.NameServers {
		// Generic validation.
		if code, err := ns.Validate(); err != nil {
			return code, err
		}

		// While creating a domain, only host objects are supported.
		if len(ns.Host) > 0 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNsForbiddenAttr
		}
		if len(ns.Object) == 0 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNSNoName
		}
		if len(ns.Object) < 2 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNSFew
		}
		if len(ns.Object) > 8 {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNSMany
		}
	}

	return STATUS_OK, nil
}

func ContainsNOOKstatus(stats []string) bool {
	nookEnums := []string{
		STATUS_TYPE_CLIENT_DELETE_PROHIBITED,
		STATUS_TYPE_CLIENT_HOLD,
		STATUS_TYPE_CLIENT_RENEW_PROHIBITED,
		STATUS_TYPE_CLIENT_TRANSFER_PROHIBITED,
		STATUS_TYPE_CLIENT_UPDATE_PROHIBITED,
		STATUS_TYPE_INACTIVE,
		STATUS_TYPE_PENDING_CREATE,
		STATUS_TYPE_PENDING_DELETE,
		STATUS_TYPE_PENDING_RENEW,
		STATUS_TYPE_PENDING_TRANSFER,
		STATUS_TYPE_PENDING_UPDATE,
		STATUS_TYPE_SERVER_DELETE_PROHIBITED,
		STATUS_TYPE_SERVER_HOLD,
		STATUS_TYPE_SERVER_RENEW_PROHIBITED,
		STATUS_TYPE_SERVER_TRANSFER_PROHIBITED,
		STATUS_TYPE_SERVER_UPDATE_PROHIBITED,
	}

	for _, s := range stats {
		if utils.InArray(s, nookEnums) {
			return true
		}
	}

	return false
}

func ValidateDomainStatus(s []string) bool {
	if len(s) > 11 {
		return false
	}

	enums := []string{
		STATUS_TYPE_CLIENT_DELETE_PROHIBITED,
		STATUS_TYPE_CLIENT_HOLD,
		STATUS_TYPE_CLIENT_RENEW_PROHIBITED,
		STATUS_TYPE_CLIENT_TRANSFER_PROHIBITED,
		STATUS_TYPE_CLIENT_UPDATE_PROHIBITED,
		STATUS_TYPE_INACTIVE,
		STATUS_TYPE_OK,
		STATUS_TYPE_PENDING_CREATE,
		STATUS_TYPE_PENDING_DELETE,
		STATUS_TYPE_PENDING_RENEW,
		STATUS_TYPE_PENDING_TRANSFER,
		STATUS_TYPE_PENDING_UPDATE,
		STATUS_TYPE_SERVER_DELETE_PROHIBITED,
		STATUS_TYPE_SERVER_HOLD,
		STATUS_TYPE_SERVER_RENEW_PROHIBITED,
		STATUS_TYPE_SERVER_TRANSFER_PROHIBITED,
		STATUS_TYPE_SERVER_UPDATE_PROHIBITED,
	}

	for _, status := range s {
		if utils.InArray(status, enums) == false {
			return false
		}
	}

	return true
}

func DomainContactTypes() []string {
	return []string{
		EPP_DOMAIN_CONTACT_TYPE_ADMIN,
		EPP_DOMAIN_CONTACT_TYPE_TECH,
		EPP_DOMAIN_CONTACT_TYPE_BILLING,
	}
}

type DomainRenewData struct {
	Name           string `xml:"name"`
	ExpirationDate string `xml:"exDate"`
}

type DomainUpdateRGPExtension struct {
	Status DomainUpdateRGPExtensionStatusObject `xml:"rgpStatus"`
}

type DomainUpdateRGPExtensionStatusObject struct {
	Value string `xml:"s,attr"`
}
