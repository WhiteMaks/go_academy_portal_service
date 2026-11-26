package util

import (
	"errors"
	"regexp"
	"strings"
)

var phoneStripPattern = regexp.MustCompile(`[\s\-]`)

var e164Pattern = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

// NormalizePhoneE164 strips spaces/dashes and checks E.164 (+[1-9], 8–15 digits total).
func NormalizePhoneE164(raw string) (string, error) {
	normalized := phoneStripPattern.ReplaceAllString(strings.TrimSpace(raw), "")
	if normalized == "" {
		return "", errors.New("phone is required")
	}
	if !e164Pattern.MatchString(normalized) {
		return "", errors.New("invalid phone format")
	}
	return normalized, nil
}

// IsPhoneE164 reports whether raw input is a valid E.164 phone after normalization.
func IsPhoneE164(raw string) bool {
	_, err := NormalizePhoneE164(raw)
	return err == nil
}
