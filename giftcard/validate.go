package giftcard

import (
	"fmt"
	"net/url"
	"strings"

	utils "github.com/fealtyx/unloq-go/giftcard/utils"
)

// Reason describes why a validation passed or failed.
type Reason string

const (
	ReasonEmptyInput         Reason = "Empty Input"
	ReasonTooShort           Reason = "Code Too Short"
	ReasonVoucherNotEligible Reason = "Voucher Not Eligible For The User"
	ReasonInvalidPhone       Reason = "Invalid Phone Number"
	ReasonInvalidDomain      Reason = "Invalid Domain"
	ReasonInvalidOrderAmount Reason = "Invalid Order Amount"
)

func (r Reason) String() string {
	if r == "" {
		return string("")
	}
	return string(r)
}

func ValidateUnloqDiscountCode(domain, phoneNumber, giftCardCode string, orderAmount float64) (is_applicable bool, is_unloq_discount_code bool, reason Reason) {
	// First, clean the gift card code
	giftCardCode = strings.ToLower(strings.TrimSpace(giftCardCode))

	// If gift card code is too short, return
	if len(giftCardCode) < 4 {
		return true, false, ""
	}

	// Check if gift card code starts with "flx" or "unq"
	if !strings.HasPrefix(giftCardCode, "flx") && !strings.HasPrefix(giftCardCode, "unq") {
		return true, false, ""
	}

	// Get last 4 characters of the code
	last4 := giftCardCode[len(giftCardCode)-4:]
	if len(last4) != 4 {
		return true, false, ""
	}

	first2 := last4[0:2]
	actualLast2 := last4[2:]

	codePrefix := giftCardCode[0 : len(giftCardCode)-4]

	// Step 1: Check if this is an Unloq gift card by checking if the last 2 characters match the expected last 2 characters
	expectedLast2 := generateLast2FromCodePrefix(codePrefix)
	if actualLast2 != expectedLast2 {
		// If the last 2 don't match the expected pattern, the gift card is not an Unloq gift card
		return true, false, ""
	}

	// // Sanitize domain
	// domain, err := sanitizeDomain(domain)
	// if err != nil {
	// 	return false, true, ReasonInvalidDomain
	// }

	// Step 2: Verify first 2 characters match expected value from phone number
	expectedFirst2 := strings.ToLower(getGiftCardCodeIdentifier(phoneNumber))

	// Case-insensitive match
	valid := strings.EqualFold(expectedFirst2, first2)
	if !valid {
		return false, true, ReasonVoucherNotEligible
	}

	return true, true, ""
}

// Helper function to generate last 2 characters from first 2 characters
func generateLast2FromFirst2(first2 string) string {
	const alphanum = "abcdefghijklmnopqrstuvwxyz0123456789"
	const charsetLen = len(alphanum)

	if len(first2) != 2 {
		return ""
	}

	// Convert hex characters to numeric byte values
	c1 := first2[0]
	c2 := first2[1]

	b1 := utils.HexCharToByte(c1)
	b2 := utils.HexCharToByte(c2)

	// Get XOR and sum of first 2 characters
	xor := b1 ^ b2
	sum := b1 + b2

	// Map XOR and sum results to character set
	idx1 := int(xor) % charsetLen
	idx2 := int(sum) % charsetLen

	return string([]byte{alphanum[idx1], alphanum[idx2]})
}

func generateLast2FromCodePrefix(codePrefix string) string {
	const alphanum = "abcdefghijklmnopqrstuvwxyz0123456789"
	const charsetLen = len(alphanum)

	if len(codePrefix) == 0 {
		return ""
	}

	prefixCodeHash := utils.GetSHA256Hash(strings.ToLower(codePrefix))

	b1 := prefixCodeHash[0]
	b2 := prefixCodeHash[1]

	// Get xor and sum of first 2 characters
	xor := b1 ^ b2
	sum := b1 + b2

	// Map xor and sum results to character set
	idx1 := xor % byte(charsetLen)
	idx2 := sum % byte(charsetLen)

	return string([]byte{alphanum[idx1], alphanum[idx2]})
}

// generates a deterministic gift card code identifier using phone number
func getGiftCardCodeIdentifier(phoneNumber string) string {

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
	input := processedPhoneNumber
	// Convert to lower case
	input = strings.ToLower(input)

	hash := utils.GetSHA256Hash(input)

	// First 2 hex characters of the hash from phone number
	c1 := hash[0]
	c2 := hash[1]

	return fmt.Sprintf("%c%c", c1, c2)
}

func sanitizeDomain(input string) (string, error) {
	input = strings.TrimSpace(input)

	// Add scheme if missing to make it a valid URL
	if !strings.Contains(input, "://") {
		input = "https://" + input
	}

	parsed, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("invalid domain: %w", err)
	}

	// Extract host and convert to lowercase
	domain := strings.ToLower(parsed.Hostname())
	if domain == "" {
		return "", fmt.Errorf("empty domain")
	}

	return domain, nil
}
