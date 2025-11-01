package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port")
		return
	}

	conn, err := net.Dial("tcp", arguments[1])
	if err != nil {
		fmt.Printf("Error connecting to server: %s\n", err)
		return
	}

	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			fmt.Printf("Error closing connection: %s\n", closeErr)
			return
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")

		input, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, input+"\n")

		output, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("->: " + output)

		if strings.TrimSpace(output) == "STOP" {
			fmt.Println("TCP Client exiting ...")
			return
		}
	}
}
