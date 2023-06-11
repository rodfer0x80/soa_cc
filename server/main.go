package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const IFACE = "127.0.0.1"
const PORT = 42000

func handle(conn net.Conn) {
	// • Establish a connection to a remote listener via net.Dial(network, address string).
	// • Initialize a Cmd via exec.Command(name string, arg ...string).
	// • Redirect Stdin and Stdout properties to utilize the net.Conn object.
	// • Run the command.
	reader := bufio.NewReader(os.Stdin)
	recvBuf := make([]byte, 1024)

	// stdin server input and send
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	fmt.Fprintf(conn, text+"\n")
	for {
		if IFACE != "127.0.0.1" {
			err := conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			if err != nil {
				log.Println("[x] SetReadDeadline failed:", err)
				conn.Close()
				return
			}
		}

		// ## builtin command
		// close connection
		if text == "exit\n" {
			log.Println("[-] Closed connection with " + fmt.Sprintf("%s", conn.RemoteAddr().String()))
			conn.Close()
			return
		}

		// receive data from client
		recvBuf = make([]byte, 1024)
		_, err := conn.Read(recvBuf[:]) // recv data
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("[!] Read timeout:", err)
				// time out
			} else {
				log.Println("[!] Read error:", err)
				conn.Close()
				// some error else, do something else, for example create new conn
			}
		}

		// stdout client output
		print(string(recvBuf))

		// stdin server input and send
		if string(recvBuf) != "" {
			text, _ = reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			fmt.Fprintf(conn, text+"\n")
		}
	}
}

func main() {
	port := PORT
	// create a pool for 100 connections
	// when connection is closed remove from pool
	// and add a fresh new one on same port listening
	address := fmt.Sprintf("%s:%d", IFACE, PORT)
	for {
		// Increase ports for connection pool
		port = port + 1
		address = fmt.Sprintf("%s:%d", IFACE, port)

		// Accept multiple connections
		listener, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatalln("[x] Unable to bind to port")
		}
		log.Println("[*] Listening on " + address)

		// Wait connections
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("[x] Unable to accept connection from " + fmt.Sprintf("%s", conn.RemoteAddr().String()))
		}
		log.Println("[+] Received connection " + fmt.Sprintf("%s", conn.RemoteAddr().String()))

		// functionality to waitGroup and finish threads so we can accept new connection on same port
		go handle(conn)
	}
}
