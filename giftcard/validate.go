package giftcard

import (
	"regexp"
	"strings"
	"time"
)

// Reason describes why a validation passed or failed.
type Reason string

const (
    ReasonUnknown        Reason = "unknown"
    ReasonValid          Reason = "valid"
    ReasonEmptyInput     Reason = "empty_input"
    ReasonInvalidFormat  Reason = "invalid_format"
    ReasonNotFound       Reason = "not_found"
    ReasonPhoneMismatch  Reason = "phone_mismatch"
    ReasonExpired        Reason = "expired"
    ReasonAlreadyRedeemed Reason = "already_redeemed"
)

func (r Reason) String() string {
    if r == "" {
        return string(ReasonUnknown)
    }
    return string(r)
}


type GiftCard struct {
	Code     string
	Phone    string // plain phone number or hashed phone number
}

var codeFormat = regexp.MustCompile(`^[A-Z0-9]{6,12}$`)


func ValidateCode(phone, code string) (bool, Reason) {
	// basic empties
	if strings.TrimSpace(phone) == "" || strings.TrimSpace(code) == "" {
		return false, ReasonEmptyInput
	}

	code = strings.ToUpper(strings.TrimSpace(code))

	if !codeFormat.MatchString(code) {
		return false, ReasonInvalidFormat
	}

	

	return true, ReasonValid
}

// normalizePhoneToDigits strips everything except digits and returns digits-only
func normalizePhoneToDigits(p string) string {
	var b strings.Builder
	for _, r := range p {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
