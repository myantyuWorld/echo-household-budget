package handler

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketManager WebSocket接続を管理する構造体
type WebSocketManager struct {
	clients map[*websocket.Conn]bool
	mutex   sync.RWMutex
}

var (
	wsManagerInstance *WebSocketManager
	wsManagerOnce     sync.Once
)

// GetWebSocketManager シングルトンインスタンスを取得
func GetWebSocketManager() *WebSocketManager {
	wsManagerOnce.Do(func() {
		wsManagerInstance = &WebSocketManager{
			clients: make(map[*websocket.Conn]bool),
		}
	})
	return wsManagerInstance
}

// NewWebSocketManager WebSocketManagerのコンストラクタ（後方互換性のため）
func NewWebSocketManager() *WebSocketManager {
	return GetWebSocketManager()
}

// AddClient クライアントを追加
func (wm *WebSocketManager) AddClient(conn *websocket.Conn) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	wm.clients[conn] = true
	log.Printf("クライアントが追加されました。現在の接続数: %d", len(wm.clients))
}

// RemoveClient クライアントを削除
func (wm *WebSocketManager) RemoveClient(conn *websocket.Conn) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	delete(wm.clients, conn)
	conn.Close()
	log.Printf("クライアントが削除されました。現在の接続数: %d", len(wm.clients))
}

// BroadcastToAll 全クライアントにメッセージをブロードキャスト
func (wm *WebSocketManager) BroadcastToAll(message []byte) {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	clientCount := len(wm.clients)
	if clientCount == 0 {
		log.Println("ブロードキャスト: 接続中のクライアントがありません")
		return
	}

	log.Printf("ブロードキャスト: %d個のクライアントにメッセージを送信", clientCount)

	for client := range wm.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("ブロードキャストエラー: %v", err)
			// エラーが発生したクライアントを削除
			go wm.RemoveClient(client)
		}
	}
}

// GetClientCount 接続中のクライアント数を取得
func (wm *WebSocketManager) GetClientCount() int {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()
	return len(wm.clients)
}

// IsClientConnected 指定されたクライアントが接続中かどうかを確認
func (wm *WebSocketManager) IsClientConnected(conn *websocket.Conn) bool {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()
	_, exists := wm.clients[conn]
	return exists
}

// BroadcastToSpecific 特定のクライアントにのみメッセージを送信
func (wm *WebSocketManager) BroadcastToSpecific(conn *websocket.Conn, message []byte) error {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	if _, exists := wm.clients[conn]; !exists {
		return fmt.Errorf("指定されたクライアントは接続されていません")
	}

	return conn.WriteMessage(websocket.TextMessage, message)
}
