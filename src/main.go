package main

import (
    "encoding/json"
    "log"
    "net/http"
    "backend/src/config"
    "backend/src/hederaService"
    "github.com/hashgraph/hedera-sdk-go/v2"
)

type API struct {
    hederaService *hederaService.HederaService
}

func NewAPI(hs *hederaService.HederaService) *API {
    return &API{hederaService: hs}
}

func (api *API) sendMessageHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        TopicID string `json:"topicId"`
        Message string `json:"message"`
    }

    json.NewDecoder(r.Body).Decode(&req)
    topicID, _ := hedera.TopicIDFromString(req.TopicID)

    txResponse, err := api.hederaService.SendMessage(topicID, req.Message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(txResponse)
}

func main() {
    cfg := config.LoadConfig()

    // Corrected: Handle both return values (HederaService instance and error)
    hs, err := hederaService.NewHederaService(cfg)
    if err != nil {
        log.Fatalf("Failed to create Hedera service: %v", err)
    }

    api := NewAPI(hs)

    http.HandleFunc("/send-message", api.sendMessageHandler)
    // Other handlers for creating topics, subscribing, adding friends, etc.

    log.Fatal(http.ListenAndServe(":8080", nil))
}
