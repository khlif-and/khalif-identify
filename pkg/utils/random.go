package utils

import (
	"fmt"
	"math/rand"
	"time"

)

func GenerateRandomHexColor() string {
	rand.Seed(time.Now().UnixNano())
	// Generate RGB Random
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	
	// Format ke Hex (Contoh: FF5733)
	// Kita return TANPA pagar (#) dulu agar mudah dipakai di URL UI Avatars
	return fmt.Sprintf("%02x%02x%02x", r, g, b)
}