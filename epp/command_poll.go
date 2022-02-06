package epp

const (
	EPP_POLL_ACKNOWLEDGE_MESSAGE_OPERATION = "ack"
	EPP_POLL_REQUEST_MESSAGE_OPERATION     = "req"
)

type PollCommand struct {
	Operation string `xml:"op,attr"`
	MessageID string `xml:"msgID,attr,omitempty"`
}

func (cmd *PollCommand) Validate() (ErrorCode, error) {
	switch cmd.Operation {
	case EPP_POLL_REQUEST_MESSAGE_OPERATION:
		if cmd.MessageID == "" {
			return STATUS_OK, nil
		} else {
			return STATUS_ERR_COMMAND_SYNTAX, ErrUnexpectedMessageID
		}
	case EPP_POLL_ACKNOWLEDGE_MESSAGE_OPERATION:
		if cmd.MessageID == "" {
			return STATUS_ERR_COMMAND_SYNTAX, ErrNoMessageID
		} else {
			return STATUS_OK, nil
		}
	default:
		return STATUS_ERR_COMMAND_SYNTAX, ErrOperation
	}
}
