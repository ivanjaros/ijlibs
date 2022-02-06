package epp

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io"
	"net"
	"time"
)

func NewServerConnection(c net.Conn) *srvConn {
	return &srvConn{conn: c}
}

// Provides handling of Length-Value protocol that the EPP uses.
// EPP uses the length header as part of the message size so it has to be
// counted in when the message is being sent and counted out when the message
// is being received.
type srvConn struct {
	conn net.Conn
}

func (c *srvConn) Read() (RequestMessage, error) {
	var req RequestMessage
	var header uint32 // int32/uint32 is represented as 4 bytes in binary form
	if err := binary.Read(c.conn, binary.BigEndian, &header); err != nil {
		return req, err
	}

	rawData := make([]byte, header-4)
	if _, err := io.ReadFull(c.conn, rawData); err != nil {
		return req, err
	}

	if err := xml.Unmarshal(rawData, &req); err != nil {
		return req, err
	}

	return req, nil
}

func (c *srvConn) Send(response ResponseMessage) error {
	data, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		return err
	}
	return c.SendRaw(data)
}

// note that the xmlData should be only the raw xml data of the marshaled
// ResponseMessage and nothing else.
func (c *srvConn) SendRaw(xmlData []byte) error {
	var data []byte
	if bytes.HasPrefix(xmlData, []byte("<?xml")) {
		data = append(make([]byte, 4), xmlData...)
	} else {
		data = append(make([]byte, 4), append([]byte(xml.Header), xmlData...)...)
	}
	binary.BigEndian.PutUint32(data, uint32(len(data)))
	_, err := c.conn.Write(data)
	return err
}

func (c *srvConn) Close() error {
	return c.conn.Close()
}

func (c *srvConn) SetDeadline(delay time.Time) error {
	return c.conn.SetDeadline(delay)
}

type ServerConn interface {
	Read() (RequestMessage, error)
	Send(ResponseMessage) error
	SendRaw([]byte) error
	Close() error
	SetDeadline(time.Time) error
}
