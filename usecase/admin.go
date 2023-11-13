package usecase

import (
	"Zhooze/domain"
	"Zhooze/helper"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(adminDetails models.AdminLogin) (domain.TokenAdmin, error) {
	adminCompareDetails, err := repository.LoginHandler(adminDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	// compare password from database and that provided from admins
	err = bcrypt.CompareHashAndPassword([]byte(adminCompareDetails.Password), []byte(adminDetails.Password))
	if err != nil {
		return domain.TokenAdmin{}, err
	}
	var adminDetailsResponse models.AdminDetailsResponse
	//  copy all details except password and sent it back to the front end
	err = copier.Copy(&adminDetailsResponse, &adminCompareDetails)
	if err != nil {
		return domain.TokenAdmin{}, err
	}

	tokenString, err := helper.GenerateTokenAdmin(adminDetailsResponse)

	if err != nil {
		return domain.TokenAdmin{}, err
	}

	return domain.TokenAdmin{
		Admin: adminDetailsResponse,
		Token: tokenString,
	}, nil
}

func DashBoard() (models.CompleteAdminDashboard, error) {
	userDetails, err := repository.DashBoardUserDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	productDetails, err := repository.DashBoardProductDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	orderDetails, err := repository.DashBoardOrder()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	totalRevenue, err := repository.TotalRevenue()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	amountDetails, err := repository.AmountDetails()
	if err != nil {
		return models.CompleteAdminDashboard{}, err
	}
	return models.CompleteAdminDashboard{
		DashboardUser:    userDetails,
		DashboardProduct: productDetails,
		DashboardOrder:   orderDetails,
		DashboardRevenue: totalRevenue,
		DashboardAmount:  amountDetails,
	}, nil
}
func ShowAllUsers(page, count int) ([]models.UserDetailsAtAdmin, error) {
	users, err := repository.ShowAllUsersIn(page, count)
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}
	return users, nil
}
func BlockedUser(id string) error {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return err
	}
	if user.Blocked {
		return errors.New("already blocked")
	} else {
		user.Blocked = true
	}
	err = repository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}
	return nil
}
func UnBlockedUser(id string) error {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return err
	}
	if user.Blocked {
		user.Blocked = false
	} else {
		return errors.New("already unblocked")
	}
	err = repository.UpdateBlockUserByID(user)
	if err != nil {
		return err
	}
	return nil
}
func GetAllOrderDetailsForAdmin(page, pagesize int) ([]models.CombinedOrderDetails, error) {
	orderDetail, err := repository.GetAllOrderDetailsBrief(page, pagesize)
	if err != nil {
		return []models.CombinedOrderDetails{}, err
	}
	return orderDetail, nil
}
func ApproveOrder(orderId int) error {
	ShipmentStatus, err := repository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}
	if ShipmentStatus == "cancelled" {
		return errors.New("the order is cancelled,cannot approve it")
	}
	if ShipmentStatus == "pending" {
		return errors.New("the order is pending,cannot approve it")
	}
	if ShipmentStatus == "processing" {
		err := repository.ApproveOrder(orderId)
		if err != nil {
			return err
		}
		return nil
	}
	// if the shipment status is not processing or cancelled. Then it is defenetely cancelled
	return nil
}
func CancelOrderFromAdmin(order_id int) error {
	ok, err := repository.CheckOrderID(order_id)
	fmt.Println(err)
	if !ok {
		return err
	}
	orderProduct, err := repository.GetProductDetailsFromOrders(order_id)
	if err != nil {
		return err
	}
	err = repository.CancelOrders(order_id)
	if err != nil {
		return err
	}
	fmt.Println("😶‍🌫️😶‍🌫️😶‍🌫️😶‍🌫️", orderProduct)
	err = repository.UpdateStockOfProduct(orderProduct)
	if err != nil {
		return err
	}
	return nil
}

//	func RefundUser(orderID string) error {
//		paymentStatus, err := repository.GetPaymentStatus(orderID)
//		if err != nil {
//			return err
//		}
//		if paymentStatus == "refund-init" {
//			paymentStatus = "refunded"
//			return repository.RefundOrder(paymentStatus, orderID)
//		}
//		return errors.New("cannot refund the order")
//	}
func FilteredSalesReport(timePeriod string) (models.SalesReport, error) {
	startTime, endTime := helper.GetTimeFromPeriod(timePeriod)
	//update date (end and start)
	saleReport, err := repository.FilteredSalesReport(startTime, endTime)

	if err != nil {
		return models.SalesReport{}, err
	}
	return saleReport, nil

}
func AddPaymentMethod(payment models.NewPaymentMethod) (domain.PaymentMethod, error) {
	exists, err := repository.CheckIfPaymentMethodAlreadyExists(payment.PaymentName)
	if err != nil {
		return domain.PaymentMethod{}, err
	}
	if exists {
		return domain.PaymentMethod{}, errors.New("payment method already exists")
	}
	paymentadd, err := repository.AddPaymentMethod(payment)
	if err != nil {
		return domain.PaymentMethod{}, err
	}
	return paymentadd, nil
}

func ListPaymentMethods() ([]domain.PaymentMethod, error) {

	categories, err := repository.ListPaymentMethods()
	if err != nil {
		return []domain.PaymentMethod{}, err
	}
	return categories, nil

}

func DeletePaymentMethod(id int) error {

	err := repository.DeletePaymentMethod(id)
	if err != nil {
		return err
	}
	return nil

}
