package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Socket struct {
	operationTimeout time.Duration
	conn             net.Conn
	r                *bufio.Reader
	w                *bufio.Writer
}

func NewSocket(conn net.Conn, operationTimeout time.Duration) *Socket {
	return &Socket{
		conn:             conn,
		operationTimeout: operationTimeout,
		r:                bufio.NewReader(conn),
		w:                bufio.NewWriter(conn),
	}
}

func (s *Socket) Read(p []byte) (n int, err error) {
	setDeadLineError := s.setDeadLine()

	if setDeadLineError != nil {
		return 0, setDeadLineError
	}

	readBytes, readError := s.r.Read(p)

	if readBytes > 0 {
		return readBytes, nil
	}

	return 0, readError
}

func (s *Socket) setDeadLine() error {
	if s.operationTimeout == 0 {
		return nil
	}

	return s.conn.SetDeadline(time.Now().Add(s.operationTimeout))
}

func (s *Socket) Close() {
	if s.conn != nil {
		err := s.conn.Close()
		if err != nil {
			fmt.Printf("connection close failed (cause: %v)", err)
		}
	}
}
