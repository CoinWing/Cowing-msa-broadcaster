package ws

import (
	"log"
)

func BroadcastToClients(messageType int, message []byte) {
	for conn := range clients {
		// Upbit로부터 받은 메시지 타입을 그대로 사용
		// BinaryMessage(2) 또는 TextMessage(1)를 그대로 전달
		err := conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("❌ 클라이언트 전송 실패:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
