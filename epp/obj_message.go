package epp

type DomainTransferMessageObject struct {
	Name              string `xml:"name"`
	Status            string `xml:"trStatus"`
	RequesteeLogin    string `xml:"reID"`
	RequestDate       string `xml:"reDate"`
	RegistrarLogin    string `xml:"acID"`
	ValidUntilDate    string `xml:"acDate"`
	NewExpirationDate string `xml:"exDate"`
	Period            int    `xml:"period"`
}
