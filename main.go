package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

// クライアントの接続を管理
var clients = struct {
	sync.Mutex
	conns map[*websocket.Conn]bool
}{conns: make(map[*websocket.Conn]bool)}

// メッセージキュー（チャネル）
var messageQueue = make(chan string, 100)

// メッセージをブロードキャストするゴルーチン
// ゴルーチン：軽量スレッド
func startBroadcastWorker() {
	for msg := range messageQueue {
		clients.Lock()
		for conn := range clients.conns {
			err := websocket.Message.Send(conn, msg)
			if err != nil {
				log.Println("Error sending message:", err)
				conn.Close()
				delete(clients.conns, conn)
			}
		}
		clients.Unlock()
	}
}

func handleWebSocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		// クライアントを登録
		clients.Lock()
		clients.conns[ws] = true
		clients.Unlock()

		defer func() {
			// クライアントを削除
			clients.Lock()
			delete(clients.conns, ws)
			clients.Unlock()
			ws.Close()
		}()

		for {
			var msg string
			// クライアントからのメッセージを受信
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				log.Println("Error receiving message:", err)
				break
			}

			// メッセージをキューに追加
			messageQueue <- fmt.Sprintf("Client: %s", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/", "public")
	e.GET("/ws", handleWebSocket)

	// ブロードキャスト用のゴルーチンを起動
	go startBroadcastWorker()

	e.Logger.Fatal(e.Start(":8080"))
}
