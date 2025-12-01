package utils

import (
	"fmt"
	"image"
	_ "image/jpeg" // Support decode JPEG
	_ "image/png"  // Support decode PNG
	"mime/multipart"

	"github.com/EdlinOrg/prominentcolor"

)

func ExtractDominantColor(file multipart.File) (string, error) {
	// 1. Decode gambar dari stream file
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// 2. Proses K-Means (Mencari 1 warna paling dominan)
	// Parameter: k=3 (jumlah cluster), mask=prominentcolor.ArgumentNoCropping
	cols, err := prominentcolor.Kmeans(img)
	if err != nil {
		return "", err
	}

	if len(cols) == 0 {
		return "#000000", nil // Fallback Black
	}

	// 3. Ambil warna pertama (paling dominan) dan format ke Hex
	bestColor := cols[0]
	hexColor := fmt.Sprintf("#%02x%02x%02x", bestColor.Color.R, bestColor.Color.G, bestColor.Color.B)

	return hexColor, nil
}