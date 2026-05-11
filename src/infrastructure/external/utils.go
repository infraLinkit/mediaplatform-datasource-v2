package external

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Concat(splitcode string, args ...string) string {

	var b bytes.Buffer

	for _, arg := range args {
		b.WriteString(arg + splitcode)
	}

	f := strings.TrimRight(b.String(), splitcode)
	return f
}

func GetFormatTime(loc *time.Location, layout string) string {

	// Standard GO Constant Format :

	// ANSIC       = "Mon Jan _2 15:04:05 2006"
	// UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	// RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	// RFC822      = "02 Jan 06 15:04 MST"
	// RFC822Z     = "02 Jan 06 15:04 -0700"
	// RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	// RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	// RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700"
	// RFC3339     = "2006-01-02T15:04:05Z07:00"
	// RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	// Kitchen     = "3:04PM"
	// // Handy time stamps.
	// Stamp      = "Jan _2 15:04:05"
	// StampMilli = "Jan _2 15:04:05.000"
	// StampMicro = "Jan _2 15:04:05.000000"
	// StampNano  = "Jan _2 15:04:05.000000000"

	// Using Manual Format :
	// 1. date yyyy-mm-dd = 2006-01-02
	// 2. time hhhh:ii:ss = 15:04:05

	t := time.Now()
	f := t.In(loc).Format(layout)

	return f
}

func GetCurrentTime(loc *time.Location, val string) time.Time {

	now := time.Now()
	t, _ := time.Parse(val, now.In(loc).Format(time.RFC3339))
	return t
}

func GetYesterday(loc *time.Location, day time.Duration, layout string) time.Time {

	now := time.Now()
	t, _ := time.Parse(layout, now.In(loc).Format(layout))
	return t.In(loc).Add((-24 * day) * time.Hour)
}

func GetTomorrow(loc *time.Location, day time.Duration, layout string) time.Time {

	now := time.Now()
	t, _ := time.Parse(layout, now.In(loc).Format(layout))
	return t.In(loc).Add((24 * day) * time.Hour)
}

func CounterZeroNumber(length int) string {

	var wordNumbers string

	for w := 0; w < length; w++ {
		wordNumbers += "0"
	}

	return wordNumbers
}

func ReduceWords(words string, start int, length int) string {

	runes := []rune(words)
	inputFmt := string(runes[start:length])

	return inputFmt
}

func GetUniqId(loc *time.Location) string {

	t := time.Now()
	var formatId = t.In(loc).Format("20060102150405.000000")
	uniqId := strings.Replace(formatId, ".", "", -1)

	return uniqId
}

func GetTrxID(loc *time.Location) string {

	t := time.Now()
	var formatId = t.In(loc).Format("20060102150405")

	return ReduceWords(strings.Replace(formatId, ".", "", -1)+strconv.Itoa(int(t.In(loc).UnixNano())), 0, 32)
}

func GetIpAddress(c *fiber.Ctx) string {

	ip := ""

	for k, v := range c.GetReqHeaders() {
		//fmt.Printf("(k1) %#v", k)

		if k == "Cf-Connecting-Ip" {
			for _, v2 := range v {
				//fmt.Printf("(2) %#v", v2)
				if v2 != "" {
					ip = v2
					break
				}
			}
			break
		}
	}

	if ip == "" {
		for k, v := range c.GetReqHeaders() {
			//fmt.Printf("(k1) %#v", k)

			if k == "X-Forwarded-For" {
				for _, v2 := range v {
					//fmt.Printf("(2) %#v", v2)
					if v2 != "" {
						ip = v2
						break
					}
				}
				break
			}
		}
	}

	if ip == "" {
		for k, v := range c.GetReqHeaders() {
			//fmt.Printf("(k1) %#v", k)

			if k == "X-Real-Ip" {
				for _, v2 := range v {
					//fmt.Printf("(2) %#v", v2)
					if v2 != "" {
						ip = v2
						break
					}
				}
				break
			}
		}
	}

	if ip == "" {
		ip = c.IP()
	}

	return ip
}

// AES_SECRET_KEY env: 16/24/32 byte → AES-128/192/256.
func getAESKey() ([]byte, error) {
	k := os.Getenv("AES_SECRET_KEY")
	if k == "" {
		return nil, errors.New("AES_SECRET_KEY env not set")
	}
	switch len(k) {
	case 16, 24, 32:
		return []byte(k), nil
	default:
		return nil, fmt.Errorf("AES_SECRET_KEY must be 16/24/32 bytes, got %d", len(k))
	}
}

func Encrypt(plaintext string) (string, error) {
	key, err := getAESKey()
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("nonce: %w", err)
	}
	ct := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ct), nil
}

func Decrypt(ciphertext string) (string, error) {
	key, err := getAESKey()
	if err != nil {
		return "", err
	}
	raw, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm: %w", err)
	}
	ns := gcm.NonceSize()
	if len(raw) < ns {
		return "", errors.New("ciphertext too short")
	}
	nonce, ct := raw[:ns], raw[ns:]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", fmt.Errorf("gcm open: %w", err)
	}
	return string(pt), nil
}

func CompressGzip(f string) string {

	// Open the original file
	originalFile, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer originalFile.Close()

	gzf := f + ".gz"
	// Create a new gzipped file
	gzippedFile, err := os.Create(gzf)
	if err != nil {
		panic(err)
	}
	defer gzippedFile.Close()

	// Create a new gzip writer
	gzipWriter := gzip.NewWriter(gzippedFile)
	defer gzipWriter.Close()

	// Copy the contents of the original file to the gzip writer
	_, err = io.Copy(gzipWriter, originalFile)
	if err != nil {
		panic(err)
	}

	// Flush the gzip writer to ensure all data is written
	gzipWriter.Flush()

	return gzf
}

func GetDateFromInt(loc *time.Location, dateInt int64) time.Time {
	t := time.Unix(dateInt, 0)
	return t.In(loc)
}

func ParseVendorDateSend(loc *time.Location, raw string) time.Time {

	fallback := time.Date(1970, 1, 1, 12, 0, 0, 0, loc)

	if len(raw) < 8 {
		return fallback
	}

	datePart := raw[:8]

	t, err := time.ParseInLocation("20060102", datePart, loc)
	if err != nil {
		return fallback
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		12, 0, 0, 0,
		loc,
	)
}

func InArray(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func AppendPartnerToTxt(filePath, content string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content + "\n"); err != nil {
		return err
	}

	return f.Sync()
}

func NormalizeCountry(country string) []string {
	switch country {
	case "SA", "KSA":
		return []string{"SA", "KSA"}
	case "LA", "LS":
		return []string{"LA", "LS"}
	case "LK", "LKA":
		return []string{"LK", "LKA"}
	case "PS", "PSE":
		return []string{"PS", "PSE"}
	case "SE", "SLE":
		return []string{"SE", "SLE"}
	case "NG", "NGA":
		return []string{"NG", "NGA"}
	case "CZ", "CZE":
		return []string{"CZ", "CZE"}
	case "OM", "OMN":
		return []string{"OM", "OMN"}
	default:
		return []string{country}
	}
}

func NormalizeCountries(country string) string {
	country = strings.ToUpper(country)

	customMap := map[string]string{
		"UAE": "AE",
		"KSA": "SA",
		"LKA": "LK",
		"NGA": "NG",
		"SLE": "SL",
		"OMN": "OM",
		"PSE": "PS",
		"CZE": "CZ",
		"ARE": "AE",
		"SAU": "SA",
	}

	if val, ok := customMap[country]; ok {
		return val
	}

	alpha3to2 := map[string]string{
		"AFG": "AF", "ALB": "AL", "DZA": "DZ", "AND": "AD", "AGO": "AO",
		"ARG": "AR", "ARM": "AM", "AUS": "AU", "AUT": "AT", "AZE": "AZ",
		"BHR": "BH", "BGD": "BD", "BLR": "BY", "BEL": "BE", "BEN": "BJ",
		"BTN": "BT", "BOL": "BO", "BIH": "BA", "BWA": "BW", "BRA": "BR",
		"BRN": "BN", "BGR": "BG", "BFA": "BF", "BDI": "BI", "CPV": "CV",
		"KHM": "KH", "CMR": "CM", "CAN": "CA", "CAF": "CF", "TCD": "TD",
		"CHL": "CL", "CHN": "CN", "COL": "CO", "COM": "KM", "COG": "CG",
		"COD": "CD", "CRI": "CR", "CIV": "CI", "HRV": "HR", "CUB": "CU",
		"CYP": "CY", "CZE": "CZ", "DNK": "DK", "DJI": "DJ", "DOM": "DO",
		"ECU": "EC", "EGY": "EG", "SLV": "SV", "GNQ": "GQ", "ERI": "ER",
		"EST": "EE", "SWZ": "SZ", "ETH": "ET", "FJI": "FJ", "FIN": "FI",
		"FRA": "FR", "GAB": "GA", "GMB": "GM", "GEO": "GE", "DEU": "DE",
		"GHA": "GH", "GRC": "GR", "GRD": "GD", "GTM": "GT", "GIN": "GN",
		"GNB": "GW", "GUY": "GY", "HTI": "HT", "HND": "HN", "HUN": "HU",
		"ISL": "IS", "IND": "IN", "IDN": "ID", "IRN": "IR", "IRQ": "IQ",
		"IRL": "IE", "ISR": "IL", "ITA": "IT", "JAM": "JM", "JPN": "JP",
		"JOR": "JO", "KAZ": "KZ", "KEN": "KE", "KIR": "KI", "KWT": "KW",
		"KGZ": "KG", "LAO": "LA", "LVA": "LV", "LBN": "LB", "LSO": "LS",
		"LBR": "LR", "LBY": "LY", "LIE": "LI", "LTU": "LT", "LUX": "LU",
		"MDG": "MG", "MWI": "MW", "MYS": "MY", "MDV": "MV", "MLI": "ML",
		"MLT": "MT", "MHL": "MH", "MRT": "MR", "MUS": "MU", "MEX": "MX",
		"FSM": "FM", "MDA": "MD", "MCO": "MC", "MNG": "MN", "MNE": "ME",
		"MAR": "MA", "MOZ": "MZ", "MMR": "MM", "NAM": "NA", "NRU": "NR",
		"NPL": "NP", "NLD": "NL", "NZL": "NZ", "NIC": "NI", "NER": "NE",
		"NGA": "NG", "PRK": "KP", "MKD": "MK", "NOR": "NO", "OMN": "OM",
		"PAK": "PK", "PLW": "PW", "PSE": "PS", "PAN": "PA", "PNG": "PG",
		"PRY": "PY", "PER": "PE", "PHL": "PH", "POL": "PL", "PRT": "PT",
		"QAT": "QA", "ROU": "RO", "RUS": "RU", "RWA": "RW", "KNA": "KN",
		"LCA": "LC", "VCT": "VC", "WSM": "WS", "SMR": "SM", "STP": "ST",
		"SAU": "SA", "SEN": "SN", "SRB": "RS", "SYC": "SC", "SLE": "SL",
		"SGP": "SG", "SVK": "SK", "SVN": "SI", "SLB": "SB", "SOM": "SO",
		"ZAF": "ZA", "KOR": "KR", "SSD": "SS", "ESP": "ES", "LKA": "LK",
		"SDN": "SD", "SUR": "SR", "SWE": "SE", "CHE": "CH", "SYR": "SY",
		"TJK": "TJ", "TZA": "TZ", "THA": "TH", "TLS": "TL", "TGO": "TG",
		"TON": "TO", "TTO": "TT", "TUN": "TN", "TUR": "TR", "TKM": "TM",
		"TUV": "TV", "UGA": "UG", "UKR": "UA", "ARE": "AE", "GBR": "GB",
		"USA": "US", "URY": "UY", "UZB": "UZ", "VUT": "VU", "VEN": "VE",
		"VNM": "VN", "YEM": "YE", "ZMB": "ZM", "ZWE": "ZW",
	}

	if val, ok := alpha3to2[country]; ok {
		return val
	}

	if len(country) == 2 {
		return country
	}

	if len(country) >= 2 {
		return country[:2]
	}

	return country
}

var channelGroupMap = map[string]string{
	// TikTok
	"tiktok":                "Mainstream TikTok Ads",
	"mainstream_tiktok":     "Mainstream TikTok Ads",
	"mainstream_tiktok_ads": "Mainstream TikTok Ads",
	"tiktok ads":            "Mainstream TikTok Ads",
	// Google
	"google":                "Mainstream Google Ads",
	"mainstream_google":     "Mainstream Google Ads",
	"mainstream_google_ads": "Mainstream Google Ads",
	"google ads":            "Mainstream Google Ads",
	"google traffic":        "Mainstream Google Ads",
	// Meta
	"meta":                "Mainstream Meta Ads",
	"mainstream_meta":     "Mainstream Meta Ads",
	"mainstream_meta_ads": "Mainstream Meta Ads",
	"facebook":            "Mainstream Meta Ads",
	"fb":                  "Mainstream Meta Ads",
	"fbmeta":              "Mainstream Meta Ads",
	// Snack Video
	"snack video":          "Mainstream Snack Video Ads",
	"snack_video":          "Mainstream Snack Video Ads",
	"snackvideo":           "Mainstream Snack Video Ads",
	// Others
	"cpa":           "CPA",
	"dsp":           "DSP",
	"sms":           "SMS",
	"telco":         "Telco Channel",
	"telco_channel": "Telco Channel",
	"telco channel": "Telco Channel",
	"s2s":           "S2S",
	"api":           "API",
}

func ResolveChannelGroup(raw string) string {
	key := strings.ToLower(strings.TrimSpace(raw))
	if key == "" || key == "na" {
		return "Other"
	}
	if group, ok := channelGroupMap[key]; ok {
		return group
	}
	// Unrecognised channel value → "Other" instead of title-casing the raw string
	return "Other"
}