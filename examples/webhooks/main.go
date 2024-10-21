package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AccountTxPayload struct {
	AccountID string `json:"account_id"`
	Lt        int64  `json:"lt"`
	TxHash    string `json:"tx_hash"`
}

func webhook(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var payload AccountTxPayload
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("%v\n", payload)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", webhook)
	http.ListenAndServe(":8092", nil)
}
