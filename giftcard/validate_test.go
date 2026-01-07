package giftcard_test

import (
	"testing"

	"github.com/fealtyx/fealtyx-go/giftcard"
	"github.com/fealtyx/fealtyx-go/giftcard/utils"
)

func TestValidateGiftCardCode(t *testing.T) {
	tests := []struct {
		name                   string // description of this test case
		domain                 string
		phoneNumber            string
		giftCardCode           string
		orderAmount            float64
		is_applicable          bool
		is_unloq_discount_code bool
		reason                 giftcard.Reason
	}{
		{
			name:                   "invalid unloq code, first 2 of last 4 characters don't match last 2",
			domain:                 "example.com",
			phoneNumber:            "+918989898989",
			giftCardCode:           "FLX1234ABCD",
			orderAmount:            100.0,
			is_applicable:          true,
			is_unloq_discount_code: false,
			reason:                 "",
		},

		{
			name:                   "invalid unloq code - missing prefix",
			domain:                 "example.com",
			phoneNumber:            "+918989898989",
			giftCardCode:           "123415EG",
			orderAmount:            100.0,
			is_applicable:          true,
			is_unloq_discount_code: false,
			reason:                 "",
		},
		{
			name:                   "invalid unloq code - short code",
			domain:                 "example.com",
			phoneNumber:            "+918989898989",
			giftCardCode:           "FL1",
			orderAmount:            100.0,
			is_applicable:          true,
			is_unloq_discount_code: false,
			reason:                 "",
		},

		{
			name:                   "valid unloq code - domain mismatch",
			domain:                 "example2.com",
			orderAmount:            100.0,
			phoneNumber:            "+918989898989",
			giftCardCode:           "FLX123415EG",
			is_applicable:          false,
			is_unloq_discount_code: true,
			reason:                 giftcard.ReasonVoucherNotEligible,
		},

		{
			name:                   "valid unloq code - phone mismatch",
			domain:                 "example.com",
			orderAmount:            100.0,
			phoneNumber:            "777777",
			giftCardCode:           "FLX123415EG",
			is_applicable:          false,
			is_unloq_discount_code: true,
			reason:                 giftcard.ReasonVoucherNotEligible,
		},

		{
			name:                   "valid unloq code with FLX prefix",
			domain:                 "example.com",
			phoneNumber:            "+918989898989",
			giftCardCode:           "FLX123415EG",
			orderAmount:            100.0,
			is_applicable:          true,
			is_unloq_discount_code: true,
			reason:                 "",
		},
		{
			name:                   "valid unloq code with UNQ prefix",
			domain:                 "example.com",
			phoneNumber:            "+918989898989",
			giftCardCode:           "UNQ123415EG",
			orderAmount:            100.0,
			is_applicable:          true,
			is_unloq_discount_code: true,
			reason:                 "",
		},
		{
			name:                   "valid unloq code with hashed phone",
			domain:                 "example.com",
			phoneNumber:            utils.GetSHA256Hash("+918989898989"),
			giftCardCode:           "FLX123415EG",
			orderAmount:            100.0,
			is_applicable:          true,
			is_unloq_discount_code: true,
			reason:                 "",
		},
		{
			name:                   "valid unloq code with non hashed phone",
			domain:                 "example.com",
			phoneNumber:            "+918989898989",
			giftCardCode:           "FLX123415EG",
			orderAmount:            100.0,
			is_applicable:          true,
			is_unloq_discount_code: true,
			reason:                 "",
		},
		{
			name:                   "valid unloq code with hashed phone but no order amount",
			domain:                 "example.com",
			phoneNumber:            utils.GetSHA256Hash("+918989898989"),
			giftCardCode:           "FLX123415EG",
			is_applicable:          true,
			is_unloq_discount_code: true,
			reason:                 "",
		},
		{
			name:                   "valid unloq code with hashed phone but empty domain",
			domain:                 "",
			phoneNumber:            utils.GetSHA256Hash("+918989898989"),
			giftCardCode:           "FLX123415EG",
			is_applicable:          false,
			is_unloq_discount_code: true,
			reason:                 giftcard.ReasonInvalidDomain,
		},
		{
			name:                   "valid unloq code with hashed phone but invalid domain",
			domain:                 "example\x00.com",
			phoneNumber:            utils.GetSHA256Hash("+918989898989"),
			giftCardCode:           "FLX123415EG",
			is_applicable:          false,
			is_unloq_discount_code: true,
			reason:                 giftcard.ReasonInvalidDomain,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, gotReason := giftcard.ValidateUnloqDiscountCode(tt.domain, tt.phoneNumber, tt.giftCardCode, tt.orderAmount)
			if got != tt.is_applicable {
				t.Errorf("ValidateGiftCardCode() = %v, want %v", got, tt.is_applicable)
			}
			if got2 != tt.is_unloq_discount_code {
				t.Errorf("ValidateGiftCardCode() = %v, want %v", got2, tt.is_unloq_discount_code)
			}
			if gotReason != tt.reason {
				t.Errorf("ValidateGiftCardCode() = %v, want %v", gotReason, tt.reason)
			}
		})
	}
}
