package helper

import (
	"Zhooze/domain"
	"Zhooze/utils/models"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func GetTimeFromPeriod(timePeriod string) (time.Time, time.Time) {
	endDate := time.Now()

	if timePeriod == "day" {
		startDate := endDate.AddDate(0, 0, -1)
		return startDate, endDate
	}

	if timePeriod == "week" {
		startDate := endDate.AddDate(0, 0, -6)
		return startDate, endDate
	}

	if timePeriod == "year" {
		startDate := endDate.AddDate(-1, 0, 0)
		return startDate, endDate
	}

	return endDate.AddDate(0, 0, -6), endDate
}

func CopyOrderDetails(orderDetails domain.Order, orderBody models.OrderIncoming) domain.Order {

	id := uuid.New().ID()
	str := strconv.Atoi(string(id))
	orderDetails.ID = str[:8]
	orderDetails.AddressID = orderBody.AddressID
	orderDetails.PaymentMethodID = orderBody.PaymentID
	orderDetails.UserID = orderBody.UserID
	orderDetails.Approval = false
	orderDetails.ShipmentStatus = "processing"
	orderDetails.PaymentStatus = "not paid"

	return orderDetails

}
