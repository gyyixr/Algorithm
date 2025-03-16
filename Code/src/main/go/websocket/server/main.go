package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// 定义WebSocket升级处理器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有CORS请求，生产环境中应进行适当的限制
		return true
	},
}

// 处理WebSocket连接的函数
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	for {
		// 读取客户端发送的消息
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message read:", err)
			break
		}

		// 这里可以处理接收到的消息，例如打印或转发
		log.Printf("Received message: %s", message)

		// 向客户端发送消息
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println("Error during message write:", err)
			break
		}
	}
}

func main() {
	// 启动HTTP服务器，并为WebSocket路径设置处理器
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Starting WebSocket server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
