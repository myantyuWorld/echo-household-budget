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

	ChatMessageTelegraphResponse struct {
		ID          int    `json:"id"`
		UserID      int    `json:"user_id"`
		UserName    string `json:"user_name"`
		Content     string `json:"content"`
		MessageType string `json:"message_type"`
		CreatedAt   string `json:"created_at"`
	}

	ChatMessageTelegraphHandler interface {
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

// WebSocketChat WebSocketを使用したチャット機能
func (h *chatMessageTelegraphHandler) WebSocketChat(c echo.Context) error {
	// パラメータの検証
	householdID, userID, err := h.validateWebSocketRequest(c)
	if err != nil {
		return err
	}

	// WebSocket接続の確立
	conn, err := h.establishWebSocketConnection(c)
	if err != nil {
		return err
	}
	defer h.wsManager.RemoveClient(conn)

	// クライアントを管理に追加
	h.wsManager.AddClient(conn)

	// 初期化処理
	if err := h.initializeChatSession(conn, householdID, userID); err != nil {
		return err
	}

	// メッセージループを開始
	return h.handleChatMessageLoop(conn, strconv.Itoa(householdID), userID)
}

// validateWebSocketRequest WebSocketリクエストのパラメータを検証する
func (h *chatMessageTelegraphHandler) validateWebSocketRequest(c echo.Context) (int, int, error) {
	householdID := c.QueryParam("householdID")
	if householdID == "" {
		return 0, 0, c.JSON(http.StatusBadRequest, map[string]string{
			"error": "householdID is required",
		})
	}

	user, ok := middleware.GetUserFromContext(c.Request().Context())
	if !ok {
		return 0, 0, c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	householdIDInt, err := strconv.Atoi(householdID)
	if err != nil {
		return 0, 0, c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid householdID format",
		})
	}

	return householdIDInt, int(user.ID), nil
}

// establishWebSocketConnection WebSocket接続を確立する
func (h *chatMessageTelegraphHandler) establishWebSocketConnection(c echo.Context) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return nil, fmt.Errorf("WebSocket接続の確立に失敗しました: %w", err)
	}

	return conn, nil
}

// initializeChatSession チャットセッションを初期化する
func (h *chatMessageTelegraphHandler) initializeChatSession(conn *websocket.Conn, householdID, userID int) error {
	// 接続確認メッセージを送信
	welcomeMsg := ChatMessageTelegraphResponse{
		ID:          1,
		UserID:      1,
		UserName:    "システム",
		Content:     "チャットに接続しました",
		MessageType: "system",
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	welcomeJSON, err := json.Marshal(welcomeMsg)
	if err != nil {
		return fmt.Errorf("ウェルカムメッセージのマーシャリングに失敗しました: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, welcomeJSON); err != nil {
		return fmt.Errorf("ウェルカムメッセージの送信に失敗しました: %w", err)
	}

	// 既存のチャットメッセージを取得
	err = h.fetchChatMessage(ChatMessageTelegraphRequest{
		MethodType:  ChatMessageTelegraphMethodTypeFetch,
		HouseholdID: householdID,
		Limit:       10,
		Offset:      0,
	}, userID)
	if err != nil {
		return fmt.Errorf("チャットメッセージの取得に失敗しました: %w", err)
	}

	return nil
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
		chatmessage := ChatMessageTelegraphResponse{
			ID:          chatMessage.ID,
			UserID:      chatMessage.UserID,
			UserName:    chatMessage.User.Name,
			Content:     chatMessage.Content,
			MessageType: string(chatMessage.MessageType),
			CreatedAt:   chatMessage.CreatedAt.Format(time.RFC3339),
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

	chatMessage := ChatMessageTelegraphResponse{
		ID:          aiChatReplyMessage.ID,
		UserID:      aiChatReplyMessage.UserID,
		UserName:    "AI",
		Content:     aiChatReplyMessage.Content,
		MessageType: string(aiChatReplyMessage.MessageType),
		CreatedAt:   aiChatReplyMessage.CreatedAt.Format(time.RFC3339),
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
