package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

func RandomUserAvatar(fullname string) string {
	return fmt.Sprintf("https://api.dicebear.com/6.x/initials/svg?seed=%s", fullname)
}

func GenerateOTP(length int) string {
	digits := "0123456789"
	var sb strings.Builder

	for range length {
		sb.WriteByte(digits[rand.Intn(len(digits))])
	}

	return sb.String()
}

func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32) // 256-bit token
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateSlug(input string) string {

	slug := strings.ToLower(input)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	suffix := strconv.Itoa(rand.Intn(1_000_000))
	slug = slug + "-" + leftPad(suffix, "0", 6)

	return slug
}

func leftPad(s string, pad string, length int) string {
	for len(s) < length {
		s = pad + s
	}
	return s
}
