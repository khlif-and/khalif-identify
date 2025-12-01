package utils

import (
	"fmt"
	"net/url"

)

// GenerateAvatarURL membuat link UI Avatars berdasarkan inisial nama/email
func GenerateAvatarURL(name, email string) string {
	cleanName := name
	if cleanName == "" {
		cleanName = email
	}
	
	// Encode agar spasi jadi %20 dsb
	encodedName := url.QueryEscape(cleanName)
	
	// Return URL
	return fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=random&color=fff&size=128", encodedName)
}