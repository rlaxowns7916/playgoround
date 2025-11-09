package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
)

type TcpServer struct {
	port     int
	wg       sync.WaitGroup
	listener net.Listener
}

func NewTcpServer(port int) *TcpServer {
	return &TcpServer{
		port:     port,
		wg:       sync.WaitGroup{},
		listener: nil,
	}
}

func (t *TcpServer) Start() error {
	if t.listener == nil {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", t.port))
		if err != nil {
			return err
		}
		t.listener = listener
	}

	return nil
}

func (t *TcpServer) Close() error {
	if t.listener == nil {
		return nil
	}

	if err := t.listener.Close(); err != nil {
		return err
	}
	return nil
}

func (t *TcpServer) Wait() {
	t.wg.Wait()
}

func (t *TcpServer) Serve(ctx context.Context, handler func(context.Context, net.Conn)) error {
	if t.listener == nil {
		return fmt.Errorf("listener not started")
	}

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				fmt.Println("server closed")
				break
			}
			fmt.Println("accept error:", err)
			continue
		}

		t.wg.Add(1)
		go func(conn net.Conn) {
			defer func() {
				t.wg.Done()
				_ = conn.Close()
			}()

			handler(ctx, conn)
		}(conn)
	}
	return nil
}
