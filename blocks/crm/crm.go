package crm

import "strconv"

func pct(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) + "%" }
