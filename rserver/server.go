package main

import (
	"chat/shared"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func main() {
	serviceName := fmt.Sprintf(":%d", shared.ServerPort)

	tcpAddr, err := net.ResolveTCPAddr("tcp", serviceName)
	if err != nil {
		log.Fatalf("error while resolving address: %s\n", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("error while creating listener: %s\n", err)
	}
	log.Printf("chat server is started...\n")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("cannot accept connection...\n")
			continue
		}
		log.Printf("client [%s] connected...\n", conn.RemoteAddr().String())

		name, err := clientHandshake(conn)
		if err != nil {
			log.Printf("cannot handshake with client: %s\n", err)
			continue
		}
		addNewClient(name, conn)
		go handleClient(name, conn)
	}
}

func handleClient(name string, conn *net.TCPConn) {
	defer func() {
		log.Printf("connection [%s] closed.\n", conn.RemoteAddr().String())
		removeFromClients(name)
		conn.Close()
	}()

	var buffer [512]byte
	for {
		n, err := conn.Read(buffer[0:])
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				break
			}
			log.Printf("cannot read data from client [%s]: %s, %T\n", conn.RemoteAddr().String(), err, err)
			continue
		}
		cr := new(shared.ClientRequest)
		if err := json.Unmarshal(buffer[:n], cr); err != nil {
			log.Printf("cannot umarshal request from client: %s\n", err)
			continue
		}
		if cr.ToUserName != "" {
			// TODO: handle for single user
			continue
		}
		sendToClientsExcept(cr.Data, cr.Name)
	}
}

func clientHandshake(conn *net.TCPConn) (string, error) {
	log.Printf("handshaking...\n")
	var buffer [512]byte
	n, err := conn.Read(buffer[0:])
	if err != nil {
		return "", err
	}
	log.Printf("name recieved: %s\n", buffer[:n])

	if _, err := conn.Write([]byte(shared.HANDSHAKING_REQREP)); err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}
