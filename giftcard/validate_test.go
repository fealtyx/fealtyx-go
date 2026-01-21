package giftcard_test

import (
	"testing"

	"github.com/fealtyx/unloq-go/giftcard"
	"github.com/fealtyx/unloq-go/giftcard/utils"
)

func TestValidateGiftCardCode(t *testing.T) {
	tests := []struct {
		name                string
		domain              string
		phoneNumber         string
		giftCardCode        string
		orderAmount         float64
		isApplicable        bool
		isUnloqDiscountCode bool
		reason              giftcard.Reason
	}{
		{
			name:                "invalid fealtyx code, first 2 of last 4 characters don't match last 2",
			domain:              "example.com",
			phoneNumber:         "+918989898989",
			giftCardCode:        "FLX1234ABCD",
			orderAmount:         100.0,
			isApplicable:        true,
			isUnloqDiscountCode: false,
			reason:              "",
		},
		{
			name:                "invalid fealtyx code - missing prefix",
			domain:              "example.com",
			phoneNumber:         "+918989898989",
			giftCardCode:        "123415EG",
			orderAmount:         100.0,
			isApplicable:        true,
			isUnloqDiscountCode: false,
			reason:              "",
		},
		{
			name:                "invalid fealtyx code - short code",
			domain:              "example.com",
			phoneNumber:         "+918989898989",
			giftCardCode:        "FL1",
			orderAmount:         100.0,
			isApplicable:        true,
			isUnloqDiscountCode: false,
			reason:              "",
		},
		// {
		//     name:                "valid fealtyx code - domain mismatch",
		//     domain:              "example2.com",
		//     phoneNumber:         "+918989898989",
		//     giftCardCode:        "FLX123494NN",
		//     orderAmount:         100.0,
		//     isApplicable:        false,
		//     isUnloqDiscountCode: true,
		//     reason:              giftcard.ReasonVoucherNotEligible,
		// },
		{
			name:                "valid fealtyx code - phone mismatch",
			domain:              "example.com",
			phoneNumber:         "+918989898989",
			giftCardCode:        "FLX1234ase2",
			orderAmount:         100.0,
			isApplicable:        false,
			isUnloqDiscountCode: true,
			reason:              giftcard.ReasonVoucherNotEligible,
		},
		{
			name:                "fealtyx deprecated code with FLX prefix, should be applicable",
			domain:              "example.com",
			phoneNumber:         "+918989898989",
			giftCardCode:        "FLX123415EG",
			orderAmount:         100.0,
			isApplicable:        true,
			isUnloqDiscountCode: false,
			reason:              "",
		},
		{
			name:                "fealtyx deprecated code with UNQ prefix, should be applicable",
			domain:              "example.com",
			phoneNumber:         "+918989898989",
			giftCardCode:        "UNQ123415EG",
			orderAmount:         100.0,
			isApplicable:        true,
			isUnloqDiscountCode: false,
			reason:              "",
		},
		{
			name:                "valid fealtyx code with FLX prefix",
			domain:              "example.com",
			phoneNumber:         "+918989898989",
			giftCardCode:        "FLX123494e2",
			orderAmount:         100.0,
			isApplicable:        true,
			isUnloqDiscountCode: true,
			reason:              "",
		},
		{
			name:                "valid fealtyx code with UNQ prefix",
			domain:              "example.com",
			phoneNumber:         "+918989898989",
			giftCardCode:        "UNQ123494me",
			orderAmount:         100.0,
			isApplicable:        true,
			isUnloqDiscountCode: true,
			reason:              "",
		},
		{
			name:                "valid fealtyx code with hashed phone",
			domain:              "example.com",
			phoneNumber:         utils.GetSHA256Hash("+918989898989"),
			giftCardCode:        "FLX123494e2",
			orderAmount:         100.0,
			isApplicable:        true,
			isUnloqDiscountCode: true,
			reason:              "",
		},
		{
			name:                "valid fealtyx code with non hashed phone",
			domain:              "example.com",
			phoneNumber:         "+918989898989",
			giftCardCode:        "FLX123494e2",
			orderAmount:         100.0,
			isApplicable:        true,
			isUnloqDiscountCode: true,
			reason:              "",
		},
		{
			name:                "valid fealtyx code with hashed phone but no order amount",
			domain:              "example.com",
			phoneNumber:         utils.GetSHA256Hash("+918989898989"),
			giftCardCode:        "FLX123494e2",
			orderAmount:         0.0,
			isApplicable:        true,
			isUnloqDiscountCode: true,
			reason:              "",
		},
		// {
		//     name:                "valid fealtyx code with hashed phone but empty domain",
		//     domain:              "",
		//     phoneNumber:         utils.GetSHA256Hash("+918989898989"),
		//     giftCardCode:        "FLX123494e2",
		//     orderAmount:         100.0,
		//     isApplicable:        false,
		//     isUnloqDiscountCode: true,
		//     reason:              giftcard.ReasonInvalidDomain,
		// },
		// {
		//     name:                "valid fealtyx code with hashed phone but invalid domain",
		//     domain:              "example\x00.com",
		//     phoneNumber:         utils.GetSHA256Hash("+918989898989"),
		//     giftCardCode:        "FLX123494NN",
		//     orderAmount:         100.0,
		//     isApplicable:        false,
		//     isUnloqDiscountCode: true,
		//     reason:              giftcard.ReasonInvalidDomain,
		// },
		{
			name:                "non-fealtyx code",
			domain:              "example.com",
			phoneNumber:         "+918989898989",
			giftCardCode:        "ABCD1234EFGH",
			orderAmount:         100.0,
			isApplicable:        true,
			isUnloqDiscountCode: false,
			reason:              "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotApplicable, gotIsUnloqCode, gotReason := giftcard.ValidateUnloqDiscountCode(
				tt.domain,
				tt.phoneNumber,
				tt.giftCardCode,
				tt.orderAmount,
			)

			if gotApplicable != tt.isApplicable {
				t.Errorf("ValidateUnloqDiscountCode() applicable = %v, want %v", gotApplicable, tt.isApplicable)
			}

			if gotIsUnloqCode != tt.isUnloqDiscountCode {
				t.Errorf("ValidateUnloqDiscountCode() isUnloqDiscountCode = %v, want %v", gotIsUnloqCode, tt.isUnloqDiscountCode)
			}

			if gotReason != tt.reason {
				t.Errorf("ValidateUnloqDiscountCode() reason = %v, want %v", gotReason, tt.reason)
			}
		})
	}
}
