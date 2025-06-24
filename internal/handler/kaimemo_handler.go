//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package handler

import (
	"echo-household-budget/internal/infrastructure/middleware"
	"echo-household-budget/internal/model"
	"echo-household-budget/internal/usecase"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	domainmodel "echo-household-budget/internal/domain/model"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type kaimemoHandler struct {
	service         usecase.KaimemoService
	shoppingUsecase usecase.ShoppingUsecase
	wsManager       *WebSocketManager
}

// WebSocketMessageProcessor WebSocketメッセージを処理する構造体
type WebSocketMessageProcessor struct {
	shoppingUsecase usecase.ShoppingUsecase
	wsManager       *WebSocketManager
}

// NewWebSocketMessageProcessor WebSocketMessageProcessorのコンストラクタ
func NewWebSocketMessageProcessor(shoppingUsecase usecase.ShoppingUsecase, wsManager *WebSocketManager) *WebSocketMessageProcessor {
	return &WebSocketMessageProcessor{
		shoppingUsecase: shoppingUsecase,
		wsManager:       wsManager,
	}
}

// ProcessMessage メッセージを処理
func (p *WebSocketMessageProcessor) ProcessMessage(msg []byte, householdID uint) error {
	var request model.TelegraphRequest
	if err := json.Unmarshal(msg, &request); err != nil {
		log.Println("JSONデコードエラー:", err)
		return fmt.Errorf("JSONデコードエラー: %w", err)
	}

	log.Println("処理中のリクエスト:", request)
	spew.Dump(request)

	switch request.MethodType {
	case model.CreateKaimemo:
		return p.handleCreateShopping(request, householdID)
	case model.RemoveKaimemo:
		return p.handleDeleteShopping(request)
	default:
		return fmt.Errorf("未対応のメソッドタイプ: %s", request.MethodType)
	}
}

// handleCreateShopping 買い物メモの作成を処理
func (p *WebSocketMessageProcessor) handleCreateShopping(request model.TelegraphRequest, householdID uint) error {
	if request.HouseholdBookID == nil || request.Tag == nil || request.Name == nil {
		return fmt.Errorf("必須パラメータが不足しています")
	}

	shopping := domainmodel.NewShoppingMemo(
		domainmodel.HouseHoldID(*request.HouseholdBookID),
		domainmodel.CategoryID(*request.Tag),
		*request.Name,
		"",
	)

	if err := p.shoppingUsecase.CreateShopping(shopping); err != nil {
		log.Printf("買い物メモ作成エラー: %v", err)
		return fmt.Errorf("買い物メモの作成に失敗しました: %w", err)
	}

	return nil
}

// handleDeleteShopping 買い物メモの削除を処理
func (p *WebSocketMessageProcessor) handleDeleteShopping(request model.TelegraphRequest) error {
	if request.ID == nil {
		return fmt.Errorf("削除対象のIDが指定されていません")
	}

	if err := p.shoppingUsecase.DeleteShopping(domainmodel.ShoppingID(*request.ID)); err != nil {
		log.Printf("買い物メモ削除エラー: %v", err)
		return fmt.Errorf("買い物メモの削除に失敗しました: %w", err)
	}

	return nil
}

// broadcastUpdatedData 更新されたデータを全クライアントにブロードキャスト
func (p *WebSocketMessageProcessor) broadcastUpdatedData(householdID uint) error {
	res, err := p.shoppingUsecase.FetchShopping(domainmodel.HouseHoldID(householdID))
	if err != nil {
		log.Printf("データ取得エラー: %v", err)
		return fmt.Errorf("データの取得に失敗しました: %w", err)
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		log.Printf("JSONマーシャリングエラー: %v", err)
		return fmt.Errorf("JSONマーシャリングに失敗しました: %w", err)
	}

	p.wsManager.BroadcastToAll(resJSON)
	return nil
}

// WebsocketTelegraph implements KaimemoHandler.
func (k *kaimemoHandler) WebsocketTelegraph(c echo.Context) error {
	tempUserID := c.QueryParam("tempUserID")
	if tempUserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "tempUserID is required",
		})
	}

	tempUserIDUint, err := strconv.ParseUint(tempUserID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid tempUserID format",
		})
	}

	// WebSocket接続の確立
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("WebSocket接続の確立に失敗しました: %w", err)
	}

	// クライアントを管理に追加
	k.wsManager.AddClient(conn)
	defer k.wsManager.RemoveClient(conn)

	// 初期データの送信
	if err := k.sendInitialData(conn, uint(tempUserIDUint)); err != nil {
		log.Printf("初期データ送信エラー: %v", err)
		return err
	}

	// メッセージプロセッサーの初期化
	processor := NewWebSocketMessageProcessor(k.shoppingUsecase, k.wsManager)

	// メッセージループ
	return k.handleMessageLoop(conn, processor, uint(tempUserIDUint))
}

// sendInitialData 初期データを送信
func (k *kaimemoHandler) sendInitialData(conn *websocket.Conn, householdID uint) error {
	res, err := k.shoppingUsecase.FetchShopping(domainmodel.HouseHoldID(householdID))
	if err != nil {
		return fmt.Errorf("初期データの取得に失敗しました: %w", err)
	}

	log.Println("初期データを送信:", res)
	spew.Dump(res)

	resJSON, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("初期データのJSONマーシャリングに失敗しました: %w", err)
	}

	return conn.WriteMessage(websocket.TextMessage, resJSON)
}

// handleMessageLoop メッセージループを処理
func (k *kaimemoHandler) handleMessageLoop(conn *websocket.Conn, processor *WebSocketMessageProcessor, householdID uint) error {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("メッセージ読み取りエラー:", err)
			break
		}

		log.Println("受信メッセージ:", string(msg))

		// メッセージの処理
		if err := processor.ProcessMessage(msg, householdID); err != nil {
			log.Printf("メッセージ処理エラー: %v", err)
			// エラーが発生しても他のクライアントには影響させない
			continue
		}

		// 更新されたデータをブロードキャスト
		if err := processor.broadcastUpdatedData(householdID); err != nil {
			log.Printf("ブロードキャストエラー: %v", err)
			// エラーが発生しても接続は維持
		}
	}

	return nil
}

// CreateKaimemoAmount implements KaimemoHandler.
func (k *kaimemoHandler) CreateKaimemoAmount(c echo.Context) error {
	req := model.CreateKaimemoAmountRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := k.service.CreateKaimemoAmount(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create kaimemo amount",
		})
	}

	return c.NoContent(http.StatusCreated)
}

// FetchKaimemoSummaryRecord implements KaimemoHandler.
func (k *kaimemoHandler) FetchKaimemoSummaryRecord(c echo.Context) error {
	ctx := c.Request().Context()
	user, ok := ctx.Value(middleware.UserKey).(*domainmodel.UserAccount)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user",
		})
	}
	spew.Dump(user)

	tempUserID := c.QueryParam("tempUserID")
	if tempUserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "TempUserID is required",
		})
	}

	res, err := k.service.FetchKaimemoSummaryRecord(tempUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch kaimemo summary record",
		})
	}

	return c.JSON(http.StatusOK, res)
}

// RemoveKaimemoAmount implements KaimemoHandler.
func (k *kaimemoHandler) RemoveKaimemoAmount(c echo.Context) error {
	req := model.RemoveKaimemoAmountRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ID is required",
		})
	}

	if err := k.service.RemoveKaimemoAmount(id, req.TempUserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to remove kaimemo",
		})
	}

	return c.NoContent(http.StatusOK)
}

// CreateKaimemo implements KaimemoHandler.
func (k *kaimemoHandler) CreateKaimemo(c echo.Context) error {
	req := model.CreateKaimemoRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	if err := k.service.CreateKaimemo(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create kaimemo",
		})
	}

	return c.NoContent(http.StatusCreated)
}

// FetchKaimemo implements KaimemoHandler.
func (k *kaimemoHandler) FetchKaimemo(c echo.Context) error {
	tempUserID := c.QueryParam("tempUserID")
	if tempUserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "TempUserID is required",
		})
	}

	res, err := k.service.FetchKaimemo(tempUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch kaimemo",
		})
	}

	return c.JSON(http.StatusOK, res)
}

// RemoveKaimemo implements KaimemoHandler.
func (k *kaimemoHandler) RemoveKaimemo(c echo.Context) error {
	req := model.RemoveKaimemoRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request",
		})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ID is required",
		})
	}

	if err := k.service.RemoveKaimemo(id, req.TempUserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to remove kaimemo",
		})
	}

	return c.NoContent(http.StatusOK)
}

type KaimemoHandler interface {
	WebsocketTelegraph(c echo.Context) error
	FetchKaimemo(c echo.Context) error
	CreateKaimemo(c echo.Context) error
	RemoveKaimemo(c echo.Context) error
	FetchKaimemoSummaryRecord(c echo.Context) error
	CreateKaimemoAmount(c echo.Context) error
	RemoveKaimemoAmount(c echo.Context) error
}

func NewKaimemoHandler(service usecase.KaimemoService, shoppingUsecase usecase.ShoppingUsecase) KaimemoHandler {
	return &kaimemoHandler{service: service, shoppingUsecase: shoppingUsecase, wsManager: GetWebSocketManager()}
}
