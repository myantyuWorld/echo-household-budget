package handler

import (
	domainmodel "echo-household-budget/internal/domain/model"
	domainservice "echo-household-budget/internal/domain/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type houseHoldHandler struct {
	service     domainservice.HouseHoldService
	userService domainservice.UserAccountService
}

type CreateShoppingRecordRequest struct {
	HouseholdID uint   `json:"householdID"`
	CategoryID  uint   `json:"categoryID"`
	Amount      int    `json:"amount"`
	Date        string `json:"date"`
	Memo        string `json:"memo"`
}

type UpdateShoppingRecordRequest struct {
	CategoryID uint   `json:"categoryID"`
	Amount     int    `json:"amount"`
	Date       string `json:"date"`
	Memo       string `json:"memo"`
}

type AddHouseHoldCategoryRequest struct {
	HouseholdID         uint   `json:"householdID" param:"householdID"`
	CategoryName        string `json:"categoryName"`
	CategoryLimitAmount int    `json:"categoryLimitAmount"`
}

type AddHouseHoldRequest struct {
	UserID      uint   `json:"id" param:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// AddHouseHold implements HouseHoldHandler.
func (h *houseHoldHandler) AddHouseHold(c echo.Context) error {
	req := AddHouseHoldRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	houseHold := domainmodel.HouseHold{
		UserID:      domainmodel.UserID(req.UserID),
		Title:       req.Title,
		Description: req.Description,
	}

	if err := h.service.AddUserHouseHold(&houseHold); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "success")
}

// AddHouseHoldCategory implements HouseHoldHandler.
func (h *houseHoldHandler) AddHouseHoldCategory(c echo.Context) error {
	req := AddHouseHoldCategoryRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.service.AddHouseHoldCategory(domainmodel.HouseHoldID(req.HouseholdID), req.CategoryName, req.CategoryLimitAmount); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "success")
}

// CreateShoppingRecord implements HouseHoldHandler.
func (h *houseHoldHandler) CreateShoppingRecord(c echo.Context) error {
	req := CreateShoppingRecordRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	shoppingAmount := domainmodel.NewShoppingAmount(domainmodel.HouseHoldID(req.HouseholdID), domainmodel.CategoryID(req.CategoryID), req.Amount, req.Date, req.Memo, 0)

	if err := h.service.CreateShoppingAmount(shoppingAmount); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "success")
}

// UpdateShoppingRecord implements HouseHoldHandler.
func (h *houseHoldHandler) UpdateShoppingRecord(c echo.Context) error {
	req := UpdateShoppingRecordRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	householdID := c.Param("householdID")
	shoppingID := c.Param("shoppingID")

	_, err := strconv.ParseUint(householdID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	shoppingIDUint, err := strconv.ParseUint(shoppingID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	shoppingAmount := &domainmodel.ShoppingAmount{
		ID:         domainmodel.ShoppingID(uint(shoppingIDUint)),
		CategoryID: domainmodel.CategoryID(req.CategoryID),
		Amount:     req.Amount,
		Date:       req.Date,
		Memo:       req.Memo,
	}

	if err := h.service.UpdateShoppingAmount(shoppingAmount); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "success")

}

// FetchShoppingRecord implements HouseHoldHandler.
func (h *houseHoldHandler) FetchShoppingRecord(c echo.Context) error {
	householdID := c.Param("householdID")
	householdIDUint, err := strconv.ParseUint(householdID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("2006-01")
	}

	input := domainservice.FetchShoppingRecordInput{
		HouseholdID: domainmodel.HouseHoldID(uint(householdIDUint)),
		Date:        date,
	}

	results, err := h.service.SummarizeShoppingAmount(input)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// RemoveShoppingRecord implements HouseHoldHandler.
func (h *houseHoldHandler) RemoveShoppingRecord(c echo.Context) error {
	householdID := c.Param("householdID")
	shoppingID := c.Param("shoppingID")

	// TODO : houseHoldIDは不要か？
	_, err := strconv.ParseUint(householdID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	shoppingIDUint, err := strconv.ParseUint(shoppingID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.service.RemoveShoppingAmount(domainmodel.ShoppingID(uint(shoppingIDUint))); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "success")
}

// ShareHouseHold implements HouseHoldHandler.
func (h *houseHoldHandler) ShareHouseHold(c echo.Context) error {
	householdID := c.Param("householdID")
	inviteUserID := c.Param("inviteUserID")

	householdIDUint, err := strconv.ParseUint(householdID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	inviteUserIDUint, err := strconv.ParseUint(inviteUserID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.service.ShareHouseHold(domainmodel.HouseHoldID(uint(householdIDUint)), domainmodel.UserID(uint(inviteUserIDUint)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "success")
}

// FetchHouseHoldUser implements HouseHoldHandler.
func (h *houseHoldHandler) FetchHouseHoldUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.userService.FetchUserAccount(domainmodel.UserID(uint(id)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// FetchHouseHold implements HouseHoldHandler.
func (h *houseHoldHandler) FetchHouseHold(c echo.Context) error {
	houseHoldID := c.Param("id")
	houseHoldIDUint, err := strconv.ParseUint(houseHoldID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	houseHold, err := h.service.FetchHouseHold(domainmodel.HouseHoldID(uint(houseHoldIDUint)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, houseHold)
}

type HouseHoldHandler interface {
	FetchHouseHold(c echo.Context) error
	FetchHouseHoldUser(c echo.Context) error
	ShareHouseHold(c echo.Context) error
	AddHouseHoldCategory(c echo.Context) error
	AddHouseHold(c echo.Context) error
	// 買い物記録
	FetchShoppingRecord(c echo.Context) error
	CreateShoppingRecord(c echo.Context) error
	UpdateShoppingRecord(c echo.Context) error
	RemoveShoppingRecord(c echo.Context) error
}

func NewHouseHoldHandler(service domainservice.HouseHoldService, userService domainservice.UserAccountService) HouseHoldHandler {
	return &houseHoldHandler{service: service, userService: userService}
}
