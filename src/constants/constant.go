package constants

/* Go time formatter: Mon Jan 2 15:04:05 MST 2006 */
const (
	BA_DateTimeLayout = "2006-01-02T15:04:05-0700"
	GA_DateTimeLayout = "2006-01-02T15:04:05-07:00"
	LA_DateTimeLayout = "2006-01-02T15:04:05"
)

var (
	Cities = map[string]string{
		"CGK": "Jakarta",
		"DPS": "Denpasar",
		"SOC": "Kabupaten Boyolali",
		"SUB": "Surabaya",
		"UPG": "Kabupaten Maros",
	}

	Timezones = map[string]string{
		"CGK": "Asia/Jakarta",
		"DPS": "Asia/Makassar",
		"SOC": "Asia/Jakarta",
		"SUB": "Asia/Jakarta",
		"UPG": "Asia/Jakarta",
	}
)
