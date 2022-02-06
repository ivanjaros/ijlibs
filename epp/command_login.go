package epp

type LoginCommand struct {
	ClientId    string                `xml:"clID"`
	Password    string                `xml:"pw"`
	NewPassword string                `xml:"newPW"`
	Options     LoginCommandOptions   `xml:"options"`
	Services    []LoginServicesObject `xml:"svcs"`
}

type LoginCommandOptions struct {
	Version  string `xml:"version"`
	Language string `xml:"lang"`
}

type LoginServicesObject struct {
	Objects    []string              `xml:"objURI"`
	Extensions LoginExtensionsObject `xml:"svcExtension"`
}

type LoginExtensionsObject struct {
	Extensions []string `xml:"extURI"`
}

func (cmd *LoginCommand) Validate() (ErrorCode, error) {
	if cmd.ClientId == "" {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoClientID
	}

	if cmd.Password == "" {
		return STATUS_ERR_COMMAND_SYNTAX, ErrNoPassword
	}

	if cmd.Options.Version != EPP_PROTO_V1 && cmd.Options.Version != EPP_PROTO_V2 {
		return STATUS_ERR_UNIMPLEMENTED_PROTOCOL, ErrVersion
	}

	if cmd.Options.Language != EPP_PROTO_LANG_EN {
		return STATUS_ERR_UNIMPLEMENTED_OPTION, ErrLanguage
	}

	return STATUS_OK, nil
}
