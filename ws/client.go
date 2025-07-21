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
			"type":           "ticker",
			"codes":          markets,
			"isOnlyRealtime": true,
		},
	}

	// orderbook 구독 메시지 생성
	orderbookSubscribeMsg := []map[string]interface{}{
		{"ticket": "orderbook-subscription"},
		{
			"type":  "orderbook",
			"codes": markets,
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

	// 메시지 수신 루프
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		// Upbit로부터 받은 Binary 데이터를 그대로 브로드캐스트
		// messageType과 msg를 모두 전달하여 Binary 형태 유지
		BroadcastToClients(messageType, msg)

		// 로깅을 위한 파싱 (선택적)
		if messageType == websocket.TextMessage {
			var result map[string]interface{}
			if err := json.Unmarshal(msg, &result); err != nil {
				log.Println("Unmarshal error:", err)
				continue
			}
			// 로그 출력 (선택적)
			// log.Printf("[%s] %s | Price: %v | Volume: %v\n",
			// 	result["type"],
			// 	result["code"],
			// 	result["trade_price"],
			// 	result["acc_trade_volume"],
			// )
		}
	}
}
