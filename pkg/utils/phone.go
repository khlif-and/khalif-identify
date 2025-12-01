package utils

import (
	"errors"

	"github.com/nyaruka/phonenumbers"

)

// FormatPhoneNumber memvalidasi dan mengubah format ke E.164 (Internasional)
// defaultRegion contoh: "ID", "US", "SG"
func FormatPhoneNumber(phone, defaultRegion string) (string, error) {
	// 1. Parse nomor
	num, err := phonenumbers.Parse(phone, defaultRegion)
	if err != nil {
		return "", errors.New("format nomor telepon tidak valid")
	}

	// 2. Validasi apakah nomor itu asli/mungkin ada
	if !phonenumbers.IsValidNumber(num) {
		return "", errors.New("nomor telepon tidak valid atau tidak terdaftar")
	}

	// 3. Format ke E.164 (Contoh: +62812345678)
	formattedPhone := phonenumbers.Format(num, phonenumbers.E164)
	
	return formattedPhone, nil
}