package event_stream

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"time"
)

type Message struct {
	Retry time.Duration
	Name  string
	ID    string
	Data  interface{}
}

// this is not memory-efficient/optimized but we are working with strings so...
func (e *Message) Encode(w io.Writer) error {
	var out string

	if e.ID != "" {
		out += "id: "
		out += e.ID
		out += "\n"
	}

	if e.Name != "" {
		out += "event: "
		out += e.Name
		out += "\n"
	}

	if e.Retry > 0 {
		out += "retry: "
		out += strconv.FormatInt(e.Retry.Milliseconds(), 10)
		out += "\n"
	}

	if e.Data != nil {
		data, err := json.Marshal(e.Data)
		if err != nil {
			return err
		}

		if bytes.Index(data, []byte("\n")) != -1 {
			chunks := bytes.Split(data, []byte("\n"))
			for k := range chunks {
				out += "data: "
				out += string(chunks[k])
				out += "\n"
			}
		} else {
			out += "data: "
			out += string(data)
			out += "\n"
		}
	}

	if out == "" {
		return nil
	}

	_, err := w.Write([]byte(out + "\n"))
	return err
}
