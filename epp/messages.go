package epp

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type ErrorCode int

func (c ErrorCode) Error() string {
	return GetResultMessage(c)
}

// https://tools.ietf.org/html/rfc5730#section-3
const (
	EPP_PROTO_V1         = "1.0"
	EPP_PROTO_V2         = "2.0"
	EPP_PROTO_LANG_EN    = "en"
	EPP_XML_TAG          = "epp"
	EPP_NAMESPACE        = "urn:ietf:params:xml:ns:epp-1.0"
	EPP_DOMAIN_OBJ_NS    = "urn:ietf:params:xml:ns:domain-1.0"
	EPP_CONTACT_OBJ_NS   = "urn:ietf:params:xml:ns:contact-1.0"
	EPP_HOST_OBJ_NS      = "urn:ietf:params:xml:ns:host-1.0"
	EPP_REGISTRAR_OBJ_NS = "urn:ietf:params:xml:ns:registrar-info-1.0"
	EPP_RGP_OBJ_NS       = "urn:ietf:params:xml:ns:rgp-1.0"
	EPP_DNSSEC_OBJ_NS    = "urn:ietf:params:xml:ns:secDNS-1.1"

	EPP_GRANSY_DOMAIN_OBJ_NS     = "http://www.subreg.cz/epp/gransy-domain-0.1"
	EPP_GRANSY_DOCUMENT_OBJ_NS   = "http://www.subreg.cz/epp/gransy-document-0.1"
	EPP_GRANSY_CONTACT_OBJ_NS    = "http://www.subreg.cz/epp/gransy-contact-0.1"
	EPP_GRANSY_COMMAND_CHECK_EXT = "http://regtonsregistry.cz/commandCheck"

	STATUS_OK                         ErrorCode = 1000 //  "Command completed successfully"
	STATUS_OK_ACTION_PENDING          ErrorCode = 1001 //  "Command completed successfully; action pending"
	STATUS_OK_NO_MESSAGES             ErrorCode = 1300 //  "Command completed successfully; no messages"
	STATUS_OK_ACK_NEEDED              ErrorCode = 1301 //  "Command completed successfully; ack to dequeue"
	STATUS_OK_END_SESSION             ErrorCode = 1500 //  "Command completed successfully; ending Session"
	STATUS_ERR_UNKNOWN_COMMAND        ErrorCode = 2000 //  "Unknown command"
	STATUS_ERR_COMMAND_SYNTAX         ErrorCode = 2001 //  "Command syntax error"
	STATUS_ERR_INVALID_COMMAND        ErrorCode = 2002 //  "Command use error"
	STATUS_ERR_MISSING_PARAM          ErrorCode = 2003 //  "Required parameter missing"
	STATUS_ERR_PARAM_RANGE            ErrorCode = 2004 //  "Parameter value range error"
	STATUS_ERR_PARAM_SYNTAX           ErrorCode = 2005 //  "Parameter value syntax error"
	STATUS_ERR_UNIMPLEMENTED_PROTOCOL ErrorCode = 2100 //  "Unimplemented protocol version"
	STATUS_ERR_UNIMPLEMENTED_COMMAND  ErrorCode = 2101 //  "Unimplemented command"
	STATUS_ERR_UNIMPLEMENTED_OPTION   ErrorCode = 2102 //  "Unimplemented option"
	STATUS_ERR_UNIMPLEMENTED_EXT      ErrorCode = 2103 //  "Unimplemented extension"
	STATUS_ERR_BILLING                ErrorCode = 2104 //  "Billing failure"
	STATUS_ERR_CANNOT_RENEW           ErrorCode = 2105 //  "Object is not eligible for renewal"
	STATUS_ERR_CANNOT_TRANSFER        ErrorCode = 2106 //  "Object is not eligible for transfer"
	STATUS_ERR_AUTHENTICATION         ErrorCode = 2200 //  "Authentication error"
	STATUS_ERR_AUTHORIZATION          ErrorCode = 2201 //  "Authorization error"
	STATUS_ERR_INVALID_AUTH           ErrorCode = 2202 //  "Invalid authorization information"
	STATUS_ERR_PENDING_TRANSFER       ErrorCode = 2300 //  "Object pending transfer"
	STATUS_ERR_NOT_PENDING_TRANSFER   ErrorCode = 2301 //  "Object not pending transfer"
	STATUS_ERR_EXISTS                 ErrorCode = 2302 //  "Object exists"
	STATUS_ERR_NOT_EXISTS             ErrorCode = 2303 //  "Object does not exist"
	STATUS_ERR_STATUS                 ErrorCode = 2304 //  "Object status prohibits operation"
	STATUS_ERR_ASSOC                  ErrorCode = 2305 //  "Object association prohibits operation"
	STATUS_ERR_PARAM_POLICY           ErrorCode = 2306 //  "Parameter value policy error"
	STATUS_ERR_UNIMLEMENTED_OBJ_SVC   ErrorCode = 2307 //  "Unimplemented object service"
	STATUS_ERR_DATA_MGMT              ErrorCode = 2308 //  "Data management policy violation"
	STATUS_ERR_COMMAND                ErrorCode = 2400 //  "Command failed"
	STATUS_ERR_COMMAND_CLOSING        ErrorCode = 2500 //  "Command failed; server closing connection"
	STATUS_ERR_AUTH_CLOSING           ErrorCode = 2501 //  "Authentication error; server closing connection"
	STATUS_ERR_SESSION_CLOSING        ErrorCode = 2502 //  "Session limit exceeded; server closing connection"

	STATUS_TYPE_CLIENT_DELETE_PROHIBITED   = "clientDeleteProhibited"
	STATUS_TYPE_CLIENT_HOLD                = "clientHold"
	STATUS_TYPE_CLIENT_RENEW_PROHIBITED    = "clientRenewProhibited"
	STATUS_TYPE_CLIENT_TRANSFER_PROHIBITED = "clientTransferProhibited"
	STATUS_TYPE_CLIENT_UPDATE_PROHIBITED   = "clientUpdateProhibited"
	STATUS_TYPE_INACTIVE                   = "inactive"
	STATUS_TYPE_OK                         = "ok"
	STATUS_TYPE_PENDING_CREATE             = "pendingCreate"
	STATUS_TYPE_PENDING_DELETE             = "pendingDelete"
	STATUS_TYPE_PENDING_RENEW              = "pendingRenew"
	STATUS_TYPE_PENDING_TRANSFER           = "pendingTransfer"
	STATUS_TYPE_PENDING_UPDATE             = "pendingUpdate"
	STATUS_TYPE_SERVER_DELETE_PROHIBITED   = "serverDeleteProhibited"
	STATUS_TYPE_SERVER_HOLD                = "serverHold"
	STATUS_TYPE_SERVER_RENEW_PROHIBITED    = "serverRenewProhibited"
	STATUS_TYPE_SERVER_TRANSFER_PROHIBITED = "serverTransferProhibited"
	STATUS_TYPE_SERVER_UPDATE_PROHIBITED   = "serverUpdateProhibited"
	STATUS_TYPE_LINKED                     = "linked"
	STATUS_TYPE_REDEMPTION_PERIOD          = "redemptionPeriod"
	STATUS_TYPE_ADD_PERIOD                 = "addPeriod"
	STATUS_TYPE_AUTO_RENEW_PERIOD          = "autoRenewPeriod"
	STATUS_TYPE_RENEW_PERIOD               = "renewPeriod"
	STATUS_TYPE_TRANSFER_PERIOD            = "transferPeriod"
	STATUS_TYPE_PENDING_RESTORE            = "pendingRestore"
	STATUS_TYPE_TRANSFER_PROHIBITED        = "transferProhibited"
)

var (
	// (alpha)TLD must consist only from letters, except IDN(xn--)
	// @see https://stackoverflow.com/questions/9071279/number-in-the-top-level-domain/53875771#53875771
	HostnameRegex = regexp.MustCompile(`^(([a-z0-9-]|xn--[[:ascii:]]+){1,64}\.)+([a-z]|xn--[[:ascii:]]+){1,64}$`)
	Ipv4Regex     = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	PhoneRegex    = regexp.MustCompile(`^(\+[0-9]{1,3}\.[0-9]{1,14})$`)
	AlphaRegex    = regexp.MustCompile(`^([a-z]|xn--[[:ascii:]]+){1,64}$`)
	TldRegex      = regexp.MustCompile(`^([a-z]|xn--[[:ascii:]]+){1,64}(\.([a-z]|xn--[[:ascii:]]+){1,64})?$`)
	DomainRegex   = regexp.MustCompile(`^([a-z0-9-]|xn--[[:ascii:]]+){1,64}\.([a-z]|xn--[[:ascii:]]+){1,64}(\.([a-z]|xn--[[:ascii:]]+){1,64})?$`)
	IdnRegex      = regexp.MustCompile(`^xn--[[:ascii:]]+$`)
)

func GetResultMessage(code ErrorCode) string {
	messages := map[ErrorCode]string{
		STATUS_OK:                         "Command completed successfully",
		STATUS_OK_ACTION_PENDING:          "Command completed successfully; action pending",
		STATUS_OK_NO_MESSAGES:             "Command completed successfully; no messages",
		STATUS_OK_ACK_NEEDED:              "Command completed successfully; ack to dequeue",
		STATUS_OK_END_SESSION:             "Command completed successfully; ending Session",
		STATUS_ERR_UNKNOWN_COMMAND:        "Unknown command",
		STATUS_ERR_COMMAND_SYNTAX:         "Command syntax error",
		STATUS_ERR_INVALID_COMMAND:        "Command use error",
		STATUS_ERR_MISSING_PARAM:          "Required parameter missing",
		STATUS_ERR_PARAM_RANGE:            "Parameter value range error",
		STATUS_ERR_PARAM_SYNTAX:           "Parameter value syntax error",
		STATUS_ERR_UNIMPLEMENTED_PROTOCOL: "Unimplemented protocol version",
		STATUS_ERR_UNIMPLEMENTED_COMMAND:  "Unimplemented command",
		STATUS_ERR_UNIMPLEMENTED_OPTION:   "Unimplemented option",
		STATUS_ERR_UNIMPLEMENTED_EXT:      "Unimplemented extension",
		STATUS_ERR_BILLING:                "Billing failure",
		STATUS_ERR_CANNOT_RENEW:           "Object is not eligible for renewal",
		STATUS_ERR_CANNOT_TRANSFER:        "Object is not eligible for transfer",
		STATUS_ERR_AUTHENTICATION:         "Authentication error",
		STATUS_ERR_AUTHORIZATION:          "Authorization error",
		STATUS_ERR_INVALID_AUTH:           "Invalid authorization, information",
		STATUS_ERR_PENDING_TRANSFER:       "Object pending transfer",
		STATUS_ERR_NOT_PENDING_TRANSFER:   "Object not pending transfer",
		STATUS_ERR_EXISTS:                 "Object exists",
		STATUS_ERR_NOT_EXISTS:             "Object does not exist",
		STATUS_ERR_STATUS:                 "Object status prohibits, operation",
		STATUS_ERR_ASSOC:                  "Object association prohibits operation",
		STATUS_ERR_PARAM_POLICY:           "Parameter value policy error",
		STATUS_ERR_UNIMLEMENTED_OBJ_SVC:   "Unimplemented object service",
		STATUS_ERR_DATA_MGMT:              "Data management policy violation",
		STATUS_ERR_COMMAND:                "Command failed",
		STATUS_ERR_COMMAND_CLOSING:        "Command failed; server closing connection",
		STATUS_ERR_AUTH_CLOSING:           "Authentication error; server closing connection",
		STATUS_ERR_SESSION_CLOSING:        "Session limit exceeded; server closing connection",
	}

	if msg, ok := messages[code]; ok {
		return msg
	}

	return "unknown status code"
}

func NewResult(code ErrorCode, msg ...interface{}) Result {
	res := Result{Code: code, Message: GetResultMessage(code)}
	if len(msg) != 0 {
		res.Message = fmt.Sprintf(fmt.Sprintf("%v", msg[0]), msg[1:]...)
	}
	return res
}

func InternalErrorResult(msg ...string) Result {
	msg = append([]string{"Internal error."}, msg...)
	// The 2400 for internal errors is according to the RFC docs.
	return Result{Code: STATUS_ERR_COMMAND, Message: strings.Join(msg, " ")}
}

func SuccessResult() Result {
	return NewResult(STATUS_OK)
}

type RequestMessage struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Hello   *Hello   `xml:"hello,omitempty"`
	Command *Command `xml:"command,omitempty"`
}

type ResponseMessage struct {
	XMLName  xml.Name  `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Greeting *Greeting `xml:"greeting,omitempty"`
	Response *Response `xml:"response,omitempty"`
}

type Hello struct{}

type Response struct {
	Result        []Result             `xml:"result"`
	Message       []string             `xml:"msg"`
	MessageQueue  *ResultMessageQueue  `xml:"msgQ,omitempty"`
	ResponseData  *ResponseData        `xml:"resData,omitempty"`
	Extension     *ResponseExtension   `xml:"extension,omitempty"`
	TransactionID *ResultTransactionID `xml:"trID,omitempty"`
}

type ResponseData struct {
	DomainCheckData  *DomainCheckData  `xml:"urn:ietf:params:xml:ns:domain-1.0 chkData,omitempty"`
	ContactCheckData *ContactCheckData `xml:"urn:ietf:params:xml:ns:contact-1.0 chkData,omitempty"`
	HostCheckData    *HostCheckData    `xml:"urn:ietf:params:xml:ns:host-1.0 chkData,omitempty"`

	ContactCreateData *ContactCreateData `xml:"urn:ietf:params:xml:ns:contact-1.0 creData,omitempty"`
	DomainCreateData  *DomainCreateData  `xml:"urn:ietf:params:xml:ns:domain-1.0 creData,omitempty"`
	HostCreateData    *HostCreateData    `xml:"urn:ietf:params:xml:ns:host-1.0 creData,omitempty"`

	ContactInfoData *ContactInfoDataObject `xml:"urn:ietf:params:xml:ns:contact-1.0 infData,omitempty"`
	DomainInfoData  *DomainInfoDataObject  `xml:"urn:ietf:params:xml:ns:domain-1.0 infData,omitempty"`
	HostInfoData    *HostInfoDataObject    `xml:"urn:ietf:params:xml:ns:host-1.0 infData,omitempty"`

	DomainRenewData *DomainRenewData `xml:"urn:ietf:params:xml:ns:domain-1.0 renData,omitempty"`

	DomainTransferData *DomainTransferMessageObject `xml:"urn:ietf:params:xml:ns:domain-1.0 trnData,omitempty"`

	MessageData *string `xml:",innerxml"`
}

type ResponseExtension struct {
	GransyContactInfo *GransyContactObject      `xml:"http://www.subreg.cz/epp/gransy-contact-0.1 infData,omitempty"`
	DomainInfoRGP     *DomainInfoRGPExtension   `xml:"urn:ietf:params:xml:ns:rgp-1.0 infData,omitempty"`
	DomainUpdateRGP   *DomainUpdateRGPExtension `xml:"urn:ietf:params:xml:ns:rgp-1.0 upData,omitempty"`
	DnsSecInfo        *DnsSecObject             `xml:"urn:ietf:params:xml:ns:secDNS-1.1 infData,omitempty"`
}

type DomainCheckData struct {
	Data []CheckDomainDataObject `xml:"cd,omitempty"`
}

type CheckDomainDataObject struct {
	Name   CheckDataObjectName    `xml:"name"`
	Reason *CheckDataObjectReason `xml:"reason,omitempty"`
}

type ContactCheckData struct {
	Data []CheckContactDataObject `xml:"cd,omitempty"`
}

type HostCheckData struct {
	Data []CheckHostDataObject `xml:"cd,omitempty"`
}

type CheckContactDataObject struct {
	ID     CheckDataObjectID      `xml:"id"`
	Reason *CheckDataObjectReason `xml:"reason,omitempty"`
}

type CheckHostDataObject struct {
	Name   CheckDataObjectName    `xml:"name"`
	Reason *CheckDataObjectReason `xml:"reason,omitempty"`
}

type CheckDataObjectName struct {
	Value     string `xml:",chardata"`
	Available int32  `xml:"avail,attr"`
}

type CheckDataObjectID struct {
	Value     string `xml:",chardata"`
	Available int32  `xml:"avail,attr"`
}

type CheckDataObjectReason struct {
	Value    string `xml:",chardata"`
	Language string `xml:"lang,attr,omitempty"`
}

type ContactCreateData struct {
	ID          string `xml:"id"`
	CreatedDate string `xml:"crDate"`
}

type DomainCreateData struct {
	Name           string `xml:"name"`
	CreatedDate    string `xml:"crDate"`
	ExpirationDate string `xml:"exDate"`
}

type HostCreateData struct {
	Name        string `xml:"name"`
	CreatedDate string `xml:"crDate"`
}

type Result struct {
	Code       ErrorCode          `xml:"code,attr"`
	Message    string             `xml:"msg"`
	Value      []interface{}      `xml:"value,omitempty"`
	ExtraValue []ResultExtraValue `xml:"extValue,omitempty"`
}

type ResultExtraValue struct {
	Value  interface{} `xml:"value"`
	Reason string      `xml:"msg,omitempty"`
}

type ResultMessageQueue struct {
	Count      int    `xml:"count,attr"`
	Id         string `xml:"id,attr"`
	QueuedDate string `xml:"qDate,omitempty"`
	Message    string `xml:"msg,omitempty"`
}

type ResultTransactionID struct {
	ClientTransactionID string `xml:"clTRID"`
	ServerTransactionID string `xml:"svTRID"`
}

type Greeting struct {
	ServerName string       `xml:"svID"`
	Date       string       `xml:"svDate"`
	Menu       GreetingMenu `xml:"svcMenu"`
}

type GreetingMenu struct {
	Version    []string               `xml:"version"`
	Language   []string               `xml:"lang"`
	Objects    []string               `xml:"objURI"`
	Extensions []GreetingExtensionURI `xml:"svcExtension,omitempty"`
	// DCP is not implemented, see https://tools.ietf.org/html/rfc5730#section-2.3
}

type GreetingExtensionURI struct {
	ExtensionURI string `xml:"extURI"`
}

type Command struct {
	Login               *LoginCommand    `xml:"login,omitempty"`
	Logout              *LogoutCommand   `xml:"logout,omitempty"`
	Check               *CheckCommand    `xml:"check,omitempty"`
	Poll                *PollCommand     `xml:"poll,omitempty"`
	Transfer            *TransferCommand `xml:"transfer,omitempty"`
	Create              *CreateCommand   `xml:"create,omitempty"`
	Delete              *DeleteCommand   `xml:"delete,omitempty"`
	Renew               *RenewCommand    `xml:"renew,omitempty"`
	Update              *UpdateCommand   `xml:"update,omitempty"`
	Info                *InfoCommand     `xml:"info,omitempty"`
	Extension           CommandExtension `xml:"extension"`
	ClientTransactionID string           `xml:"clTRID"`
}

type CommandExtension struct {
	GransyContactCreate *GransyContactObject      `xml:"http://www.subreg.cz/epp/gransy-contact-0.1 create,omitempty"`
	GransyContactUpdate *GransyContactObject      `xml:"http://www.subreg.cz/epp/gransy-contact-0.1 update,omitempty"`
	UpdateDomainRGP     *UpdateDomainRGPExtension `xml:"urn:ietf:params:xml:ns:rgp-1.0 update,omitempty"`
	DnsSecCreate        *DnsSecObject             `xml:"urn:ietf:params:xml:ns:secDNS-1.1 create,omitempty"`
	DnsSecUpdate        *DnsSecUpdateExtension    `xml:"urn:ietf:params:xml:ns:secDNS-1.1 update,omitempty"`
}

func (ext *CommandExtension) Validate() (ErrorCode, error) {
	if ext.GransyContactCreate != nil {
		if code, err := ext.GransyContactCreate.Validate(); err != nil {
			return code, err
		}
	}

	if ext.GransyContactUpdate != nil {
		if code, err := ext.GransyContactUpdate.Validate(); err != nil {
			return code, err
		}
	}

	if ext.UpdateDomainRGP != nil {
		if code, err := ext.UpdateDomainRGP.Validate(); err != nil {
			return code, err
		}
	}

	if ext.DnsSecCreate != nil {
		if code, err := ext.DnsSecCreate.Validate(); err != nil {
			return code, err
		}
	}

	if ext.DnsSecUpdate != nil {
		if code, err := ext.DnsSecUpdate.Validate(); err != nil {
			return code, err
		}
	}

	return STATUS_OK, nil
}

func (cmd *Command) Validate() (ErrorCode, error) {
	if cmd.Login == nil && cmd.ClientTransactionID == "" {
		return STATUS_ERR_INVALID_COMMAND, ErrNoClientID
	}

	if code, err := cmd.Extension.Validate(); err != nil {
		return code, err
	}

	var count int

	switch {
	case cmd.Login != nil:
		if code, err := cmd.Login.Validate(); err != nil {
			return code, err
		}
		count++
	case cmd.Logout != nil:
		if code, err := cmd.Logout.Validate(); err != nil {
			return code, err
		}
		count++
	case cmd.Check != nil:
		if code, err := cmd.Check.Validate(); err != nil {
			return code, err
		}
		count++
	case cmd.Poll != nil:
		if code, err := cmd.Poll.Validate(); err != nil {
			return code, err
		}
		count++
	case cmd.Transfer != nil:
		if code, err := cmd.Transfer.Validate(); err != nil {
			return code, err
		}
		count++
	case cmd.Create != nil:
		if code, err := cmd.Create.Validate(); err != nil {
			return code, err
		}
		count++
	case cmd.Delete != nil:
		if code, err := cmd.Delete.Validate(); err != nil {
			return code, err
		}
		count++
	case cmd.Renew != nil:
		if code, err := cmd.Renew.Validate(); err != nil {
			return code, err
		}
		count++
	case cmd.Update != nil:
		if code, err := cmd.Update.Validate(); err != nil {
			return code, err
		}
		count++
	case cmd.Info != nil:
		if code, err := cmd.Info.Validate(); err != nil {
			return code, err
		}
		count++
	}

	if count == 0 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoCommand
	}

	if count > 1 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrTooManyCommands
	}

	return STATUS_OK, nil
}

type InfoDataStatusObject struct {
	Status string `xml:"s,attr"`
}

type InnerXML struct {
	Content string `xml:",innerxml"`
}

// According to documentation, the encoder will convert "true" or "false" into boolean value.
type InnerXMLBool struct {
	Value bool `xml:",chardata"`
}

// EPP uses "2006-01-02"(YYYY-MM-DD) format for dates
func DateFromString(d string) (time.Time, error) {
	return time.Parse("2006-01-02", d)
}

func DateToString(d time.Time) string {
	return d.Format("2006-01-02")
}

func DateTimeToString(d time.Time) string {
	return d.Format(time.RFC3339)
}
