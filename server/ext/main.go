package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const IFACE = "127.0.0.1"
const PORT = 42000
const WORKERS = 4

func handle(wg *sync.WaitGroup, shell_lock *sync.Mutex, port int) {
	// • Establish a connection to a remote listener via net.Dial(network, address string).
	// • Initialize a Cmd via exec.Command(name string, arg ...string).
	// • Redirect Stdin and Stdout properties to utilize the net.Conn object.
	// • Run the command.
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", IFACE, port)

	// Accept multiple connections
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("[x] Unable to bind to port " + fmt.Sprintf("%d", port))
	}
	log.Println("[*] Listening on " + address)

	// Wait connections
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalln("[x] Unable to accept connection from " + fmt.Sprintf("%s", conn.RemoteAddr().String()))
	}
	log.Println("[+] Received connection " + fmt.Sprintf("%s", conn.RemoteAddr().String()))

	shell_lock.Lock()

	// stdin server input and send
	reader := bufio.NewReader(os.Stdin)
	recvBuf := make([]byte, 1024)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	fmt.Fprintf(conn, text+"\n")
	for {
		if IFACE != "127.0.0.1" {
			err := conn.SetReadDeadline(time.Now().Add(time.Second))
			if err != nil {
				log.Println("[x] SetReadDeadline failed:", err)
				conn.Close()
				shell_lock.Unlock()
				time.Sleep(time.Second * 30)
				handle(wg, shell_lock, port)
			}
		}

		// receive data from client
		recvBuf = make([]byte, 1024)
		_, err := conn.Read(recvBuf[:]) // recv data
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("[!] Read timeout:", err)
				// time out
				time.Sleep(time.Second * 60)
			} else {
				log.Println("[!] Read error:", err)
				conn.Close()
				shell_lock.Unlock()
				time.Sleep(time.Second * 30)
				handle(wg, shell_lock, port)
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

		// ## builtin command
		// close connection
		if text == "exit" {
			log.Println("[-] Closed connection with " + fmt.Sprintf("%s", conn.RemoteAddr().String()))
			conn.Close()
			shell_lock.Unlock()
			time.Sleep(time.Second * 30)
			handle(wg, shell_lock, port)
		}
	}
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	//pipe := make(chan bool)
	shell_lock := new(sync.Mutex)
	var port int = PORT

	// create a pool for 100 connections
	// when connection is closed remove from pool
	// and add a fresh new one on same port listening
	// functionality to waitGroup and finish threads so we can accept new connection on same port
	for id := 1; id <= WORKERS; id++ {
		wg.Add(1)

		port = port + 1

		go handle(wg, shell_lock, port)

		// Increase ports for connection pool
	}
	wg.Wait()
}
