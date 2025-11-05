package giftcard_test

import (
	"testing"

	"github.com/fealtyx/fealtyx-go/giftcard"
	"github.com/fealtyx/fealtyx-go/giftcard/utils"
)

func TestValidateGiftCardCode(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		domain       string
		phoneNumber  string
		giftCardCode string
		want         bool
		want2        giftcard.Reason
	}{
		{
			name: "valid code",
			domain: "example.com",
		    phoneNumber: "+918989898989",
		    giftCardCode: "FLX123415EG",
		    want: true,
		    want2: giftcard.ReasonValid,
		},
		{
			name: "invalid code - wrong suffix",
			domain: "example.com",
		    phoneNumber: "+918989898989",
		    giftCardCode: "FLX1234ABCD",
		    want: false,
		    want2: giftcard.ReasonInvalidFormat,
		},
		{
			name: "invalid code - invalid phone",
			domain: "example.com",
		    phoneNumber: "777777",
		    giftCardCode: "FLX123415EG",
		    want: false,
		    want2: giftcard.ReasonInvalidPhone,
		},
		{
			name: "invalid code - missing prefix",
			domain: "example.com",
		    phoneNumber: "+918989898989",
					    giftCardCode: "123415EG",
		    want: false,
		    want2: giftcard.ReasonInvalidFormat,			
		},
		{
			name: "invalid code - short code",
			domain: "example.com",
		    phoneNumber: "+918989898989",
		    giftCardCode: "FL1",
		    want: false,
		    want2: giftcard.ReasonTooShort,
		},
		{
			name :"valid code with hashed phone",
			domain: "example.com",
		    phoneNumber: utils.GetSHA256Hash("+918989898989"),
			giftCardCode: "FLX123415EG",
			want: true,
			want2: giftcard.ReasonValid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := giftcard.ValidateGiftCardCode(tt.domain, tt.phoneNumber, tt.giftCardCode)
			if got != tt.want {
				t.Errorf("ValidateGiftCardCode() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("ValidateGiftCardCode() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
