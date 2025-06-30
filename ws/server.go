package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // CORS 허용 (실 서비스에선 제한 필요)
	},
}

var clients = make(map[*websocket.Conn]bool) // 모든 연결된 클라이언트 저장

func ServeClient(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()


	clients[conn] = true
	log.Println("📡 클라이언트 연결됨")

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("클라이언트 연결 종료:", err)
			delete(clients, conn)
			break
		}
	}
}