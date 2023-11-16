package usecase

import (
	"Zhooze/config"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"fmt"

	"github.com/razorpay/razorpay-go"
)

func PaymentAlreadyPaid(orderID int) (bool, error) {
	AlreadyPayed, err := repository.PaymentAlreadyPaid(orderID)
	if err != nil {
		return false, err
	}
	return AlreadyPayed, nil
}

func MakePaymentRazorPay(orderID int) (models.CombinedOrderDetails, string, error) {
	cfg, _ := config.LoadConfig()
	combinedOrderDetails, err := repository.GetOrderDetailsByOrderId(orderID)
	if err != nil {
		return models.CombinedOrderDetails{}, "", err
	}

	client := razorpay.NewClient(cfg.KEY_ID_FOR_PAY, cfg.SECRET_KEY_FOR_PAY)

	data := map[string]interface{}{
		"amount":   int(combinedOrderDetails.FinalPrice) * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {

		return models.CombinedOrderDetails{}, "", err
	}

	razorPayOrderID := body["id"].(string)

	err = repository.AddRazorPayDetails(orderID, razorPayOrderID)
	if err != nil {
		return models.CombinedOrderDetails{}, "", err
	}

	return combinedOrderDetails, razorPayOrderID, nil

}

func SavePaymentDetails(orderID int, paymentID string) error {
	status, err := repository.CheckPaymentStatus(orderID)
	if err != nil {
		return err
	}
	if status == "not paid" {
		err = repository.UpdatePaymentDetails(orderID, paymentID)
		if err != nil {
			return err
		}
		err := repository.UpdateShipmentAndPaymentByOrderID("processing", "paid", orderID)
		if err != nil {
			return err
		}
		return nil
	}
	fmt.Println("❌already paid❌")
	return errors.New("already paid")
}
