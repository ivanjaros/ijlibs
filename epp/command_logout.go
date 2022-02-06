package epp

type LogoutCommand struct{}

func (cmd *LogoutCommand) Validate() (ErrorCode, error) {
	return STATUS_OK, nil
}
