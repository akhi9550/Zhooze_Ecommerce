package repository

import (
	"Zhooze/db"
	"Zhooze/utils/models"
	"errors"
)

func CouponExist(couponName string) (bool, error) {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM coupons WHERE coupon = ?", couponName).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func CouponValidity(couponName string) (bool, error) {
	var validity bool
	err := db.DB.Raw("SELECT validity FROM coupons WHERE coupon = ?", couponName).Scan(&validity).Error
	if err != nil {
		return false, err
	}

	return validity, nil

}

func CouponRevalidateIfExpired(couponName string) (bool, error) {
	var isValid bool
	err := db.DB.Raw("SELECT validity FROM coupons WHERE coupon = ?", couponName).Scan(&isValid).Error
	if err != nil {
		return false, err
	}

	if isValid {
		return true, nil
	}

	err = db.DB.Exec("UPDATE coupons SET validity = true WHERE coupon = ?", couponName).Error
	if err != nil {
		return false, err
	}
	return false, nil
}

func AddCoupon(coupon models.AddCoupon) error {
	err := db.DB.Exec("INSERT INTO coupons (coupon,discount_percentage,minimum_price,validity) VALUES (?, ?, ?, ?)", coupon.Coupon, coupon.DiscountPercentage, coupon.MinimumPrice, true).Error
	if err != nil {
		return nil
	}
	return nil

}

func GetCoupon() ([]models.Coupon, error) {
	var coupons []models.Coupon
	err := db.DB.Raw("SELECT id,coupon,discount_percentage,minimum_price,Validity FROM coupons").Scan(&coupons).Error
	if err != nil {
		return []models.Coupon{}, err
	}
	return coupons, nil
}

func ExistCoupon(couponID int) (bool, error) {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM coupons WHERE id = ?", couponID).Scan(&count).Error
	if err != nil {
		return false, errors.New("the offer already exists")
	}

	return count > 0, nil
}

func CouponAlreadyExpired(couponID int) error {
	var valid bool
	err := db.DB.Raw("SELECT validity FROM coupons WHERE id = ?", couponID).Scan(&valid).Error
	if err != nil {
		return err
	}
	if valid {
		err := db.DB.Exec("UPDATE coupons SET validity = false WHERE id = ?", couponID).Error
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("already expired")
}

func GetCouponMinimumAmount(coupon string) (float64, error) {

	var MinDiscountPrice float64
	err := db.DB.Raw("SELECT minimum_price FROM coupons WHERE coupon = ?", coupon).Scan(&MinDiscountPrice).Error
	if err != nil {
		return 0.0, err
	}
	return MinDiscountPrice, nil
}

func DidUserAlreadyUsedThisCoupon(coupon string, userID int) (bool, error) {
	var count int
	err := db.DB.Raw("SELECT COUNT(*) FROM used_coupons WHERE coupon_id = (SELECT id FROM coupons WHERE coupon = ?) AND user_id = ?", coupon, userID).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func UpdateUsedCoupon(coupon string, userID int) (bool, error) {
	var couponID uint
	err := db.DB.Raw("SELECT id FROM coupons WHERE coupon = ?", coupon).Scan(&couponID).Error
	if err != nil {
		return false, err
	}

	var count int
	// if a coupon have already been added, replace the order with current coupon and delete the existing coupon
	err = db.DB.Raw("SELECT count(*) FROM used_coupons WHERE user_id = ? AND used = false", userID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		err = db.DB.Exec("DELETE FROM used_coupons WHERE user_id = ? AND used = false", userID).Error
		if err != nil {
			return false, err
		}
	}
	err = db.DB.Exec("INSERT INTO used_coupons (coupon_id,user_id,used) VALUES (?, ?, false)", couponID, userID).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
