package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket 연결 업그레이더 설정
var upgrader = websocket.Upgrader{
	// CORS(Cross-Origin Resource Sharing) 정책 설정
	CheckOrigin: func(r *http.Request) bool {
		return true // 모든 도메인에서의 접근을 허용 (개발 환경용)
	},
}

// 연결된 모든 클라이언트를 관리하는 전역 맵
// key: WebSocket 연결 객체, value: 연결 상태 (true/false)
var clients = make(map[*websocket.Conn]bool)

// HTTP 요청을 WebSocket으로 업그레이드하고 클라이언트 연결을 처리하는 핸들러
func HttpToWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// HTTP 연결을 WebSocket으로 업그레이드
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 업그레이드 실패:", err)
		return
	}
	// 함수 종료 시 연결 정리
	defer conn.Close()

	// 새로운 클라이언트를 clients 맵에 등록
	clients[conn] = true
	log.Println("📡 새 클라이언트가 연결되었습니다")

	// 클라이언트 연결 유지 및 메시지 수신 대기 루프
	for {
		// 클라이언트로부터 메시지 읽기 (실제로는 연결 상태 확인용)
		_, _, err := conn.ReadMessage()
		if err != nil {
			// 클라이언트 연결이 끊어진 경우
			log.Println("클라이언트 연결 종료:", err)
			// clients 맵에서 해당 연결 제거 (메모리 누수 방지)
			delete(clients, conn)
			break // 루프 종료하여 함수 종료
		}
		// 참고: 실제로는 클라이언트가 메시지를 보내지 않고,
		// 서버에서 Upbit 데이터를 일방적으로 브로드캐스트함
	}
}
