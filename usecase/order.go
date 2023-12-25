package usecase

import (
	"Zhooze/domain"
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
	"fmt"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/jung-kurt/gofpdf"
)

func OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error) {
	var orderBody models.OrderIncoming
	err := copier.Copy(&orderBody, &orderFromCart)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	orderBody.UserID = userID
	cartExist, err := repository.DoesCartExist(userID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if !cartExist {
		return domain.OrderSuccessResponse{}, errors.New("cart empty can't order")
	}

	addressExist, err := repository.AddressExist(orderBody)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !addressExist {
		return domain.OrderSuccessResponse{}, errors.New("address does not exist")
	}
	PaymentExist, err := repository.PaymentExist(orderBody)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !PaymentExist {
		return domain.OrderSuccessResponse{}, errors.New("paymentmethod does not exist")
	}
	cartItems, err := repository.GetAllItemsFromCart(orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	total, err := repository.TotalAmountInCart(orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	discount_price, err := repository.GetCouponDiscountPrice(int(orderBody.UserID), total)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	err = repository.UpdateCouponDetails(discount_price, orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	FinalPrice := total - discount_price
	if orderBody.PaymentID == 3 {
		wallectAmount, err := repository.WallectAmount(userID)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}
		if FinalPrice >= wallectAmount {
			return domain.OrderSuccessResponse{}, errors.New("this much of amount not available in wallet")
		}
	}

	order_id, err := repository.OrderItems(orderBody, FinalPrice)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	if orderBody.PaymentID == 3 {
		if err := repository.UpdateWallectAfterOrder(userID, FinalPrice); err != nil {
			return domain.OrderSuccessResponse{}, err
		}
		reason := "Amount debited for purchasing products"
		err = repository.UpdateHistoryForDebit(userID, order_id, FinalPrice, reason)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}
	}
	if err := repository.AddOrderProducts(order_id, cartItems); err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	orderSuccessResponse, err := repository.GetBriefOrderDetails(order_id, orderBody.PaymentID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	var orderItemDetails domain.OrderItem
	for _, c := range cartItems {
		orderItemDetails.ProductID = c.ProductID
		orderItemDetails.Quantity = c.Quantity
		err := repository.UpdateCartAfterOrder(userID, int(orderItemDetails.ProductID), orderItemDetails.Quantity)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}
	}
	return orderSuccessResponse, nil
}
func GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {

	fullOrderDetails, err := repository.GetOrderDetails(userId, page, count)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return fullOrderDetails, nil

}

func CancelOrders(orderID int, userId int) error {
	userTest, err := repository.UserOrderRelationship(orderID, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New("the order is not done by this user")
	}
	orderProductDetails, err := repository.GetProductDetailsFromOrders(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := repository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("item already delivered, cannot cancel")
	}

	if shipmentStatus == "pending" || shipmentStatus == "returned" || shipmentStatus == "return" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("the order is in " + message + ", so no point in cancelling")
	}

	if shipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled, so no point in cancelling")
	}
	err = repository.CancelOrders(orderID)
	if err != nil {
		return err
	}
	payment_status, err := repository.PaymentStatus(orderID)
	if err != nil {
		return err
	}
	err = repository.UpdateQuantityOfProduct(orderProductDetails)
	if err != nil {
		return err
	}
	amount, err := repository.TotalAmountFromOrder(orderID)
	if err != nil {
		return err
	}
	if payment_status == "refunded" {
		err = repository.UpdateAmountToWallet(userId, amount)
		if err != nil {
			return err
		}
		reason := "Amount credited for cancellation of order by user"
		err := repository.UpdateHistory(userId, orderID, amount, reason)
		if err != nil {
			return err
		}
	}

	return nil

}

func Checkout(userID int) (models.CheckoutDetails, error) {
	allUserAddress, err := repository.GetAllAddresses(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	paymentDetails, err := repository.GetAllPaymentOption(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	cartItems, err := repository.DisplayCart(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}
	grandTotal, err := repository.GetTotalPrice(userID)
	if err != nil {
		return models.CheckoutDetails{}, err
	}

	return models.CheckoutDetails{
		AddressInfoResponse: allUserAddress,
		Payment_Method:      paymentDetails,
		Cart:                cartItems,
		Total_Price:         grandTotal.FinalPrice,
	}, nil
}
func PaymentMethodID(order_id int) (int, error) {
	id, err := repository.PaymentMethodID(order_id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func ExecutePurchaseCOD(orderID int) error {
	err := repository.OrderExist(orderID)
	if err != nil {
		return err
	}
	shipmentStatus, err := repository.GetShipmentStatus(orderID)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("item  delivered, cannot pay")
	}
	if shipmentStatus == "order placed" {
		return errors.New("item placed, cannot pay")
	}
	if shipmentStatus == "cancelled" || shipmentStatus == "returned" || shipmentStatus == "return" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("the order is in" + message + "so can't paid")
	}
	if shipmentStatus == "processing" {
		return errors.New("the order is already paid")
	}
	err = repository.UpdateOrder(orderID)
	if err != nil {
		return err
	}

	return nil

}
func PrintInvoice(orderId int) (*gofpdf.Fpdf, error) {

	if orderId < 1 {
		return nil, errors.New("enter a valid order id")
	}

	order, err := repository.GetDetailedOrderThroughId(orderId)
	if err != nil {
		return nil, err
	}

	items, err := repository.GetItemsByOrderId(orderId)
	if err != nil {
		return nil, err
	}

	fmt.Println("order details ", order)
	fmt.Println("itemssss", items)
	fmt.Println("order status", order.ShipmentStatus)
	if order.ShipmentStatus != "order placed" {
		return nil, errors.New("wait for the invoice until the product is received")
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(31, 73, 125)
	pdf.Cell(0, 20, "Invoice")
	pdf.Ln(20)

	pdf.SetFont("Arial", "I", 14)
	pdf.SetTextColor(51, 51, 51)
	pdf.Cell(0, 10, "Customer Details")
	pdf.Ln(10)
	customerDetails := []string{
		"Name: " + order.Firstname,
		"House Name: " + order.HouseName,
		"Street: " + order.Street,
		"State: " + order.State,
		"City: " + order.City,
	}
	for _, detail := range customerDetails {
		pdf.Cell(0, 10, detail)
		pdf.Ln(10)
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, "Item", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Total Price", "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.SetFillColor(255, 255, 255)
	for _, item := range items {
		pdf.CellFormat(40, 10, item.ProductName, "1", 0, "L", true, 0, "")
		pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(item.TotalPrice, 'f', 2, 64), "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 10, strconv.Itoa(int(item.Quantity)), "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(item.TotalPrice, 'f', 2, 64), "1", 0, "C", true, 0, "")
		pdf.Ln(10)
	}
	pdf.Ln(10)

	var totalPrice float64
	for _, item := range items {
		totalPrice += item.TotalPrice
	}

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Total Price:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(totalPrice, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	OfferApplied := totalPrice - order.FinalPrice

	fmt.Println("offer Applied", OfferApplied)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Offer Applied:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(OfferApplied, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetFillColor(217, 217, 217)
	pdf.CellFormat(120, 10, "Final Amount:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, "$"+strconv.FormatFloat(order.FinalPrice, 'f', 2, 64), "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	return pdf, nil
}
