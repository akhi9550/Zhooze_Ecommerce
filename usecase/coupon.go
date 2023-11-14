package usecase

import (
	"Zhooze/repository"
	"Zhooze/utils/models"
	"errors"
)

func AddCoupon(coupon models.AddCoupon) (string, error) {

	// if coupon already exist and if it is expired revalidate it. else give back an error message saying the coupon already exist
	couponExist, err := repository.CouponExist(coupon.Coupon)
	if err != nil {
		return "", err
	}

	if couponExist {
		alreadyValid, err := repository.CouponRevalidateIfExpired(coupon.Coupon)
		if err != nil {
			return "", err
		}

		if alreadyValid {
			return "The coupon which is valid already exists", nil
		}

		return "Made the coupon valid", nil

	}

	err = repository.AddCoupon(coupon)
	if err != nil {
		return "", err
	}
	return "successfully added the coupon", nil
}
func GetCoupon() ([]models.Coupon, error) {
	coupons, err := repository.GetCoupon()
	if err != nil {
		return []models.Coupon{}, err
	}
	return coupons, nil
}
func ExpireCoupon(couponID int) error {
	couponExist, err := repository.ExistCoupon(couponID)
	if err != nil {
		return err
	}
	// if it exists expire it, if already expired send back relevant message
	if couponExist {
		err = repository.CouponAlreadyExpired(couponID)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("coupon does not exist")

}

func ApplyCoupon(coupon string, userID int) error {

	cartExist, err := repository.DoesCartExist(userID)
	if err != nil {
		return err
	}

	if !cartExist {
		return errors.New("cart empty, can't apply coupon")
	}

	couponExist, err :=repository.CouponExist(coupon)
	if err != nil {
		return err
	}

	if !couponExist {
		return errors.New("coupon does not exist")
	}

	couponValidity, err :=repository.CouponValidity(coupon)
	if err != nil {
		return err
	}

	if !couponValidity {
		return errors.New("coupon expired")
	}

	minDiscountPrice, err := repository.GetCouponMinimumAmount(coupon)
	if err != nil {
		return err
	}

	totalPriceFromCarts, err :=repository.GetTotalPriceFromCart(userID)
	if err != nil {
		return err
	}

	// if the total Price is less than minDiscount price don't allow coupon to be added
	if totalPriceFromCarts < minDiscountPrice {
		return errors.New("coupon cannot be added as the total amount is less than minimum amount for coupon")
	}

	userAlreadyUsed, err := repository.DidUserAlreadyUsedThisCoupon(coupon, userID)
	if err != nil {
		return err
	}

	if userAlreadyUsed {
		return errors.New("user already used this coupon")
	}

	couponStatus, err := repository.UpdateUsedCoupon(coupon, userID)
	if err != nil {
		return err
	}

	if couponStatus {
		return nil
	}
	return errors.New("could not add the coupon")

}
