package websocketService

import (
    "net/http"
    "github.com/gorilla/websocket"
    "log"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type WebSocketService struct {
    Clients   map[*websocket.Conn]bool
    Broadcast chan []byte
}

func NewWebSocketService() *WebSocketService {
    return &WebSocketService{
        Clients:   make(map[*websocket.Conn]bool),
        Broadcast: make(chan []byte),
    }
}

func (ws *WebSocketService) HandleConnections(w http.ResponseWriter, r *http.Request) {
    wsConn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket Upgrade Error: %v\n", err)
        return
    }
    defer wsConn.Close()

    ws.Clients[wsConn] = true

    for {
        messageType, message, err := wsConn.ReadMessage()
        if err != nil {
            log.Printf("WebSocket Read Error: %v\n", err)
            delete(ws.Clients, wsConn)
            break
        }

        log.Printf("Received message: %s\n", string(message))
        
        if messageType == websocket.TextMessage {
            ws.Broadcast <- message
        }
    }
}

func (ws *WebSocketService) HandleMessages() {
    for {
        msg := <-ws.Broadcast
        for client := range ws.Clients {
            err := client.WriteMessage(websocket.TextMessage, msg)
            if err != nil {
                log.Printf("WebSocket Write Error: %v\n", err)
                client.Close()
                delete(ws.Clients, client)
            }
        }
    }
}
