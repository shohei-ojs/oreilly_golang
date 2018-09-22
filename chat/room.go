package main

import (
	"log"
)

type room struct {
	// forwardは他のクライアントへ転送するメッセージ保持チャネル
	forward chan []byte
	// joinはチャットルームに参加しようとしているクライアントのためのチャネル
	join chan *client
	// leaveはチャットルームから退出しようとしているクライアントのためのチャネル
	leave chan *client
	// clientsは在室している全てのクライアント
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool)
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true
		case client := <-r.leave:
			// 退出
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// 全てのクライアントにメッセージ送信
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージを送信
				default:
					// 送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)
var upgrader = &websocket.Upgrader{ReadBufferSize:
		socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWrite, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client
	defer func() {r.leave <- client}()
	go client.write()
	client.read()
}