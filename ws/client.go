package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func StartUpbitWebSocket() {
	url := "wss://api.upbit.com/websocket/v1"
	dialer := websocket.Dialer{
		Subprotocols: []string{"binary"},
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal("WebSocket dial error:", err)
	}
	defer conn.Close()

	// 모든 마켓 받아오기
	markets := GetAllMarkets()

	// ticker 구독 메시지 생성
	tickerSubscribeMsg := []map[string]interface{}{
		{"ticket": "ticker-subscription"},
		{
			"type":             "ticker",
			"codes":            markets,
			"is_only_realtime": true,
		},
	}

	// orderbook 구독 메시지 생성
	orderbookSubscribeMsg := []map[string]interface{}{
		{"ticket": "orderbook-subscription"},
		{
			"type":             "orderbook",
			"codes":            markets,
			"is_only_realtime": true,
		},
	}

	// ticker 구독
	tickerMsgBytes, _ := json.Marshal(tickerSubscribeMsg)
	if err := conn.WriteMessage(websocket.TextMessage, tickerMsgBytes); err != nil {
		log.Fatal("Ticker subscribe error:", err)
	}

	// orderbook 구독
	orderbookMsgBytes, _ := json.Marshal(orderbookSubscribeMsg)
	if err := conn.WriteMessage(websocket.TextMessage, orderbookMsgBytes); err != nil {
		log.Fatal("Orderbook subscribe error:", err)
	}

	log.Println("📡 Upbit WebSocket 구독 완료 (ticker + orderbook)")

	// 메시지 수신 루프 - Upbit로부터 실시간 데이터를 계속 수신
	for {
		// WebSocket으로부터 메시지 읽기
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		// 받은 데이터를 모든 클라이언트에게 브로드캐스트
		BroadcastToClients(messageType, msg)
	}
}
