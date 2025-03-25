package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type AccountTxPayload struct {
	AccountID string `json:"account_id"`
	Lt        int64  `json:"lt"`
	TxHash    string `json:"tx_hash"`
}

type MempoolEvent struct {
	EventType string `json:"event_type"`
	Boc       string `json:"boc"`
}

var (
	mu              sync.Mutex
	hashes          = make(map[string]int)
	messagesCounter = 0
)

func webhookMempool(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	fmt.Printf("req: %v\n", req.Method)
	//for name, header := range req.Header {
	//	fmt.Printf("header: %v -> %#v\n", name, header)
	//
	//}
	var payload MempoolEvent
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := sha256.New()
	hash.Write([]byte(payload.Boc))
	payloadHash := hash.Sum(nil)
	payloadHex := hex.EncodeToString(payloadHash)

	w.WriteHeader(http.StatusOK)

	mu.Lock()
	defer mu.Unlock()

	hashes[payloadHex] += 1
	messagesCounter += 1

	//value := hashes[payloadHex]
	//if value > 2 {
	//	//fmt.Printf("payload hex %v -> %v\n", payloadHex, value)
	//}
	if messagesCounter%100 == 0 {
		x := 0
		for _, value := range hashes {
			if value > 2 {
				x += 1
			}
		}
		percent := float64(x) / float64(len(hashes)) * 100.0
		fmt.Printf("%0.2f\n", percent)
	}

	if len(hashes) > 100_000 {
		hashes = make(map[string]int)
	}
}

func webhook(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	fmt.Printf("req: %v\n", req.Method)
	//for name, header := range req.Header {
	//	fmt.Printf("header: %v -> %#v\n", name, header)
	//
	//}
	var payload AccountTxPayload
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("%v %v %v\n", payload.AccountID, payload.Lt, payload.TxHash)

}

func main() {
	http.HandleFunc("/webhook", webhook)
	http.HandleFunc("/webhook-mempool", webhookMempool)
	http.ListenAndServe(":8090", nil)
}
