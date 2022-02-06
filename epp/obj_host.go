package epp

import (
	"github.com/ivanjaros/jslibs/epp/utils"
	"net"
)

type HostObject struct {
	Name string   `xml:"name"`
	IPs  []string `xml:"addr"`
}

func (obj *HostObject) ValidateCreate() (ErrorCode, error) {
	if utils.LengthRange(obj.Name, 1, 255) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNSName
	}

	if utils.NumRange(len(obj.IPs), 0, 10) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNSAddrCount
	}

	for _, ip := range obj.IPs {
		if net.ParseIP(ip) == nil {
			return STATUS_ERR_COMMAND_SYNTAX, ErrIP
		}
	}

	if len(obj.IPs) != len(utils.ArrayUnique(obj.IPs)) {
		return STATUS_ERR_COMMAND_SYNTAX, ErrIPDuplicate
	}

	return STATUS_OK, nil
}

type HostInfoDataObject struct {
	Name         string                 `xml:"name"`
	StorageID    string                 `xml:"roid"`
	Statuses     []InfoDataStatusObject `xml:"status"`
	Owner        string                 `xml:"clid"`
	Creator      string                 `xml:"crid"`
	UpdatedBy    string                 `xml:"upID,omitempty"`
	CreatedDate  string                 `xml:"crDate"`
	UpdatedDate  string                 `xml:"upDate,omitempty"`
	TransferDate string                 `xml:"trdate,omitempty"`
	IPs          []HostInfoIPObject     `xml:"addr"`
}

type HostInfoIPObject struct {
	Address string `xml:",chardata"`
	Type    string `xml:"type,attr"`
}
