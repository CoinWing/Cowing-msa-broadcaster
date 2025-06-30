package ws

import (
	"log"
)

func BroadcastToClients(message []byte) {
	for conn := range clients {
		err := conn.WriteMessage(1, message) // 1 = TextMessage
		if err != nil {
			log.Println("❌ 클라이언트 전송 실패:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}