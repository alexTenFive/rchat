package main

import (
	"chat/client/ui"
	"chat/shared"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

const clientServerPort = 2222

var (
	name        *string = flag.String("username", "", "input your username")
	msgReceiver         = make(chan shared.TerminalData, 20)
	msgSender           = make(chan string, 20)
	chQuit              = make(chan bool)
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	if *name == "" {
		*name += strconv.Itoa(rand.Intn(1e6))
	}

	go ui.InitUI(msgReceiver, msgSender, chQuit)
	go func(chQuit chan bool) {
		<-chQuit
		os.Exit(0)
	}(chQuit)

	serverName := fmt.Sprintf(":%d", shared.ServerPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverName)
	if err != nil {
		printData(fmt.Sprintf("cannot resolve address: %s", err))
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		printData(fmt.Sprintf("cannot connect to the remote server: %s", err))
		os.Exit(1)
	}

	// send client username to server
	if _, err = conn.Write([]byte(*name)); err != nil {
		printData(fmt.Sprintf("cannot send username and handshake: %s", err))
		os.Exit(1)
	}
	// read confirmation from server
	var buffer [512]byte
	n, err := conn.Read(buffer[0:])
	if err != nil {
		printData(fmt.Sprintf("cannot read handshake confirm from server: %s", err))
		os.Exit(1)
	}
	if string(buffer[:n]) != shared.HANDSHAKING_REQREP {
		printData("confirmation from server not equal")
		os.Exit(1)
	}

	// read chat messages
	go func(conn *net.TCPConn) {
		var buffer [512]byte
		for {
			n, err := conn.Read(buffer[0:])
			if err != nil {
				if err == io.EOF {
					printData(fmt.Sprintf("connection closed: %s", err))
					break
				}
				if _, ok := err.(*net.OpError); ok {
					printData(fmt.Sprintf("connection closed: %s", err))
					break
				}
				printData(fmt.Sprintf("cannot read data from server: %s", err))
				continue
			}
			msg := new(shared.Message)
			if err := json.Unmarshal(buffer[:n], msg); err != nil {
				printData(fmt.Sprintf("cannot unmarshal data from server: %s", err))
				continue
			}
			fmt.Println(msg.Name)
			printData(fmt.Sprintf("[%s]<%s>: %s",
				msg.Time.Format("2006-01-02 15:04:05"),
				msg.Name,
				msg.Message))
		}
	}(conn)
	// send messages
	for {
		msg := <-msgSender
		req := &shared.ClientRequest{
			Name: *name,
			Data: msg,
		}
		bt, err := json.Marshal(req)
		if err != nil {
			printData(fmt.Sprintf("cannot marshal request: %s", err))
			continue
		}
		conn.Write(bt)
	}
}

func printData(msg string) {
	msgReceiver <- shared.TerminalData{Message: msg}
}
