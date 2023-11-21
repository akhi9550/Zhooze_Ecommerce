package repository

import (
	"Zhooze/db"
	"Zhooze/utils/models"
)

func GetWallet(userID int) (models.WalletAmount, error) {
	var walletAmount models.WalletAmount
	err := db.DB.Raw("SELECT amount FROM wallets WHERE user_id = ?", userID).Scan(&walletAmount).Error
	if err != nil {
		return models.WalletAmount{}, err
	}
	return walletAmount, nil
}
func GetWalletHistory(userID int) ([]models.WalletHistory, error) {
	var history []models.WalletHistory
	err := db.DB.Raw("SELECT id,order_id,reason,amount FROM wallet_histories WHERE user_id = ?", userID).Scan(&history).Error
	if err != nil {
		return []models.WalletHistory{}, err
	}
	return history, nil
}
