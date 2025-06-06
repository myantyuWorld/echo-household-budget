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
}

var clients = make(map[*websocket.Conn]bool)

// FYI. GoでWebSocketを使いチャットサーバー構築 | https://qiita.com/TetsuyaFukunaga/items/4c83a8dedd34e65ffbdc
// WebsocketTelegraph implements KaimemoHandler.
func (k *kaimemoHandler) WebsocketTelegraph(c echo.Context) error {
	tempUserID := c.QueryParam("tempUserID")
	if tempUserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "tempUserID is required",
		})
	}

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	clients[conn] = true
	defer conn.Close()

	tempUserIDUint, err := strconv.ParseUint(tempUserID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid tempUserID format",
		})
	}

	res, err := k.shoppingUsecase.FetchShopping(domainmodel.HouseHoldID(uint(tempUserIDUint)))
	fmt.Println("res, err := k.shoppingUsecase.FetchShopping(domainmodel.HouseHoldID(uint(tempUserIDUint)))")
	spew.Dump(res)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch kaimemo",
		})
	}
	resJSON, _ := json.Marshal(res)
	conn.WriteMessage(websocket.TextMessage, resJSON)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("読み取りエラー:", err)
			break
		}
		fmt.Println("msg", msg)

		// TODO : 処理区分を渡して、それに応じて登録・削除を分ける必要がある
		var request model.TelegraphRequest
		if err := json.Unmarshal(msg, &request); err != nil {
			log.Println("JSONデコードエラー:", err)
			continue
		}
		fmt.Println("request", request)
		spew.Dump(request)

		if request.MethodType == "1" {
			shopping := domainmodel.NewShoppingMemo(domainmodel.HouseHoldID(*request.HouseholdBookID), domainmodel.CategoryID(*request.Tag), *request.Name, "")

			if err := k.shoppingUsecase.CreateShopping(shopping); err != nil {
				fmt.Println("err", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to create shopping",
				})
			}

		} else if request.MethodType == "2" {
			// TODO : shoppingUsecaseで削除する
			if err := k.shoppingUsecase.DeleteShopping(domainmodel.ShoppingID(*request.ID)); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to delete shopping",
				})
			}

			// if err := k.service.RemoveKaimemo(*request.ID, tempUserID); err != nil {
			// 	return c.JSON(http.StatusInternalServerError, map[string]string{
			// 		"error": "Failed to remove kaimemo",
			// 	})
			// }
		}

		res, err := k.shoppingUsecase.FetchShopping(domainmodel.HouseHoldID(tempUserIDUint))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch shopping",
			})
		}
		resJSON, _ := json.Marshal(res)

		// 全クライアントにメッセージをブロードキャスト
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, resJSON); err != nil {
				log.Printf("ブロードキャストエラー: %v", err)
				delete(clients, client)
				client.Close()
			}
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
	return &kaimemoHandler{service: service, shoppingUsecase: shoppingUsecase}
}
