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

	// 구독 메시지 생성
	subscribeMsg := []map[string]interface{}{
		{"ticket": "all-market"},
		{
			"type":            "ticker",
			"codes":           markets,
			"isOnlyRealtime":  true,
		},
	}

	msgBytes, _ := json.Marshal(subscribeMsg)

	if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		log.Fatal("Write error:", err)
	}

	conn.WriteMessage(websocket.TextMessage, msgBytes)


	// 메시지 수신 루프
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(msg, &result); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		// front에 전달
		BroadcastToClients(msg)

		// 필터링 해서 출력
		// fmt.Printf("[Ticker] %s | Price: %v | Volume: %v\n",
		// 	result["code"],
		// 	result["trade_price"],
		// 	result["acc_trade_volume"],
		// )

		// 전체 출력 시 사용
		// fmt.Printf("%+v\n", result) 
	}
}