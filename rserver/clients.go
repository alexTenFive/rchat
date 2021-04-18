package main

import (
	"chat/shared"
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

var (
	clients      map[string]*net.TCPConn
	clientsMutex sync.RWMutex
)

func init() {
	clients = make(map[string]*net.TCPConn)
}
func addNewClient(name string, conn *net.TCPConn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	if _, ok := clients[name]; ok {
		return
	}
	clients[name] = conn
}

func getClients() map[string]*net.TCPConn {
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()
	return clients
}

func removeFromClients(name string) {
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()
	delete(clients, name)
}

func sendToClientsExcept(message, exceptName string) {
	for _, conn := range getClients() {
		// if name == exceptName {
		// 	continue
		// }
		msg := &shared.Message{
			Time:    time.Now(),
			Name:    exceptName,
			Message: message,
		}
		bt, err := json.Marshal(msg)
		if err != nil {
			log.Printf("cannot marshal message to client [%s]: %s\n", message, err)
			continue
		}
		conn.Write(bt)
	}
}
