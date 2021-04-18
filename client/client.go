package main

import (
	"bufio"
	"chat/shared"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

const clientServerPort = 2222

var name *string = flag.String("username", "", "input your username")

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	if *name == "" {
		*name += strconv.Itoa(rand.Intn(1e6))
	}

	log.Printf("chat is starting...\n")
	serverName := fmt.Sprintf(":%d", shared.ServerPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverName)
	if err != nil {
		log.Fatalf("cannot resolve address: %s\n", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("cannot connect to the remote server: %s\n", err)
	}

	// send client username to server
	if _, err = conn.Write([]byte(*name)); err != nil {
		log.Fatalf("cannot send username and handshake: %s\n", err)
	}
	// read confirmation from server
	var buffer [512]byte
	n, err := conn.Read(buffer[0:])
	if err != nil {
		log.Fatalf("cannot read handshake confirm from server: %s\n", err)
	}
	if string(buffer[:n]) != shared.HANDSHAKING_REQREP {
		log.Fatalf("confirmation from server not equal")
	}
	log.Printf("server accept your client\n")

	// read chat messages
	go func(conn *net.TCPConn) {
		for {
			n, err := conn.Read(buffer[0:])
			if err != nil {
				if err == io.EOF {
					log.Printf("connection closed: %s\n", err)
					break
				}
				if _, ok := err.(*net.OpError); ok {
					log.Printf("connection closed: %s\n", err)
					break
				}
				log.Printf("cannot read data from server: %s\n", err)
				continue
			}
			msg := new(shared.Message)
			if err := json.Unmarshal(buffer[:n], msg); err != nil {
				log.Printf("cannot unmarshal data from server: %s\n", err)
			}
			fmt.Printf("%s[%s]<%s>: %s%s\n", shared.T_COLOR_BLUE,
				msg.Time.Format("2006-01-02 15:04:05"),
				msg.Name,
				msg.Message,
				shared.T_COLOR_CLOSING_TAG)
		}
	}(conn)
	// send messages
	for {
		fmt.Print("\rSend message: ")
		buffer := bufio.NewReader(os.Stdin)
		msg, err := buffer.ReadString('\n')
		if err != nil {
			log.Printf("cannot read message from buffer: %s\n", err)
			continue
		}
		req := &shared.ClientRequest{
			Name: *name,
			Data: msg,
		}
		bt, err := json.Marshal(req)
		if err != nil {
			log.Printf("cannot marshal request: %s\n", err)
			continue
		}
		conn.Write(bt)
	}
}
