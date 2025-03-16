package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	// WebSocket服务器的URL
	url := "ws://localhost:8080/ws"

	// 创建一个新的Dialer，它将被用来建立WebSocket连接
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// 使用Dialer建立WebSocket连接
	ws, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Failed to dial WebSocket server:", err)
	}
	defer ws.Close()

	// 发送消息到服务器
	message := "Hello, WebSocket server!"
	err = ws.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Failed to send message:", err)
		return
	}

	// 接收服务器的响应
	for {
		_, response, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error during message read:", err)
			break
		}
		fmt.Printf("Received message from server: %s\n", response)
	}
}
