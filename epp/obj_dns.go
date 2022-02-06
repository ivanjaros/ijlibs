package epp

import (
	"github.com/ivanjaros/jslibs/epp/utils"
)

// DNS SEC algorithms according to https://tools.ietf.org/html/rfc8624
const (
	DnsSec_Key_Alg_RSAMD5             = 1
	DnsSec_Key_Alg_DSA                = 3
	DnsSec_Key_Alg_RSASHA1            = 5
	DnsSec_Key_Alg_DSA_NSEC3_SHA1     = 6
	DnsSec_Key_Alg_RSASHA1_NSEC3_SHA1 = 7
	DnsSec_Key_Alg_RSASHA256          = 8
	DnsSec_Key_Alg_RSASHA512          = 10
	DnsSec_Key_Alg_ECC_GOST           = 12
	DnsSec_Key_Alg_ECDSAP256SHA256    = 13
	DnsSec_Key_Alg_ECDSAP384SHA384    = 14
	DnsSec_Key_Alg_ED25519            = 15
	DnsSec_Key_Alg_ED448              = 16

	// Delegation Signer Digest Algorithms
	DnsSec_Digest_Alg_SHA1            = 1
	DnsSec_Digest_Alg_SHA256          = 2
	DnsSec_Digest_Alg_GOST_R_34_11_94 = 3
	DnsSec_Digest_Alg_SHA384          = 4
)

func IsDnsSecKeyAlg(alg int) bool {
	switch alg {
	case DnsSec_Key_Alg_RSAMD5,
		DnsSec_Key_Alg_DSA,
		DnsSec_Key_Alg_RSASHA1,
		DnsSec_Key_Alg_DSA_NSEC3_SHA1,
		DnsSec_Key_Alg_RSASHA1_NSEC3_SHA1,
		DnsSec_Key_Alg_RSASHA256,
		DnsSec_Key_Alg_RSASHA512,
		DnsSec_Key_Alg_ECC_GOST,
		DnsSec_Key_Alg_ECDSAP256SHA256,
		DnsSec_Key_Alg_ECDSAP384SHA384,
		DnsSec_Key_Alg_ED25519,
		DnsSec_Key_Alg_ED448:
		return true
	default:
		return false
	}
}

func IsDnsSecDigestAlg(alg int) bool {
	switch alg {
	case DnsSec_Digest_Alg_SHA1,
		DnsSec_Digest_Alg_SHA256,
		DnsSec_Digest_Alg_GOST_R_34_11_94,
		DnsSec_Digest_Alg_SHA384:
		return true
	default:
		return false
	}
}

func DnsSecDigestLengthMatch(alg int, value string) bool {
	switch alg {
	case DnsSec_Digest_Alg_SHA1:
		return len(value) == 40
	case DnsSec_Digest_Alg_SHA256:
		return len(value) == 64
	case DnsSec_Digest_Alg_GOST_R_34_11_94:
		return len(value) == 64
	case DnsSec_Digest_Alg_SHA384:
		return len(value) == 94
	default:
		return false
	}
}

type DnsSecObject struct {
	MaxSignatureLife *uint                `xml:"maxSigLife,omitempty"`
	DsData           []DnsSecDsDataObject `xml:"dsData,omitempty"`
	KeyData          []DnsSecDsKeyObject  `xml:"keyData,omitempty"`
}

func (obj *DnsSecObject) Validate() (ErrorCode, error) {
	// unsupported via https://tools.ietf.org/html/rfc5910#section-5.2.1
	if obj.MaxSignatureLife != nil {
		return STATUS_ERR_UNIMPLEMENTED_OPTION, ErrDnsMaxSigLife
	}

	// Either dsData or keyData objects can be present, but not both
	if len(obj.DsData) > 0 && len(obj.KeyData) > 0 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDNSBothDSTypes
	}

	for k := range obj.DsData {
		if c, e := obj.DsData[k].Validate(); e != nil {
			return c, e
		}
		for i := range obj.DsData {
			if i != k && obj.DsData[k].Equals(obj.DsData[i]) {
				return STATUS_ERR_COMMAND_SYNTAX, ErrDNSDuplicate
			}
		}
	}

	for k := range obj.KeyData {
		if c, e := obj.KeyData[k].Validate(); e != nil {
			return c, e
		}
		for i := range obj.KeyData {
			if i != k && obj.KeyData[k].Equals(obj.KeyData[i]) {
				return STATUS_ERR_COMMAND_SYNTAX, ErrDNSDuplicate
			}
		}
	}

	return STATUS_OK, nil
}

type DnsSecDsDataObject struct {
	KeyTag      uint               `xml:"keyTag"`
	Algorithm   uint               `xml:"alg"`
	DigestType  uint               `xml:"digestType"`
	DigestValue string             `xml:"digest"`
	KeyData     *DnsSecDsKeyObject `xml:"keyData,omitempty"`
}

func (obj *DnsSecDsDataObject) Validate() (ErrorCode, error) {
	// https://tools.ietf.org/html/rfc8624#section-3.1
	if IsDnsSecKeyAlg(int(obj.Algorithm)) == false {
		return STATUS_ERR_PARAM_RANGE, ErrDnsAlgorithm
	}

	// https://tools.ietf.org/html/rfc8624#section-3.3
	if IsDnsSecDigestAlg(int(obj.DigestType)) == false {
		return STATUS_ERR_PARAM_RANGE, ErrDnsDigestType
	}

	if DnsSecDigestLengthMatch(int(obj.DigestType), obj.DigestValue) == false {
		return STATUS_ERR_PARAM_RANGE, ErrDnsDigest
	}

	if obj.KeyData != nil {
		if c, e := obj.KeyData.Validate(); e != nil {
			return c, e
		}
	}

	return STATUS_OK, nil
}

func (obj *DnsSecDsDataObject) Equals(cpr DnsSecDsDataObject) bool {
	if obj.KeyData != nil && cpr.KeyData != nil && obj.KeyData.Equals(*cpr.KeyData) == false {
		return false
	}
	return obj.DigestValue == cpr.DigestValue && obj.DigestType == cpr.DigestType && obj.Algorithm == cpr.Algorithm && obj.KeyTag == cpr.KeyTag
}

type DnsSecDsKeyObject struct {
	Flags     uint   `xml:"flags"`
	Protocol  uint   `xml:"protocol"`
	Algorithm uint   `xml:"alg"`
	PublicKey string `xml:"pubKey"`
}

func (obj *DnsSecDsKeyObject) Validate() (ErrorCode, error) {
	// https://tools.ietf.org/html/rfc4034#section-2.2
	if utils.IntInArray(int(obj.Flags), []int{0, 256, 257}) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDnsFlag
	}

	// https://tools.ietf.org/html/rfc4034#section-2.1.2
	if obj.Protocol != 3 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDnsProtocol
	}

	// https://tools.ietf.org/html/rfc8624#section-3.1
	if IsDnsSecKeyAlg(int(obj.Algorithm)) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDnsAlgorithm
	}

	if utils.IsBase64(obj.PublicKey) == false {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDnsInvalidPubKey
	}

	return STATUS_OK, nil
}

func (obj *DnsSecDsKeyObject) Equals(that DnsSecDsKeyObject) bool {
	return obj.Flags == that.Flags && obj.Protocol == that.Protocol && obj.Algorithm == that.Algorithm && obj.PublicKey == that.PublicKey
}

type DnsSecUpdateExtension struct {
	Urgent       bool                      `xml:"urgent,attr"`
	AddAction    *DnsSecObject             `xml:"add,omitempty"`
	RemoveAction *DnsSecUpdateRemoveAction `xml:"rem,omitempty"`
	ChangeAction *DnsSecUpdateChangeAction `xml:"chg,omitempty"`
}

func (obj *DnsSecUpdateExtension) Validate() (ErrorCode, error) {
	// unsupported via https://tools.ietf.org/html/rfc5910#section-5.2.5
	if obj.Urgent {
		return STATUS_ERR_UNIMPLEMENTED_OPTION, ErrDnsUrgent
	}

	var counter int

	if obj.AddAction != nil {
		counter++

		if c, e := obj.AddAction.Validate(); e != nil {
			return c, e
		}
	}

	if obj.RemoveAction != nil {
		counter++

		if c, e := obj.RemoveAction.Validate(); e != nil {
			return c, e
		}
	}

	if obj.ChangeAction != nil {
		counter++

		if c, e := obj.ChangeAction.Validate(); e != nil {
			return c, e
		}
	}

	// at least one action has to be provided: https://tools.ietf.org/html/rfc5910#section-5.2.5
	if counter == 0 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrDnsUpdateNothing
	}

	return STATUS_OK, nil
}

type DnsSecUpdateChangeAction struct {
	MaxSignatureLife uint `xml:"maxSigLife,omitempty"`
}

func (obj *DnsSecUpdateChangeAction) Validate() (ErrorCode, error) {
	// unsupported via https://tools.ietf.org/html/rfc5910#section-5.2.1
	return STATUS_ERR_UNIMPLEMENTED_OPTION, ErrDnsMaxSigLife
}

type DnsSecUpdateRemoveAction struct {
	All     InnerXMLBool         `xml:"all"`
	DsData  []DnsSecDsDataObject `xml:"dsData,omitempty"`
	KeyData []DnsSecDsKeyObject  `xml:"keyData,omitempty"`
}

func (obj *DnsSecUpdateRemoveAction) Validate() (ErrorCode, error) {
	// if the "all" option is provided, no other values should be present
	if obj.All.Value && len(obj.DsData) > 0 && len(obj.KeyData) > 0 {
		return STATUS_ERR_COMMAND_SYNTAX, ErrCommandSyntax
	}

	for k := range obj.DsData {
		if c, e := obj.DsData[k].Validate(); e != nil {
			return c, e
		}
		for i := range obj.DsData {
			if i != k && obj.DsData[k].Equals(obj.DsData[i]) {
				return STATUS_ERR_COMMAND_SYNTAX, ErrDNSDuplicate
			}
		}
	}

	for k := range obj.KeyData {
		if c, e := obj.KeyData[k].Validate(); e != nil {
			return c, e
		}
		for i := range obj.KeyData {
			if i != k && obj.KeyData[k].Equals(obj.KeyData[i]) {
				return STATUS_ERR_COMMAND_SYNTAX, ErrDNSDuplicate
			}
		}
	}

	return STATUS_OK, nil
}
