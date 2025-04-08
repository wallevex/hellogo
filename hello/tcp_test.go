package hello

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"
)

func handleConnection(conn net.Conn) {
	for {
		lengthBuf := make([]byte, 4)
		_, err := conn.Read(lengthBuf)
		if err != nil {
			return
		}

		msgLength := binary.BigEndian.Uint32(lengthBuf) // 读取长度
		msgBuf := make([]byte, msgLength)
		_, err = conn.Read(msgBuf)
		if err != nil {
			return
		}

		fmt.Println("Received:", string(msgBuf))
	}
}

func TestTCPServer(t *testing.T) {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	defer listen.Close()
	fmt.Println("Server listening on port 8080")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn) // 处理每个连接
	}
}

func randString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func sendLoop(ctx context.Context) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8080", 60*time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		msg := randString(10)

		var buf bytes.Buffer
		binary.Write(&buf, binary.BigEndian, int32(len(msg))) // 4字节长度前缀
		buf.WriteString(msg)                                  // 消息体
		conn.Write(buf.Bytes())                               // 发送
		fmt.Printf("Send: %s\n", msg)
	}
}

func TestTCPClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	sendLoop(ctx)
}
