package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"proxy-broadcast/ws"
)

func main() {
	log.Println("Starting Upbit WebSocket client and frontend server...")

	// 프론트 요청을 받을 WebSocket 핸들러 등록
	http.HandleFunc("/ws", ws.ServeClient)

	// WebSocket 서버 실행 (8081 포트에서 프론트 요청 받음)
	go func() {
		log.Println("WebSocket 서버 시작: ws://localhost:8081/ws")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatal("ListenAndServe error:", err)
		}
	}()

	// Upbit WebSocket 수신기 실행 (서버 입장에서 클라이언트 역할)
	go ws.StartUpbitWebSocket()

	// 종료 대기
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	log.Println("Interrupt received. Shutting down.")
}