package scanner

// countCode(): html을 스캔해 code별 카운팅.
func countCode(html, pageURL, code string) int {
	n := 0
	for _, f := range ScanURL(html, pageURL).Findings {
		if f.Code == code {
			n++
		}
	}
	return n
}

// findFirst(): code에 해당하는 첫 탐지.
func findFirst(html, pageURL, code string) (Finding, bool) {
	for _, f := range ScanURL(html, pageURL).Findings {
		if f.Code == code {
			return f, true
		}
	}
	return Finding{}, false
}
