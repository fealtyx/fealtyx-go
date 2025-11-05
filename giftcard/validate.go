package giftcard

import (
	"fmt"
	"regexp"
	"strings"
    utils "github.com/fealtyx/fealtyx-go/giftcard/utils"
)

// Reason describes why a validation passed or failed.
type Reason string

const (
    ReasonUnknown        Reason = "unknown"
    ReasonValid          Reason = "valid"
    ReasonEmptyInput     Reason = "empty_input"
	ReasonTooShort      Reason = "too_short_code"
    ReasonInvalidFormat  Reason = "invalid_format"
    ReasonInvalidPhone   Reason = "invalid_phone"
)

func (r Reason) String() string {
    if r == "" {
        return string(ReasonUnknown)
    }
    return string(r)
}

func ValidateGiftCardCode(domain, phoneNumber, giftCardCode string) (bool, Reason) {
	if len(giftCardCode) < 4 {
		// code too short to contain a valid suffix
		return false, ReasonInvalidFormat
	}

	if !checkValidPhoneNumber(phoneNumber) {
		return false, ReasonInvalidPhone
	}

	expectedSuffix := strings.ToLower(getGiftCardCodeIdentifier(domain, phoneNumber))
	actualSuffix := strings.ToLower(giftCardCode[len(giftCardCode)-4:])

	// Case-insensitive match since generation uppercases code
	valid := strings.EqualFold(expectedSuffix, actualSuffix)

	if !valid {
		return false, ReasonInvalidFormat
	}

	return true, ReasonValid
}


func getGiftCardCodeIdentifier(domain, phoneNumber string) string {

	var processedPhoneNumber string

	//check if phone number is already hashed
	if len(phoneNumber) == 64  && utils.IsSHA256Hash(phoneNumber) {
		processedPhoneNumber = phoneNumber
	} else {
		if !strings.HasPrefix(phoneNumber, "+91") {
			// If the phone number doesn't start with +91, add it
			phoneNumber = "+91" + phoneNumber
		}
		// Hash the phone number using SHA-256
		processedPhoneNumber = utils.GetSHA256Hash(phoneNumber)
	}

	// Character set for alphanumeric (A-Z + 0-9)
	const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const charsetLen = len(alphanum)
	input := domain + processedPhoneNumber
	// Convert to lower case
	input = strings.ToLower(input)

	hash := utils.GetSHA256Hash(input)

	// First 2 hex characters of the hash
	c1 := hash[0]
	c2 := hash[1]

	// Convert to 0–15 values
	b1 := utils.HexCharToByte(c1)
	b2 := utils.HexCharToByte(c2)

	// Derive two new characters from the first two
	xor := b1 ^ b2
	idx1 := int(xor) % charsetLen
	idx2 := int(b1+b2) % charsetLen

	return fmt.Sprintf("%c%c%c%c", c1, c2, alphanum[idx1], alphanum[idx2])
}


func checkValidPhoneNumber(phone string) bool {
	phone = strings.TrimSpace(phone)


	if utils.IsSHA256Hash(phone) {
		return true
	}

	// Case 1: +91 + 10 digits → total 13 chars
	if strings.HasPrefix(phone, "+91") {
		re := regexp.MustCompile(`^\+91[1-9]\d{9}$`)
		return re.MatchString(phone)
	}

	// Case 2: plain 10 digits (no +91)
	re := regexp.MustCompile(`^[1-9]\d{9}$`)
	return re.MatchString(phone)
}
