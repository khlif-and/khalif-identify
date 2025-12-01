package utils

import (
	"fmt"

	"github.com/biter777/countries"

)

type Country struct {
	Name     string `json:"name"`
	ISOCode  string `json:"iso_code"`  // ID, US, SG
	DialCode string `json:"dial_code"` // +62, +1
}

func GetCountryList() []Country {
	var list []Country

	// Library ini punya daftar semua negara (All)
	allCountries := countries.All()

	for _, c := range allCountries {
		// Validasi: Hanya ambil negara yang punya kode telepon valid
		// (Menghindari wilayah antartika atau pulau kosong)
		if c.CallCodes() != nil && len(c.CallCodes()) > 0 {
			
			// Ambil Call Code pertama (biasanya negara cuma punya 1 kode utama)
			dialCode := c.CallCodes()[0]

			list = append(list, Country{
				Name:     c.String(),      // Contoh: "Indonesia"
				ISOCode:  c.Alpha2(),      // Contoh: "ID"
				DialCode: fmt.Sprintf("+%s", dialCode), // Format: "+62"
			})
		}
	}

	return list
}