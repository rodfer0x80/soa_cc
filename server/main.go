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

func handle(listener net.Listener) {
	// stdin server input and send
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("[!] Fatal: " + fmt.Sprintf("%s", err))
		return
	}
	fmt.Println("[+] Connection estabilished: [" + conn.RemoteAddr().String() + "]")
	for {
		err := conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		if err != nil {
			fmt.Println("[x] Thread Fatal: net.Conn.SetReadDeadline")
			fmt.Println(err)
			fmt.Println("[-] Closed connection: [" + conn.RemoteAddr().String() + "]")
			defer conn.Close()
			return
		}

		reader := bufio.NewReader(os.Stdin)
		recvBuf := make([]byte, 1024)
		fmt.Printf("[%s]$ ", conn.RemoteAddr().String())
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "exit" {
			fmt.Fprintf(conn, text+"\n")
			defer conn.Close()
			fmt.Println("[-] Closed connection: [" + conn.RemoteAddr().String() + "]")
			return
		} else {
			fmt.Fprintf(conn, text+"\n")
		}

		// receive data from client
		_, err = conn.Read(recvBuf[:]) // recv data
		res := string(recvBuf)
		if err != nil {
			// ? timeout error : wait
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Println("[x] Thread Info: netErr.Timeout")
				fmt.Println(err)
				fmt.Println("[-] Closed connection: [" + conn.RemoteAddr().String() + "]")
				defer conn.Close()
				return
			} else {
				fmt.Println("[x] Thread Fatal: net.Conn.Read")
				fmt.Println(err)
				fmt.Println("[-] Closed connection: [" + conn.RemoteAddr().String() + "]")
				defer conn.Close()
				return
			}
		} else {
			if res == "exit" {
				fmt.Println("[-] Client closed connection: [" + conn.RemoteAddr().String() + "]")
				defer conn.Close()
				return
			}
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
		handle(listener)
	}
}
