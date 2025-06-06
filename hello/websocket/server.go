package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// 升级 HTTP 连接为 WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域连接，实际部署需更严格配置
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// 启动定时器，每隔 1 秒发送数据
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			// 构造你要发的数据（比如时间戳）
			msg := fmt.Sprintf("当前时间: %s", t.Format("15:04:05"))
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("WriteMessage error:", err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("服务启动于 http://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
