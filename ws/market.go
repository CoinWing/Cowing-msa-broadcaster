package ws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 마켓 정보 가져오는 REST API
func GetAllMarkets() []string {
	resp, err := http.Get("https://api.upbit.com/v1/market/all?isDetails=false")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data []map[string]interface{}
	json.Unmarshal(body, &data)

	var markets []string
	for _, m := range data {
		if market, ok := m["market"].(string); ok && strings.HasPrefix(market, "KRW-") {
			markets = append(markets, market)
		}
	}
	fmt.Printf("KRW Markets: %v\n", markets)
	return markets
}