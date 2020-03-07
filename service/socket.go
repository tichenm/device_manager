package service

import (
	"github.com/gorilla/websocket"
)

type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Id     string
	Socket *websocket.Conn
	Send   chan []byte
}







