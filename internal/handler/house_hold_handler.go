package handler

import (
	domainmodel "echo-household-budget/internal/domain/model"
	domainservice "echo-household-budget/internal/domain/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type houseHoldHandler struct {
	service     domainservice.HouseHoldService
	userService domainservice.UserAccountService
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
}

func NewHouseHoldHandler(service domainservice.HouseHoldService, userService domainservice.UserAccountService) HouseHoldHandler {
	return &houseHoldHandler{service: service, userService: userService}
}
