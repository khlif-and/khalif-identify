package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/url"

)

// ProcessProfileImageResult menampung hasil pemrosesan gambar
type ProcessProfileImageResult struct {
	DominantColor string
	AvatarURL     string // Terisi hanya jika pakai UI Avatar
}

// HandleProfileImageLogic menentukan warna dominan dan avatar URL
func HandleProfileImageLogic(file multipart.File, name, email string) (ProcessProfileImageResult, error) {
	var result ProcessProfileImageResult

	if file != nil {
		// A. KASUS ADA FILE: Ekstrak Warna
		extractedColor, err := ExtractDominantColor(file)
		if err == nil {
			result.DominantColor = extractedColor
		} else {
			result.DominantColor = "#000000" // Fallback
		}

		// Reset pointer file agar bisa diupload nanti oleh Azure
		file.Seek(0, io.SeekStart)

	} else {
		// B. KASUS TIDAK ADA FILE: Generate Random
		randomHex := GenerateRandomHexColor()
		result.DominantColor = "#" + randomHex

		cleanName := name
		if cleanName == "" {
			cleanName = email
		}
		encodedName := url.QueryEscape(cleanName)

		// Set URL Avatar
		result.AvatarURL = fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=%s&color=fff&size=128", encodedName, randomHex)
	}

	return result, nil
}