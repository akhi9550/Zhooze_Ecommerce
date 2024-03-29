package helper

import (
	"gorm.io/gorm"
)

func GetCouponDiscountPrice(userID int, TotalPrice float64, DB *gorm.DB) (float64, error) {
	var count int
	err := DB.Raw("SELECT COUNT(*) FROM used_coupons WHERE user_id = ? AND used = false", userID).Scan(&count).Error
	if err != nil {
		return 0.0, err
	}

	if count < 0 {
		return 0.0, nil
	}

	type CouponDetails struct {
		DiscountPercentage int
		MinimumPrice       float64
	}

	// take the discount percentage and minimum price to check the condition ( !! Actually this is not needed. As all the conditions were checked while adding the coupon !!)
	// just discount percentage would work fine - should refactor this in the future
	var coup CouponDetails
	err = DB.Raw("select discount_percentage,minimum_price from coupons where id = (select coupon_id from used_coupons where user_id = ? and used = false)", userID).Scan(&coup).Error
	if err != nil {
		return 0.0, err
	}

	var totalPrice float64
	err = DB.Raw("select COALESCE(SUM(total_price), 0) from carts where user_id = ?", userID).Scan(&totalPrice).Error
	if err != nil {
		return 0.0, err
	}

	if totalPrice < coup.MinimumPrice {
		return 0.0, nil
	}

	return ((float64(coup.DiscountPercentage) * totalPrice) / 100), nil

}
func GetReferralDiscountPrice(FinalPrice float64, userID int, DB *gorm.DB) (float64, error) {
	var count int
	err := DB.Raw("SELECT COUNT(*) FROM referrals WHERE user_id = ?", userID).Scan(&count).Error
	if err != nil {
		return 0.0, err
	}
	if count < 0 {
		return 0.0, nil
	}
	var Price float64
	err = DB.Raw("SELECT referral_amount FROM referrals WHERE user_id = ? ", userID).Scan(&Price).Error
	if err != nil {
		return 0.0, err
	}
	if FinalPrice < Price {
		return 0.0, nil
	}
	return (FinalPrice - Price), nil

}
