package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Tüm origin'lere izin ver (Geliştirme aşaması için)
		},
	}
	clients = make(map[*websocket.Conn]bool) // Tüm aktif istemcileri saklar
	lock    sync.Mutex                       // Eşzamanlı erişimi kontrol etmek için
)

// handleMessage tüm bağlı istemcilere mesajı yayınlar
func handleMessage(messageType int, message []byte) {
	lock.Lock()
	defer lock.Unlock()
	for client := range clients {
		if err := client.WriteMessage(messageType, message); err != nil {
			log.Printf("Error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Yeni istemciyi kaydet
	lock.Lock()
	clients[conn] = true
	lock.Unlock()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("Received: %s\n", message)
		handleMessage(messageType, message)
	}

	// İstemci bağlantıyı kapattığında
	lock.Lock()
	delete(clients, conn)
	lock.Unlock()
}

func main() {
	http.HandleFunc("/ws", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
