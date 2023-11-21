package domain

import "gorm.io/gorm"

type Coupons struct {
	ID                 uint    `json:"id" gorm:"uniquekey; not null"`
	Coupon             string  `json:"coupon" gorm:"coupon"`
	DiscountPercentage int     `json:"discount_percentage"`
	Validity           bool    `json:"validity"`
	MinimumPrice       float64 `json:"minimum_price"`
}

type UsedCoupon struct {
	ID       uint    `json:"id" gorm:"uniquekey not null"`
	CouponID uint    `json:"coupon_id"`
	Coupons  Coupons `json:"-" gorm:"foreignkey:CouponID"`
	UserID   uint    `json:"user_id"`
	Users    User   `json:"-" gorm:"foreignkey:UserID"`
	Used     bool    `json:"used"`
}

type Referral struct {
	gorm.Model
	UserID         uint    `json:"user_id" gorm:"uniquekey; not null"`
	Users          User   `json:"-" gorm:"foreignkey:UserID"`
	ReferralCode   string  `json:"referral_code"`
	ReferralAmount float64 `json:"referral_amount"`
	ReferredUserID uint    `json:"referred_user_id"`
}
