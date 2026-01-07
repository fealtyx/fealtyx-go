package giftcard

import (
	"fmt"
	"net/url"
	"strings"

	utils "github.com/fealtyx/fealtyx-go/giftcard/utils"
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

func ValidateUnloqDiscountCode(domain, phoneNumber, giftCardCode string, orderAmount float64) (is_applicable bool, is_fealtyx_discount_code bool, reason Reason) {
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

	// Step 1: Check if this is a Fealtyx gift card by checking if the last 2 characters match the expected last 2 characters
	expectedLast2 := generateLast2FromFirst2(first2)
	if actualLast2 != expectedLast2 {
		// If the last 2 don't match the expected pattern, the gift card is not a Fealtyx gift card
		return true, false, ""
	}

	// Sanitize domain
	domain, err := sanitizeDomain(domain)
	if err != nil {
		return false, true, ReasonInvalidDomain
	}

	// Step 2: Verify full code - generate expected suffix
	expectedSuffix := strings.ToLower(getGiftCardCodeIdentifier(domain, phoneNumber))

	// Case-insensitive match
	valid := strings.EqualFold(expectedSuffix, last4)
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

func getGiftCardCodeIdentifier(domain, phoneNumber string) string {

	hashedPhoneNumber := utils.GetPhoneNumberHash(phoneNumber)
	// Character set for alphanumeric (A-Z + 0-9)
	const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const charsetLen = len(alphanum)
	input := domain + hashedPhoneNumber
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
