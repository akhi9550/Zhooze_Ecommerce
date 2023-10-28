package usecase

import (
	"Zhooze/config"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"

	"github.com/razorpay/razorpay-go"
)

func MakePaymentRazorPay(orderID string, userID int) (models.CombinedOrderDetails, string, error) {
	cfg, _ := config.LoadConfig()
	combainedOrderDetails, err := repository.GetOrderDetailsByOrderId(orderID)
	if err != nil {
		return models.CombinedOrderDetails{}, "", err
	}
	client := razorpay.NewClient(cfg.KEY_ID_fOR_PAY, cfg.SECRET_KEY_FOR_PAY)

	data := map[string]interface{}{
		"amount":   int(combainedOrderDetails.FinalPrice) * 100,
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
	return models.CombinedOrderDetails{}, razorPayOrderID, nil
}

func SavePaymentDetails(paymentID string, orderID string) error {
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
	return errors.New("already paid")
}
