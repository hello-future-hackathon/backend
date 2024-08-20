package main

import (
    "encoding/json"
    "log"
    "net/http"
    "backend/src/config"
    "backend/src/service/hederaService"
    "backend/src/service/websocketService"
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

    go func() {
        err := api.hederaService.SubscribeToTopic(topicID, func(tm hedera.TopicMessage) {
            log.Printf("Received message: %s\n", string(tm.Contents))
        })
        if err != nil {
            log.Printf("Failed to subscribe to topic %s: %v", topicID.String(), err)
        }
    }()

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(topicID)
}


func main() {
    cfg := config.LoadConfig()

    wsService := websocketService.NewWebSocketService()
    hs, err := hederaService.NewHederaService(cfg, wsService)
    if err != nil {
        log.Fatalf("Failed to create Hedera service: %v", err)
    }

    api := NewAPI(hs)

    http.HandleFunc("/ws", wsService.HandleConnections)
    go wsService.HandleMessages()

    http.HandleFunc("/send-message", api.sendMessageHandler)
    http.HandleFunc("/create-topic", api.CreateTopic)

    existingTopicID := hedera.TopicID{Shard: 0, Realm: 0, Topic: 4703702}
    go func() {
        err := hs.SubscribeToTopic(existingTopicID, func(tm hedera.TopicMessage) {
            log.Printf("Received message: %s\n", string(tm.Contents))
            wsService.Broadcast <- tm.Contents
        })
        if err != nil {
            log.Fatalf("Failed to subscribe to topic %s: %v", existingTopicID.String(), err)
        }
    }()

    log.Fatal(http.ListenAndServe(":8080", nil))
}
