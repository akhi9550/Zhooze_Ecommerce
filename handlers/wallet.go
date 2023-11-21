package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Wallet Details
// @Description Wallet TotalPrice from User Profile
// @Tags Wallet
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/wallet   [GET]
func GetWallet(c *gin.Context) {
	userID, _ := c.Get("user_id")
	WalletDetails, err := usecase.GetWallet(userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Wallet Details", WalletDetails, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary Wallet History Details
// @Description Wallet Details from User Profile
// @Tags Wallet
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /user/wallet/history   [GET]
func WalletHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	WalletDetails, err := usecase.GetWalletHistory(userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Wallet Histories", WalletDetails, nil)
	c.JSON(http.StatusOK, success)
}
