package repository

import (
	"Zhooze/db"
	"Zhooze/domain"
	"Zhooze/helper"
	"Zhooze/utils/models"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func LoginHandler(adminDetails models.AdminLogin) (domain.Admin, error) {
	var details domain.Admin
	if err := db.DB.Raw("SELECT * FROM users WHERE email=? AND isadmin= true", adminDetails.Email).Scan(&details).Error; err != nil {
		return domain.Admin{}, err
	}
	return details, nil
}
func DashBoardUserDetails() (models.DashBoardUser, error) {
	var userDetails models.DashBoardUser
	err := db.DB.Raw("SELECT COUNT(*) FROM users WHERE isadmin='false'").Scan(&userDetails.TotalUsers).Error
	if err != nil {
		return models.DashBoardUser{}, nil
	}
	err = db.DB.Raw("SELECT COUNT(*)  FROM users WHERE blocked=true").Scan(&userDetails.BlockedUser).Error
	if err != nil {
		return models.DashBoardUser{}, nil
	}
	return userDetails, nil
}

func DashBoardProductDetails() (models.DashBoardProduct, error) {
	var productDetails models.DashBoardProduct
	err := db.DB.Raw("SELECT COUNT(*) FROM products").Scan(&productDetails.TotalProducts).Error
	if err != nil {
		return models.DashBoardProduct{}, nil
	}
	err = db.DB.Raw("SELECT COUNT(*) FROM products WHERE stock<=0").Scan(&productDetails.OutofStockProduct).Error
	if err != nil {
		return models.DashBoardProduct{}, nil
	}
	return productDetails, nil
}
func ShowAllUsersIn(page, count int) ([]models.UserDetailsAtAdmin, error) {
	var user []models.UserDetailsAtAdmin
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * count
	err := db.DB.Raw("SELECT id,firstname,lastname,email,phone FROM users WHERE isadmin='false' limit ? offset ?", count, offset).Scan(&user).Error
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}
	return user, nil
}
func GetUserByID(id string) (domain.User, error) {
	user_id, err := strconv.Atoi(id)
	if err != nil {
		return domain.User{}, err
	}
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM users WHERE id=?", user_id).Scan(&count).Error; err != nil {

		return domain.User{}, err
	}
	if count < 1 {
		return domain.User{}, errors.New("user for the given id does not exist")

	}
	var userDetails domain.User
	if err := db.DB.Raw("SELECT * FROM users WHERE id=?", user_id).Scan(&userDetails).Error; err != nil {
		return domain.User{}, err
	}
	return userDetails, nil
}

func UpdateBlockUserByID(user domain.User) error {
	err := db.DB.Exec("UPDATE users SET blocked = ? WHERE id = ?", user.Blocked, user.ID).Error
	if err != nil {
		fmt.Println("Error updating user:", err)
		return err
	}
	return nil
}
func DashBoardOrder() (models.DashboardOrder, error) {
	var orderDetail models.DashboardOrder
	err := db.DB.Raw("SELECT COUNT(*) FROM orders WHERE payment_status= 'paid' AND approval =true").Scan(&orderDetail.CompletedOrder).Error
	if err != nil {
		return models.DashboardOrder{}, err
	}
	err = db.DB.Raw("SELECT COUNT(*) FROM orders WHERE shipment_status='pending' OR shipment_status = 'processing'").Scan(&orderDetail.PendingOrder).Error
	if err != nil {
		return models.DashboardOrder{}, err
	}
	err = db.DB.Raw("select count(*) from orders where shipment_status = 'cancelled'").Scan(&orderDetail.CancelledOrder).Error
	if err != nil {
		return models.DashboardOrder{}, nil
	}

	err = db.DB.Raw("select count(*) from orders").Scan(&orderDetail.TotalOrder).Error
	if err != nil {
		return models.DashboardOrder{}, nil
	}

	err = db.DB.Raw("select COALESCE(SUM(quantity),0) from carts").Scan(&orderDetail.TotalOrderItem).Error
	if err != nil {
		return models.DashboardOrder{}, nil
	}
	return orderDetail, nil
}
func TotalRevenue() (models.DashboardRevenue, error) {
	var revenueDetails models.DashboardRevenue
	startTime := time.Now().AddDate(0, 0, -1)
	endTime := time.Now()
	err := db.DB.Raw("SELECT COALESCE(SUM(final_price),0) FROM orders WHERE payment_status = 'paid' AND approval = true AND created_at >=? AND created_at <=?", startTime, endTime).Scan(&revenueDetails.TodayRevenue).Error
	if err != nil {
		return models.DashboardRevenue{}, nil
	}
	startTime, endTime = helper.GetTimeFromPeriod("month")
	err = db.DB.Raw("SELECT COALESCE (SUM(final_price),0) FROM orders WHERE payment_status = 'paid' AND approval = true AND created_at >=? AND created_at <=?", startTime, endTime).Scan(&revenueDetails.MonthRevenue).Error
	if err != nil {
		return models.DashboardRevenue{}, nil
	}
	startTime, endTime = helper.GetTimeFromPeriod("year")
	err = db.DB.Raw("SELECT COALESCE (SUM(final_price),0) FROM orders WHERE payment_status = 'paid' AND approval = true AND created_at >=? AND created_at <=?", startTime, endTime).Scan(&revenueDetails.YearRevenue).Error
	if err != nil {
		return models.DashboardRevenue{}, nil
	}
	return revenueDetails, nil
}

func AmountDetails() (models.DashboardAmount, error) {
	var amountDetails models.DashboardAmount
	err := db.DB.Raw("SELECT COALESCE (SUM(final_price),0) FROM orders WHERE payment_status = 'paid' AND approval = true").Scan(&amountDetails.CreditedAmount).Error
	if err != nil {
		return models.DashboardAmount{}, nil
	}
	err = db.DB.Raw("SELECT COALESCE(SUM(final_price),0) FROM orders WHERE payment_status = 'not paid' AND shipment_status = 'processing' OR shipment_status = 'pending' OR shipment_status = 'order placed'").Scan(&amountDetails.PendingAmount).Error
	if err != nil {
		return models.DashboardAmount{}, nil
	}
	return amountDetails, nil
}
func FilteredSalesReport(startTime time.Time, endTime time.Time) (models.SalesReport, error) {
	var salesReport models.SalesReport
	result := db.DB.Raw("SELECT COALESCE(SUM(final_price),0) FROM orders WHERE payment_status='paid' AND approval = true AND created_at >= ? AND created_at <= ?", startTime, endTime).Scan(&salesReport.TotalSales)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	result = db.DB.Raw("SELECT COUNT(*) FROM orders").Scan(&salesReport.TotalOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	result = db.DB.Raw("SELECT COUNT(*) FROM orders WHERE payment_status = 'paid' and approval = true and created_at >= ? AND created_at <= ?", startTime, endTime).Scan(&salesReport.CompletedOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	result = db.DB.Raw("SELECT COUNT(*) FROM orders WHERE shipment_status = 'processing' AND approval = false AND created_at >= ? AND created_at<=?", startTime, endTime).Scan(&salesReport.PendingOrders)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	var productID int
	result = db.DB.Raw("SELECT product_id FROM order_items GROUP BY product_id order by SUM(quantity) DESC LIMIT 1").Scan(&productID)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	result = db.DB.Raw("SELECT name FROM products WHERE id = ?", productID).Scan(&salesReport.TrendingProduct)
	if result.Error != nil {
		return models.SalesReport{}, result.Error
	}
	return salesReport, nil
}
func AddPaymentMethod(pay models.NewPaymentMethod) (domain.PaymentMethod, error) {
	var payment string
	if err := db.DB.Raw("INSERT INTO payment_methods (payment_name) VALUES (?) RETURNING payment_name", pay.PaymentName).Scan(&payment).Error; err != nil {
		return domain.PaymentMethod{}, err
	}
	var paymentResponse domain.PaymentMethod
	err := db.DB.Raw("SELECT id, payment_name FROM payment_methods WHERE payment_name = ?", payment).Scan(&paymentResponse).Error
	if err != nil {
		return domain.PaymentMethod{}, err
	}
	return paymentResponse, nil

}

func ListPaymentMethods() ([]domain.PaymentMethod, error) {
	var model []domain.PaymentMethod
	err := db.DB.Raw("SELECT * FROM payment_methods").Scan(&model).Error
	if err != nil {
		return []domain.PaymentMethod{}, err
	}

	return model, nil
}

func CheckIfPaymentMethodAlreadyExists(payment string) (bool, error) {
	var count int64
	err := db.DB.Raw("SELECT COUNT(*) FROM payment_methods WHERE payment_name = $1", payment).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DeletePaymentMethod(id int) error {
	var count int
	if err := db.DB.Raw("SELECT COUNT(*) FROM payment_methods WHERE id=?", id).Scan(&count).Error; err != nil {
		return err
	}
	if count < 1 {
		return errors.New("payment for given id does not exist")
	}

	if err := db.DB.Exec("DELETE FROM payment_methods WHERE id=?", id).Error; err != nil {
		return err
	}
	return nil
}
