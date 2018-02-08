package api

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"time"
)

type RPCType byte

const (
	RPCRiff RPCType = 0
	RPCLog          = 1
)

func NewClient(address string) (*rpc.Client, error) {
	conn, err := net.DialTimeout("tcp", address, time.Second*10)
	if err != nil {
		return nil, err
	}
	conn.Write([]byte{byte(RPCRiff)})
	codec := NewGobClientCodec(conn)
	client := rpc.NewClientWithCodec(codec)
	return client, nil
}

func TimeoutCoder(f func(interface{}) error, e interface{}, msg string) error {
	endChan := make(chan error, 1)
	go func() { endChan <- f(e) }()
	select {
	case e := <-endChan:
		return e
	case <-time.After(time.Minute):
		return fmt.Errorf("Timeout %s", msg)
	}
}

//server
func NewGobServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	buf := bufio.NewWriter(conn)
	codec := &gobServerCodec{
		rwc:    conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewEncoder(buf),
		encBuf: buf,
	}
	return codec
}

type gobServerCodec struct {
	rwc    io.ReadWriteCloser
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
	closed bool
}

func (c *gobServerCodec) ReadRequestHeader(r *rpc.Request) error {
	return TimeoutCoder(c.dec.Decode, r, "server read request header")
}

func (c *gobServerCodec) ReadRequestBody(body interface{}) error {
	return TimeoutCoder(c.dec.Decode, body, "server read request body")
}

func (c *gobServerCodec) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	if err = TimeoutCoder(c.enc.Encode, r, "server write response"); err != nil {
		if c.encBuf.Flush() == nil {
			log.Println("rpc: gob error encoding response:", err)
			c.Close()
		}
		return
	}
	if err = TimeoutCoder(c.enc.Encode, body, "server write response body"); err != nil {
		if c.encBuf.Flush() == nil {
			log.Println("rpc: gob error encoding body:", err)
			c.Close()
		}
		return
	}
	return c.encBuf.Flush()
}

func (c *gobServerCodec) Close() error {
	if c.closed {
		// Only call c.rwc.Close once; otherwise the semantics are undefined.
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

//client
func NewGobClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	buf := bufio.NewWriter(conn)
	codec := &gobClientCodec{
		rwc:    conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewEncoder(buf),
		encBuf: buf,
	}
	return codec
}

type gobClientCodec struct {
	rwc    io.ReadWriteCloser
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
}

func (c *gobClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	if err = TimeoutCoder(c.enc.Encode, r, "client write request"); err != nil {
		return
	}
	if err = TimeoutCoder(c.enc.Encode, body, "client write request body"); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *gobClientCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.dec.Decode(r)
}

func (c *gobClientCodec) ReadResponseBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *gobClientCodec) Close() error {
	return c.rwc.Close()
}
