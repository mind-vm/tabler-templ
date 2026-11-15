package reference

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

// Formats of GS-73: 1.234,50 €, TT.MM.JJJJ, 07:30, ranges with an en dash.

// Decimal is a quantity in thousandths (12,5 = 12500).
type Decimal int64

func DecimalOf(f float64) Decimal { return Decimal(math.Round(f * 1000)) }

// Times multiplies the quantity by a price in cents and rounds to cents.
func (d Decimal) Times(cents int64) int64 {
	return int64(math.Round(float64(d) * float64(cents) / 1000))
}

// String is the German form without trailing zeros: 12, 12,5, 1.200,25.
func (d Decimal) String() string {
	neg := d < 0
	if neg {
		d = -d
	}
	whole, frac := int64(d)/1000, int64(d)%1000
	s := thousands(whole)
	if frac != 0 {
		f := strings.TrimRight(strconv.FormatInt(1000+frac, 10)[1:], "0")
		s += "," + f
	}
	if neg {
		s = "-" + s
	}
	return s
}

var errNumber = errors.New("keine Zahl")

// ParseDecimal reads "1.234,5", "12,5" or "12.5" (one dot with up to two digits after it is a decimal point).
func ParseDecimal(s string) (Decimal, error) {
	f, err := parseNumber(s)
	if err != nil {
		return 0, err
	}
	return DecimalOf(f), nil
}

// ParseCents reads an amount in euros ("48,50", "1.200", "48,50 €") as cents.
func ParseCents(s string) (int64, error) {
	f, err := parseNumber(strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "€")))
	if err != nil {
		return 0, err
	}
	return int64(math.Round(f * 100)), nil
}

func parseNumber(s string) (float64, error) {
	s = strings.ReplaceAll(strings.TrimSpace(s), " ", "")
	if s == "" {
		return 0, errNumber
	}
	if strings.Contains(s, ",") {
		s = strings.ReplaceAll(s, ".", "")
		s = strings.Replace(s, ",", ".", 1)
	} else if i := strings.LastIndex(s, "."); i >= 0 && (len(s)-i-1 == 3 || strings.Count(s, ".") > 1) {
		s = strings.ReplaceAll(s, ".", "") // 1.200 is a thousand separator
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, errNumber
	}
	return f, nil
}

func thousands(n int64) string {
	s := strconv.FormatInt(n, 10)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "." + s[i:]
	}
	return s
}

// Euro formats cents: 1.234,50 €.
func Euro(cents int64) string {
	neg := cents < 0
	if neg {
		cents = -cents
	}
	s := thousands(cents/100) + "," + strconv.FormatInt(100+cents%100, 10)[1:] + " €"
	if neg {
		s = "−" + s
	}
	return s
}

// EuroField is the amount in a field, without the sign: 48,50.
func EuroField(cents int64) string { return strings.TrimSuffix(Euro(cents), " €") }

// Date is TT.MM.JJJJ; the zero time is "".
func Date(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02.01.2006")
}

// Period is „01.04.2026 – 31.03.2028“.
func Period(from, to time.Time) string { return Date(from) + " – " + Date(to) }

func parseDate(s string) (time.Time, error) {
	return time.Parse("02.01.2006", strings.TrimSpace(s))
}

// fold makes text comparable for the search: lower case, ä = ae, ö = oe, ü = ue, ß = ss (GS-50).
func fold(s string) string {
	return folder.Replace(strings.ToLower(s))
}

var folder = strings.NewReplacer("ä", "ae", "ö", "oe", "ü", "ue", "ß", "ss")

// matches: every word of the query occurs in one of the fields, ignoring case and umlauts (GS-50).
func matches(query string, fields ...string) bool {
	text := fold(strings.Join(fields, " "))
	for _, w := range strings.Fields(fold(query)) {
		if !strings.Contains(text, w) {
			return false
		}
	}
	return true
}
