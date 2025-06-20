package handler

import (
	"echo-household-budget/internal/infrastructure/middleware"
	"echo-household-budget/internal/usecase"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type (
	ChatMessageTelegraphMethodType string

	ChatMessageTelegraphRequest struct {
		MethodType  ChatMessageTelegraphMethodType `json:"methodType" validate:"required"`
		HouseholdID int                            `json:"householdID" param:"householdID" validate:"required"`
		Message     string                         `json:"message"`
		Limit       int                            `json:"limit" query:"limit"`
		Offset      int                            `json:"offset" query:"offset"` // デバッグ用
	}

	WebSocketChatMessage struct {
		ID          int    `json:"id"`
		UserID      int    `json:"user_id"`
		UserName    string `json:"user_name"`
		Content     string `json:"content"`
		MessageType string `json:"message_type"`
		CreatedAt   string `json:"created_at"`
	}

	ChatMessageTelegraphHandler interface {
		Handle(c echo.Context) error
		WebSocketChat(c echo.Context) error
	}

	chatMessageTelegraphHandler struct {
		wsManager                  *WebSocketManager
		registerChatMessageUsecase usecase.RegisterChatMessageUsecase
		fetchChatMessageUsecase    usecase.FetchChatMessageUsecase
	}
)

const (
	ChatMessageTelegraphMethodTypeFetch    ChatMessageTelegraphMethodType = "fetch"
	ChatMessageTelegraphMethodTypeRegister ChatMessageTelegraphMethodType = "register"
)

// NewChatMessageTelegraphHandler チャットメッセージテレグラフハンドラーのコンストラクタ
func NewChatMessageTelegraphHandler(registerChatMessageUsecase usecase.RegisterChatMessageUsecase, fetchChatMessageUsecase usecase.FetchChatMessageUsecase) ChatMessageTelegraphHandler {
	return &chatMessageTelegraphHandler{
		wsManager:                  GetWebSocketManager(),
		registerChatMessageUsecase: registerChatMessageUsecase,
		fetchChatMessageUsecase:    fetchChatMessageUsecase,
	}
}

// Handle 通常のHTTPリクエストを処理
func (h *chatMessageTelegraphHandler) Handle(c echo.Context) error {
	// 通常のHTTPリクエスト処理（必要に応じて実装）
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Chat message telegraph handler",
	})
}

// WebSocketChat WebSocketを使用したチャット機能
func (h *chatMessageTelegraphHandler) WebSocketChat(c echo.Context) error {
	fmt.Println("=== WebSocketChat メソッドが呼ばれました ===")
	fmt.Println("リクエストURL:", c.Request().URL.String())
	fmt.Println("リクエストメソッド:", c.Request().Method)
	fmt.Println("リクエストヘッダー:", c.Request().Header)

	householdID := c.QueryParam("householdID")
	fmt.Println("householdID:", householdID)
	if householdID == "" {
		fmt.Println("householdIDが空です")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "householdID is required",
		})
	}

	user, ok := middleware.GetUserFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	userID := int(user.ID)
	fmt.Println("userID:", userID)

	// WebSocket接続の確立
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			fmt.Println("CheckOrigin呼ばれました:", r.Header.Get("Origin"))
			return true
		},
	}

	fmt.Println("WebSocket接続をアップグレード中...")
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println("WebSocket接続エラー:", err)
		return fmt.Errorf("WebSocket接続の確立に失敗しました: %w", err)
	}
	fmt.Println("WebSocket接続が確立されました")

	// クライアントを管理に追加
	h.wsManager.AddClient(conn)
	defer h.wsManager.RemoveClient(conn)

	// 接続確認メッセージを送信
	welcomeMsg := WebSocketChatMessage{
		ID:          1,
		UserID:      1,
		UserName:    "システム",
		Content:     "チャットに接続しました",
		MessageType: "system",
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	welcomeJSON, _ := json.Marshal(welcomeMsg)
	fmt.Println("ウェルカムメッセージを送信:", string(welcomeJSON))
	conn.WriteMessage(websocket.TextMessage, welcomeJSON)

	householdIDInt, err := strconv.Atoi(householdID)
	if err != nil {
		fmt.Println("householdIDの変換エラー:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid householdID format",
		})
	}

	err = h.fetchChatMessage(ChatMessageTelegraphRequest{
		MethodType:  ChatMessageTelegraphMethodTypeFetch,
		HouseholdID: householdIDInt,
		Limit:       10,
		Offset:      0,
	}, userID)
	if err != nil {
		fmt.Println("チャットメッセージ取得エラー:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch chat messages",
		})
	}

	// メッセージループ
	fmt.Println("メッセージループを開始します")
	return h.handleChatMessageLoop(conn, householdID, userID)
}

// handleChatMessageLoop チャットメッセージループを処理
func (h *chatMessageTelegraphHandler) handleChatMessageLoop(conn *websocket.Conn, householdID string, userID int) error {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("チャットメッセージ読み取りエラー:", err)
			break
		}

		log.Println("受信チャットメッセージ:", string(msg))

		// メッセージの処理
		if err := h.processChatMessage(msg, householdID, userID); err != nil {
			log.Printf("チャットメッセージ処理エラー: %v", err)
			// エラーが発生しても他のクライアントには影響させない
			continue
		}
	}

	return nil
}

// processChatMessage チャットメッセージを処理
func (h *chatMessageTelegraphHandler) processChatMessage(msg []byte, householdID string, userID int) error {
	var request ChatMessageTelegraphRequest
	if err := json.Unmarshal(msg, &request); err != nil {
		log.Println("チャットメッセージJSONデコードエラー:", err)
		return fmt.Errorf("JSONデコードエラー: %w", err)
	}

	// websocketの処理区分による処理分岐
	switch request.MethodType {
	case ChatMessageTelegraphMethodTypeFetch:
		return h.fetchChatMessage(request, userID)
	case ChatMessageTelegraphMethodTypeRegister:
		return h.registerChatMessage(request, userID)
	}

	return nil
}

func (h *chatMessageTelegraphHandler) fetchChatMessage(request ChatMessageTelegraphRequest, userID int) error {
	// limit, offsetに則って、取得したデータをブロードキャストする
	// チャットメッセージを作成
	// メッセージをJSONに変換
	fetchChatMessageInput := usecase.FetchChatMessageInput{
		HouseholdID: request.HouseholdID,
		Limit:       request.Limit,
		Offset:      request.Offset,
	}

	fetchChatMessageOutput, err := h.fetchChatMessageUsecase.Execute(fetchChatMessageInput)
	if err != nil {
		return fmt.Errorf("チャットメッセージ取得エラー: %w", err)
	}

	for _, chatMessage := range fetchChatMessageOutput.ChatMessages {
		chatmessage := WebSocketChatMessage{
			ID:          chatMessage.ID,
			UserID:      chatMessage.UserID,
			UserName:    chatMessage.User.Name,
			Content:     chatMessage.Content,
			MessageType: string(chatMessage.MessageType),
			CreatedAt:   time.Now().Format(time.RFC3339),
		}

		messageJSON, err := json.Marshal(chatmessage)
		if err != nil {
			log.Printf("チャットメッジJSONマーシャリングエラー: %v", err)
			return fmt.Errorf("JSONマーシャリングに失敗しました: %w", err)
		}

		h.wsManager.BroadcastToAll(messageJSON)
	}

	return nil
}

func (h *chatMessageTelegraphHandler) registerChatMessage(request ChatMessageTelegraphRequest, userID int) error {
	input := usecase.RegisterChatMessageInput{
		HouseholdID: request.HouseholdID,
		UserID:      userID,
		Message:     request.Message,
	}

	aiChatReplyMessage, err := h.registerChatMessageUsecase.Execute(input)
	if err != nil {
		return fmt.Errorf("チャットメッセージ登録エラー: %w", err)
	}

	chatMessage := WebSocketChatMessage{
		ID:          aiChatReplyMessage.ID,
		UserID:      aiChatReplyMessage.UserID,
		UserName:    "AI",
		Content:     aiChatReplyMessage.Content,
		MessageType: string(aiChatReplyMessage.MessageType),
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	// メッセージをJSONに変換
	messageJSON, err := json.Marshal(chatMessage)
	if err != nil {
		log.Printf("チャットメッセージJSONマーシャリングエラー: %v", err)
		return fmt.Errorf("JSONマーシャリングに失敗しました: %w", err)
	}

	// 全クライアントにブロードキャスト
	h.wsManager.BroadcastToAll(messageJSON)
	return nil
}
