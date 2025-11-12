package giftcard

import (
	"fmt"
	utils "github.com/fealtyx/fealtyx-go/giftcard/utils"
	"regexp"
	"strings"
)

// Reason describes why a validation passed or failed.
type Reason string

const (
	ReasonEmptyInput             Reason = "Empty Input"
	ReasonTooShort               Reason = "Code Too Short"
	ReasonVoucherNotEligible     Reason = "Voucher Not Eligible For The User"
	ReasonInvalidPhone           Reason = "Invalid Phone"
	ReasonInvalidPartnerEntityId Reason = "Invalid Partner Entity Id"
	ReasonInvalidOrderAmount     Reason = "Invalid Order Amount"
)

func (r Reason) String() string {
	if r == "" {
		return string("")
	}
	return string(r)
}

func ValidateGiftCardCode(domain, phoneNumber, giftCardCode string, orderAmount float64) (is_applicable bool, is_fealtyx_discount_code bool, reason Reason) {
	if len(giftCardCode) < 4 {
		// code too short to contain a valid suffix
		return true, false, ""
	}

	giftCardCode = strings.ToLower(giftCardCode)

	if !strings.HasPrefix(giftCardCode, "flx") {
		return true, false, ""
	}

	if !checkValidPhoneNumber(phoneNumber) {
		return true, false, ReasonInvalidPhone
	}

	if domain==""{
		return true, false , ReasonInvalidPartnerEntityId
	}

	expectedSuffix := strings.ToLower(getGiftCardCodeIdentifier(domain, phoneNumber))
	actualSuffix := strings.ToLower(giftCardCode[len(giftCardCode)-4:])

	// Case-insensitive match since generation uppercases code
	valid := strings.EqualFold(expectedSuffix, actualSuffix)

	if !valid {
		return false, true, ReasonVoucherNotEligible
	}

	return true, true, ""
}

func getGiftCardCodeIdentifier(domain, phoneNumber string) string {

	var processedPhoneNumber string

	//check if phone number is already hashed
	if utils.IsSHA256Hash(phoneNumber) {
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
