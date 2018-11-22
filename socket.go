package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
)

var ActiveClients = make(map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	webSocket *websocket.Conn
	clientIP  net.Addr
}

func main() {
	addr := flag.String("p", ":3000", "address where the server listen on")
	flag.Parse()

  http.HandleFunc("/socket", handleSocket)

  fmt.Printf("start socket server on %s", *addr)

  err := http.ListenAndServe(*addr, nil)
  if err != nil {
    fmt.Println(err)
  }
}

func handleSocket(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")

  fmt.Println("[socket start]")
  fmt.Println("[ActiveClients]", ActiveClients)
  ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
  if _, ok := err.(websocket.HandshakeError); ok {
    http.Error(w, "Not a websocket handshake", 400)
    return
  } else if err != nil {
    fmt.Println(err)
    return
  }

  client := ws.RemoteAddr()
  sockCli := ClientConn{ws, client}
  addClient(sockCli)
  go readMsg(ws, sockCli)
}

// 添加客户端
func addClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[cc] = 0
	ActiveClientsRWMutex.Unlock()
}

//删除客户端
func deleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, cc)
	ActiveClientsRWMutex.Unlock()
}

//广播信息
func broadcastMessage(messageType int, message []byte) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	for client, _ := range ActiveClients {
		if err := client.webSocket.WriteMessage(messageType, message); err != nil {
			return
		}
	}
}

//接收消息
func readMsg(ws *websocket.Conn, sockCli ClientConn) {
	fmt.Println("[start readMsg]")
	for {
		_, bp, err := ws.ReadMessage()
		if err != nil {
			deleteClient(sockCli)
			return
		}

		fmt.Println("[readMsg]", 1, string(bp))
		fmt.Println("[ActiveClients]", ActiveClients)
		broadcastMessage(1, bp)
	}
}