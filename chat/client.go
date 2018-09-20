package main

import (
	"github.com/gorilla/websocket"
)

// clientはチャットを行っている一人のユーザー
type client struct {
	// socketはこのクライアントのwebsocket
	socket *websocket.Conn
	// sendはメッセージが送られるチャネル
	send chan []byte
	// roomはこのクライアントが参加しているチャットルーム
	room *room
}