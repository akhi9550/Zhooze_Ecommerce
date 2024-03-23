package usecase

import (
	"Zhooze/repository"
	"Zhooze/utils/models"
)

func GetWallet(userID int) (models.WalletAmount, error) {
	return repository.GetWallet(userID)
}

func GetWalletHistory(userID int)([]models.WalletHistory,error){
	return repository.GetWalletHistory(userID)

}