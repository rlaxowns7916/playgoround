package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := fmt.Sprintf(":%s", arguments[1])
	serverSocket, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Error listening on port %s: %s\n", port, err)
		return
	}

	wg := sync.WaitGroup{}

	go func() {
		for {
			clientSocket, err := serverSocket.Accept()
			if err == nil {
				wg.Add(1)
				go echo(clientSocket, &wg)
			}
		}
	}()

	<-ctx.Done()
	serverSocket.Close()
	wg.Wait()
	fmt.Println("server shutdown")

}

func echo(clientSocket net.Conn, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		clientSocket.Close()
	}()

	r := bufio.NewReader(clientSocket)
	w := bufio.NewWriter(clientSocket)

	for {
		input, err := r.ReadString('\n')
		if len(input) > 0 {
			if "QUIT" == strings.TrimSpace(input) {
				return
			}
			fmt.Printf(" >> [%s] %s\n", clientSocket.RemoteAddr(), input)

			_, err = w.WriteString(input)
			if err != nil {
				fmt.Printf("Error writing to client: %s\n", err)
				return
			}

			if err = w.Flush(); err != nil {
				return
			}
			fmt.Printf(" << [%s] %s\n", clientSocket.RemoteAddr(), input)
		}

		if err != nil {
			switch err {
			case io.EOF:
				fmt.Println("Client closed the connection")
			default:
				fmt.Println("Error reading from client:", err)
			}

			return
		}
	}
}
