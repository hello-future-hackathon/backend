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

func (api *API) CreateTopic(w http.ResponseWriter, r *http.Request) {
    topicID, err := api.hederaService.CreateTopic()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(topicID)
}

func main() {
    cfg := config.LoadConfig()

    hs, err := hederaService.NewHederaService(cfg)
    if err != nil {
        log.Fatalf("Failed to create Hedera service: %v", err)
    }

    api := NewAPI(hs)

    http.HandleFunc("/send-message", api.sendMessageHandler)
    http.HandleFunc("/create-topic", api.CreateTopic)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
