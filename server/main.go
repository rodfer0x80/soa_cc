package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const HOST = "127.0.0.1"
const PORT = 42000

func handle(conn net.Conn) {
	// stdin server input and send
	for {
		reader := bufio.NewReader(os.Stdin)
		recvBuf := make([]byte, 1024)
		fmt.Printf("[%s]$ ", conn.RemoteAddr().String())
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "exit" {
			fmt.Fprintf(conn, text+"\n")
			conn.Close()
			fmt.Println("[-] Closed connection: [" + conn.RemoteAddr().String() + "]")
			return
		} else {
			fmt.Fprintf(conn, text+"\n")
		}

		err := conn.SetReadDeadline(time.Now().Add(time.Second * 60))
		if err != nil {
			fmt.Println("[x] Thread Fatal: net.Conn.SetReadDeadline")
			fmt.Println(err)
			conn.Close()
			return
		}

		// receive data from client
		_, err = conn.Read(recvBuf[:]) // recv data
		res := string(recvBuf)
		if err != nil {
			// ? timeout error : wait
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("[x] Thread Info: netErr.Timeout")
				time.Sleep(time.Second * 60)
				//handle(conn)
			} else {
				fmt.Println("[x] Thread Fatal: net.Conn.Read")
				fmt.Println(err)
				conn.Close()
				return
			}
		} else {
			print(res)
		}
	}
}

func main() {
	address := fmt.Sprintf("%s:%d", HOST, PORT)

	listener, err := net.Listen("tcp4", address)
	if err != nil {
		fmt.Println("[!] Fatal: " + fmt.Sprintf("%s", err))
		return
	}
	fmt.Println("[*] Listening: [" + address + "]")

	defer listener.Close()
	rand.Seed(time.Now().Unix())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[!] Fatal: " + fmt.Sprintf("%s", err))
			return
		}
		fmt.Println("[+] Connection estabilished: [" + conn.RemoteAddr().String() + "]")

		handle(conn)
	}
}
